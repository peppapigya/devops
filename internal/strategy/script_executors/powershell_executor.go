package script_executors

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/util"

	"github.com/gin-gonic/gin"
)

type PowerShellExecutor struct{}

func (p PowerShellExecutor) ExecuteStream(c *gin.Context, script *dto.ExecutorScript, onEvent func(util.StreamEvent)) (map[string][]*util.ExecutorResult, error) {
	//TODO implement me
	panic("implement me")
}

func (p PowerShellExecutor) GetFileSuffix() string {
	return ""
}

func (p PowerShellExecutor) Prepare(script *dto.ExecutorScript) ([]string, string) {
	return nil, ""
}

func (p PowerShellExecutor) GetType() string {
	return "powershell"
}

func (p PowerShellExecutor) Validate(script *dto.ExecutorScript) error {
	return nil
}

func (p PowerShellExecutor) Execute(c *gin.Context, script *dto.ExecutorScript) (map[string][]*util.ExecutorResult, error) {
	return nil, nil
}
