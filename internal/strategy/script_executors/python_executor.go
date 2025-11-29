package script_executors

import (
	"k8s-platform-go/internal/dal/dto"

	"github.com/gin-gonic/gin"
)

type PythonExecutor struct{}

func (p PythonExecutor) GetType() string {
	//TODO implement me
	panic("implement me")
}

func (p PythonExecutor) Validate(script *dto.ExecutorScript) error {
	//TODO implement me
	panic("implement me")
}

func (p PythonExecutor) Execute(c *gin.Context, script *dto.ExecutorScript) (ExecutorResult, error) {
	//TODO implement me
	panic("implement me")
}
