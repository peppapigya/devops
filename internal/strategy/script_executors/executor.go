package script_executors

import (
	"k8s-platform-go/internal/dal/dto"

	"github.com/gin-gonic/gin"
)

type ScriptExecutor interface {

	// GetType 获取执行器类型
	GetType() string

	// Validate 校验脚本
	Validate(script *dto.ExecutorScript) error

	// Execute 执行脚本
	Execute(c *gin.Context, script *dto.ExecutorScript) (ExecutorResult, error)
}

type ExecutorResult struct {
	Success   bool
	Output    string
	Error     error
	ExistCode int
	Duration  int64 //单位毫秒
}
