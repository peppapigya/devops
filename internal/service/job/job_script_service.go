package job

import (
	"bytes"
	"fmt"
	"io"
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper/job"
	"k8s-platform-go/internal/service/host"
	"k8s-platform-go/internal/strategy"
	"k8s-platform-go/internal/util"
	"k8s-platform-go/pkg/consts"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
)

type JobScriptService struct {
	jobScriptMapper   *job.JobScriptMapper
	jobPlanLogService *JobPlanLogService
	hostService       *host.HostService
}

func NewJobScriptService(jobPlanLogService *JobPlanLogService, jobScriptMapper *job.JobScriptMapper, hostService *host.HostService) *JobScriptService {
	return &JobScriptService{
		jobScriptMapper:   jobScriptMapper,
		jobPlanLogService: jobPlanLogService,
		hostService:       hostService,
	}
}

func (s *JobScriptService) CreateJobScript(c *gin.Context, req dto.JobScriptSaveRequest) {
	script := &model.JobScript{
		Name:       req.Name,
		Type:       req.Type,
		Category:   req.Category,
		Content:    req.Content,
		Parameters: req.Parameters,
		Timeout:    uint32(req.Timeout),
		WorkDir:    &req.WorkDir,
		Env:        &req.Env,
	}
	if script.Category == "" {
		script.Category = "default"
	}
	err := s.jobScriptMapper.InsertJobScript(script)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScriptService) UpdateJobScript(c *gin.Context, req dto.JobScriptSaveRequest) {
	script, err := s.jobScriptMapper.GetJobScriptById(req.ID)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	if script == nil {
		common.Fail(c, common.BadRequest)
		return
	}

	script.Name = req.Name
	script.Type = req.Type
	script.Category = req.Category
	script.Content = req.Content
	script.Parameters = req.Parameters
	script.Timeout = uint32(req.Timeout)
	script.WorkDir = &req.WorkDir
	script.Env = &req.Env

	err = s.jobScriptMapper.UpdateJobScript(script)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScriptService) DeleteJobScript(c *gin.Context, id int64) {
	err := s.jobScriptMapper.DeleteJobScript(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScriptService) GetJobScriptPage(c *gin.Context, req dto.JobScriptPageRequest) {
	pageResult, err := s.jobScriptMapper.GetJobScriptPage(req)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, pageResult)
}

func (s *JobScriptService) GetJobScriptById(id int64) (*model.JobScript, error) {
	script, err := s.jobScriptMapper.GetJobScriptById(id)
	if err != nil {
		return nil, err
	}
	return script, nil
}

func (s *JobScriptService) GetJobScriptSelect(condition string) ([]*model.JobScript, error) {

	jobScripts, err := s.jobScriptMapper.SelectListByCondition(condition)

	if err != nil {
		return nil, err
	}

	return jobScripts, nil
}

func (s *JobScriptService) ExecuteJobScript(c *gin.Context, script dto.ExecutorScript) (map[string][]*util.ExecutorResult, error) {
	// 如果使用的脚本库，将脚本数据库中的数据填充到请求参数中
	if script.ScriptId != 0 {
		if err := s.getScriptFromDatabase(&script); err != nil {
			return nil, err
		}
	}

	// 获取对应的执行器
	factory := strategy.GetExecutor(script.Type)

	if factory == nil {
		return nil, common.ScriptFactoryNotExist
	}

	execute, err := factory.Execute(c, &script)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", execute)
	// 添加执行日志
	if err := s.insertJobLogs(execute, consts.ScriptExecuteTypeSingle, consts.JobTypeManual); err != nil {
		return nil, err
	}
	return execute, nil
}

// 添加执行日志
func (s *JobScriptService) insertJobLogs(executeResult map[string][]*util.ExecutorResult, method string, typeName string) error {
	jobPlanLogs := make([]*model.JobPlanLog, 0)
	for _, results := range executeResult {
		var builder strings.Builder
		var totalTime int64
		status := true
		hostId := results[0].HostId
		for _, res := range results {
			builder.WriteString(res.Output)
			builder.WriteString("\n")
			status = status && res.Success
			totalTime += res.Duration
		}
		output := builder.String()
		jobPlanLog := &model.JobPlanLog{
			HostID:     uint32(hostId),
			Method:     method,
			Type:       typeName,
			Status:     status,
			TotalTime:  strconv.FormatInt(totalTime, 10),
			ReturnCode: int32(results[len(results)-1].ExitCode),
			Output:     &output,
		}
		jobPlanLogs = append(jobPlanLogs, jobPlanLog)
	}
	return s.jobPlanLogService.InsertJobPlanLogBatch(jobPlanLogs)
}

// 从数据库获取脚本信息并填充到请求参数中
func (s *JobScriptService) getScriptFromDatabase(script *dto.ExecutorScript) error {
	// 查询脚本信息
	scriptInfo, err := s.GetJobScriptById(script.ScriptId)
	if err != nil {
		return common.ScriptNotExist
	}
	script.Content = scriptInfo.Content
	script.Type = scriptInfo.Type
	script.Name = scriptInfo.Name
	script.Content = scriptInfo.Content
	if script.Parameters == "" {
		script.Parameters = scriptInfo.Parameters
	}
	script.TimeOut = int(scriptInfo.Timeout)
	if script.WorkDir == "" {
		script.WorkDir = *scriptInfo.WorkDir
	}
	if script.Env == nil {
		script.Env = scriptInfo.Env
	}
	if script.WorkDir == "" {
		script.WorkDir = *scriptInfo.WorkDir
	}
	return nil
}

func (s *JobScriptService) DistributeJobScript(distribute dto.DistributeJobScript) (map[string][]*dto.DistributeResult, error) {
	// 根据主机id查询主机的所有信息
	hosts, err := s.hostService.GetHostByIds(distribute.HostIds)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	resultChan := make(chan dto.DistributeResult, len(hosts))
	results := make(map[string][]*dto.DistributeResult)
	// 准备需要传输的内容
	content, fileName, err := s.prepareContent(distribute)
	if err != nil {
		return nil, err
	}

	// 分发给每个主机
	for _, hostInfo := range hosts {
		wg.Add(1)
		go func(hostInfo util.HostInfo) {
			defer wg.Done()
			result := dto.DistributeResult{
				HostID:  uint32(hostInfo.ID),
				Address: hostInfo.Address,
				Success: true,
				Message: fmt.Sprintf("分发任务成功，主机ID：%d", hostInfo.ID),
			}
			start := time.Now()
			defer func() {
				result.Duration = time.Since(start).String()
				resultChan <- result
			}()

			// 分发单个主机
			if err := s.distribute(hostInfo, distribute, content, fileName); err != nil {
				result.Message = err.Error()
				result.Success = false
				return
			}
		}(*hostInfo)
	}

	wg.Wait()
	close(resultChan)

	for result := range resultChan {
		results[result.Address] = append(results[result.Address], &result)
	}
	return results, nil
}

// 准备要传输的内容
func (s *JobScriptService) prepareContent(distribute dto.DistributeJobScript) (io.Reader, string, error) {
	var content io.Reader
	var fileName string
	// 如果使用的脚本库
	if distribute.Id > 0 {
		scriptContent, err := s.getScriptContentById(distribute.Id)
		if err != nil {
			return nil, "", err
		}
		content = strings.NewReader(scriptContent)
		fileName = fmt.Sprintf("script_%d.sh", distribute.Id)
	} else {
		// 如果传输的是文件
		file, err := distribute.File.Open()
		if err != nil {
			return nil, "", err
		}
		defer func() { _ = file.Close() }()
		// 将文件内容读入内存
		fileContent, err := io.ReadAll(file)
		if err != nil {
			return nil, "", err
		}
		content = bytes.NewReader(fileContent)
		fileName = distribute.File.Filename
	}
	return content, fileName, nil
}

// 从数据库获取脚本内容
func (s *JobScriptService) getScriptContentById(scriptId int64) (string, error) {
	script, err := s.jobScriptMapper.GetJobScriptById(scriptId)
	if err != nil {
		return "", err
	}
	return script.Content, nil
}

// 分发脚本到单个主机
func (s *JobScriptService) distribute(hostInfo util.HostInfo, distribute dto.DistributeJobScript, content io.Reader, fileName string) error {
	//  建立ssh连接
	connection, err := hostInfo.Connection()
	defer func() { _ = connection.Close() }()
	if err != nil {
		return fmt.Errorf("建立ssh连接失败: %s", err)
	}
	filePath := path.Join(distribute.RemotePath, fileName)
	// 创建sftp
	sftpClient, err := sftp.NewClient(connection)
	defer func() { _ = sftpClient.Close() }()
	if err != nil {
		return fmt.Errorf("创建sftp客户端失败: %s", err)
	}

	// 如果需要备份文件
	if distribute.Backup {
		if err := s.backupFile(sftpClient, filePath); err != nil {
			return fmt.Errorf("备份文件失败: %s", err)
		}
	}

	// 是否需要覆盖文件，不覆盖文件的话就使用redis缓存的值
	if err := s.overwriteFile(sftpClient, filePath, distribute.Overwrite); err != nil {
		return err
	}

	// 创建远端文件
	remoteFile, err := sftpClient.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建远程文件失败: %s", err)
	}
	defer func() { _ = remoteFile.Close() }()

	// 传输文件内容
	_, err = util.Copy(remoteFile, content, nil)
	if err != nil {
		return fmt.Errorf("传输文件内容失败: %s", err)
	}

	// 设置文件权限
	if err := s.setFilePermission(sftpClient, filePath, distribute.Permission); err != nil {
		return fmt.Errorf("设置文件权限失败: %s", err.Error())
	}

	// 设置文件所有者，todo @dxg 目前前端还没实现该功能,后续完善

	return nil
}

func (s *JobScriptService) backupFile(sftpClient *sftp.Client, remotePath string) error {
	// 检查目录是否存在，不存在则创建
	if err := sftpClient.MkdirAll(path.Dir(remotePath)); err != nil {
		return err
	}
	// 文件不存在没必要备份
	if _, err := sftpClient.Stat(remotePath); os.IsNotExist(err) {
		return nil
	}
	// 备份文件
	backupFileName := fmt.Sprintf("%s.%s.backup", remotePath, time.Now().Format("20251206_142001"))
	srcFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %s", err)
	}
	defer func() { _ = srcFile.Close() }()

	destFile, err := sftpClient.Create(backupFileName)
	if err != nil {
		return fmt.Errorf("创建备份文件失败: %s", err)
	}
	defer func() { _ = destFile.Close() }()

	// 复制文件
	_, err = util.Copy(destFile, srcFile, nil)
	return err
}

// 覆盖文件
func (s *JobScriptService) overwriteFile(sftpClient *sftp.Client, remotePath string, overWrite bool) error {
	if err := sftpClient.MkdirAll(path.Dir(remotePath)); err != nil {
		return err
	}
	// 如果文件不存在或者不需要覆盖，直接返回
	if _, err := sftpClient.Stat(remotePath); err == nil && !overWrite {
		return fmt.Errorf("文件已存在请选择覆盖模式：%s", err)
	}
	return nil
}

// 设置文件权限
func (s *JobScriptService) setFilePermission(client *sftp.Client, remotePath, permission string) error {
	perm, err := strconv.ParseUint(permission, 8, 32)
	if err != nil {
		return fmt.Errorf("无效的权限值: %s", err)
	}

	return client.Chmod(remotePath, os.FileMode(perm))
}
