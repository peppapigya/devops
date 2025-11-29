package script_executors

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service"

	"github.com/gin-gonic/gin"
)

// shell 策略执行器

type ShellExecutor struct {
	hostService *service.HostService
}

func NewShellExecutor(hostService *service.HostService) *ShellExecutor {
	return &ShellExecutor{
		hostService: hostService,
	}
}

func (s ShellExecutor) GetType() string {
	//TODO implement me
	panic("implement me")
}

func (s ShellExecutor) Validate(script *dto.ExecutorScript) error {
	//TODO implement me
	panic("implement me")
}

func (s ShellExecutor) Execute(c *gin.Context, script *dto.ExecutorScript) (ExecutorResult, error) {
	//TODO implement me
	panic("implement me")
}
