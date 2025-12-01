package script_executors

import (
	"fmt"
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/config"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/mapper"
	"k8s-platform-go/internal/util"
	"time"

	"github.com/gin-gonic/gin"
)

// shell 策略执行器

type ShellExecutor struct {
	hostMapper *mapper.HostMapper
}

func (s ShellExecutor) GetFileSuffix() string {
	return ".sh"
}

func (s ShellExecutor) Prepare(script *dto.ExecutorScript) ([]string, string) {
	defaultScript := DefaultScriptExecutor{}
	return defaultScript.Prepare(script)
}

func NewShellExecutor(hostMapper *mapper.HostMapper) *ShellExecutor {
	return &ShellExecutor{
		hostMapper: hostMapper,
	}
}

func (s ShellExecutor) GetType() string {
	return "shell"
}

func (s ShellExecutor) Validate(script *dto.ExecutorScript) error {
	return nil
}

func (s ShellExecutor) Execute(c *gin.Context, script *dto.ExecutorScript) (map[string][]*util.ExecutorResult, error) {
	// 1. 根据主机id列表去查找所有的详细的信息
	if len(script.HostIds) <= 0 {
		common.Fail(c, common.HostIdsEmpty)
		return nil, common.HostIdsEmpty
	}
	hosts, err := s.hostMapper.SelectByIds(script.HostIds)
	if err != nil {
		return nil, common.NewErrorCode(500, "查询主机列表失败")
	}
	// 2. 进行ssh连接到对应的主机
	hostInfos := make([]*util.HostInfo, 0, len(hosts))
	for _, host := range hosts {
		hostInfo := &util.HostInfo{
			Address:  host.Address,
			Port:     int(host.HostPort),
			Username: host.Username,
			Password: *host.HostPassword,
			Timeout:  time.Duration(script.TimeOut) * time.Second,
		}
		hostInfos = append(hostInfos, hostInfo)
	}

	// 3. 添加批量执行的命令和参数
	commands, scriptFileName := s.Prepare(script)
	executorScript := fmt.Sprintf("bash %s %s", scriptFileName, script.Parameters)
	commands = append(commands, executorScript)
	return util.BatchExecuteCommands(hostInfos, commands, &util.BatchConfig{
		MaxConcurrent: config.GetGlobalConfig().Job.Script.GlobalTimeout,
		GlobalTimeout: time.Duration(script.TimeOut),
	})
}
