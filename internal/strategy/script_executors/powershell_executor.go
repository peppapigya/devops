package script_executors

import (
	"k8s-platform-go/internal/dal/dto"

	"github.com/gin-gonic/gin"
)

type PowerShellExecutor struct{}

func (p PowerShellExecutor) GetType() string {
	//TODO implement me
	panic("implement me")
}

func (p PowerShellExecutor) Validate(script *dto.ExecutorScript) error {
	//TODO implement me
	panic("implement me")
}

func (p PowerShellExecutor) Execute(c *gin.Context, script *dto.ExecutorScript) (ExecutorResult, error) {
	//TODO implement me
	panic("implement me")
}
