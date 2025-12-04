package job

import (
	"fmt"
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper/job"
	"k8s-platform-go/internal/strategy"
	"k8s-platform-go/internal/util"

	"github.com/gin-gonic/gin"
)

type JobScriptService struct {
	jobScriptMapper *job.JobScriptMapper
}

func NewJobScriptService(jobScriptMapper *job.JobScriptMapper) *JobScriptService {
	return &JobScriptService{
		jobScriptMapper: jobScriptMapper,
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
	return execute, nil
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

func (s *JobScriptService) DistributeJobScript(c *gin.Context, distribute dto.DistributeJobScript) (map[string][]*dto.DistributeResult, error) {
	return nil, nil
}
