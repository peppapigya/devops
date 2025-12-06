package host

import (
	"fmt"
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper/host"
	"k8s-platform-go/internal/util"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

type HostService struct {
	hostMapper *host.HostMapper
	context    *gin.Context
}

type item struct {
	name string
	cmd  string
}

func NewHostService(hostMapper *host.HostMapper) *HostService {
	return &HostService{
		hostMapper: hostMapper,
	}
}

func (hostService *HostService) GetHostPage(request dto.HostPageRequest, c *gin.Context) {
	pageResult, err := hostService.hostMapper.SelectPageByCondition(request)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, pageResult)
}

func (hostService *HostService) CreateHost(request dto.CreateHostDTO, c *gin.Context) {
	hostInfo := &model.Host{
		HostName:     request.HostName,
		Address:      request.Address,
		HostPort:     request.HostPort,
		Username:     request.Username,
		HostPassword: &request.HostPassword,
		Remark:       &request.Remark,
	}

	err := hostService.hostMapper.InsertHost(hostInfo)
	if err != nil {
		log.Printf("添加主机失败: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (hostService *HostService) UpdateHost(request dto.UpdateHostDTO, c *gin.Context) {
	hostInfo := &model.Host{
		ID:       request.ID,
		HostName: request.HostName,
		Address:  request.Address,
		HostPort: request.HostPort,
		Username: request.Username,
		Remark:   &request.Remark,
	}

	// 如果密码不为空，则更新密码
	if request.HostPassword != "" {
		hostInfo.HostPassword = &request.HostPassword
	}

	err := hostService.hostMapper.UpdateHost(hostInfo)
	if err != nil {
		log.Printf("更新主机失败: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (hostService *HostService) DeleteHost(id int, c *gin.Context) {
	err := hostService.hostMapper.DeleteHostById(id)
	if err != nil {
		log.Printf("删除主机失败: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (hostService *HostService) TestConnection(id int) (string, *common.ErrorCode) {
	// 这里应该实现实际的连接测试逻辑
	// 目前返回模拟结果
	hostInfo, err := hostService.hostMapper.SelectHostById(id)
	if err != nil {
		return "", nil
	}

	if hostInfo == nil {
		return "主机不存在", common.HostNotExist
	}
	remoteHostInfo := &util.HostInfo{
		Address:  hostInfo.Address,
		Port:     int(hostInfo.HostPort),
		Username: hostInfo.Username,
		Password: *hostInfo.HostPassword,
		Timeout:  5 * time.Second,
	}
	// 测试ping主机
	res, err := remoteHostInfo.TestSSHConnection()
	if err != nil || !res {
		return "连接失败", common.HostUnreachable
	}
	// 模拟连接测试结果
	return "连接成功", nil
}

func (hostService *HostService) InspectHost(id int) (interface{}, *common.ErrorCode) {
	hostInfo, err := hostService.hostMapper.SelectHostById(id)
	if err != nil {
		return nil, common.ServerError
	}

	if hostInfo == nil {
		return nil, common.ServerError
	}
	remoteHostInfo := &util.HostInfo{
		Address:  hostInfo.Address,
		Port:     int(hostInfo.HostPort),
		Username: hostInfo.Username,
		Password: *hostInfo.HostPassword,
		Timeout:  5 * time.Second,
	}
	info, err := CheckHostInfo(remoteHostInfo)
	if err != nil {
		return nil, common.NewErrorCode(500, err.Error())
	}
	return info, nil
}

func (hostService *HostService) GetHostSelectList() ([]*dto.HostsSelectInfo, error) {
	list, err := hostService.hostMapper.SelectHostSelectList()
	if err != nil {
		return nil, err
	}
	var hosts []*dto.HostsSelectInfo
	// 转化为dto
	hosts, _ = common.ConvertList(list, func(host *model.Host) (*dto.HostsSelectInfo, error) {
		return &dto.HostsSelectInfo{
			Id:    int(host.ID),
			Label: host.HostName + "(" + host.Address + ")",
			Value: strconv.Itoa(int(host.ID)),
		}, nil
	})
	return hosts, nil

}

func CheckHostInfo(hostInfo *util.HostInfo) (interface{}, error) {
	// 远程连接
	client, err := hostInfo.Connection()
	if err != nil {
		log.Printf("连接失败: %v", err)
		return nil, err
	}
	defer func(client *ssh.Client) {
		_ = client.Close()
	}(client)

	// 声明需要执行的命令
	commands := []item{
		{"主机名：", "hostname"},
		{"操作系统信息：", "cat /etc/os-release"},
		{"内核信息：", "uname -r"},
		{"uptime", "uptime"},
		{"负载信息：", "cat /proc/loadavg || w"},
		{"cpu信息：", "cat /proc/cpuinfo | head -n 50"},
		{"内存：", "free -m"},
		{"磁盘信息：", "df -hP"},
		{"挂载详情", "mount | head -n 200"},
		{"top", "COLUMNS=200 top -b -n1 | head -n 60"},
		{"进程信息：", "ps -eo pid,ppid,user,%cpu,%mem,state,cmd --sort=-%cpu | head -n 50"},
		{"net_if", "ip -o addr"},
		{"netstat", "ss -tuanp | head -n 100 || netstat -an | head -n 100"},
		{"proc_count", "ls /proc | wc -l"},
	}

	result, err := ExecuteCommand(client, commands)
	data, err := util.BeautifyRawData(result)
	if err != nil {
		log.Printf("美化数据失败: %v", err)
		return result, err
	}
	fmt.Printf("数据美化成功: %v", data)
	return data, nil
}

func ExecuteCommand(client *ssh.Client, commands []item) (map[string]string, error) {
	// 执行命令
	result := make(map[string]string, len(commands))
	session, err := client.NewSession()
	if err != nil {
		log.Printf("创建会话失败: %v", err)
		return result, err
	}
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	sessionPool := make(chan struct{}, 5)

	for _, it := range commands {
		wg.Add(1)
		go func(cmdItem item) {
			defer wg.Done()

			sessionPool <- struct{}{}
			defer func() { <-sessionPool }()

			session, err := client.NewSession()
			if err != nil {
				mutex.Lock()
				result[cmdItem.name] = "创建会话失败"
				mutex.Unlock()
				return
			}
			defer func(session *ssh.Session) {
				_ = session.Close()
			}(session)

			output, err := session.Output(cmdItem.cmd)
			mutex.Lock()
			if err != nil {
				result[cmdItem.name] = "执行命令失败: " + err.Error()
			} else {
				result[cmdItem.name] = string(output)
			}
			mutex.Unlock()
		}(it)
	}

	wg.Wait()
	return result, nil
}

func (hostService *HostService) GetHostByIds(ids []uint32) ([]*util.HostInfo, error) {
	hosts, err := hostService.hostMapper.SelectHostByIds(ids)
	if err != nil {
		return nil, err
	}
	list, err := common.ConvertList(hosts, func(host *model.Host) (*util.HostInfo, error) {
		return &util.HostInfo{
			ID:       int(host.ID),
			Address:  host.Address,
			Port:     int(host.HostPort),
			Username: host.Username,
			Password: *host.HostPassword,
			Timeout:  20 * time.Second,
		}, nil
	})
	return list, err

}
