package script_executors

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/util"

	"github.com/gin-gonic/gin"
)

type PythonExecutor struct{}

func (p PythonExecutor) ExecuteStream(c *gin.Context, script *dto.ExecutorScript, onEvent func(util.StreamEvent)) (map[string][]*util.ExecutorResult, error) {
	//TODO implement me
	panic("implement me")
}

func (p PythonExecutor) GetFileSuffix() string {
	//TODO implement me
	panic("implement me")
}

func (p PythonExecutor) Prepare(script *dto.ExecutorScript) ([]string, string) {
	return []string{"python"}, script.Content
}

func (p PythonExecutor) GetType() string {
	return "python"
}

func (p PythonExecutor) Validate(script *dto.ExecutorScript) error {
	return nil
}

func (p PythonExecutor) Execute(c *gin.Context, script *dto.ExecutorScript) (map[string][]*util.ExecutorResult, error) {
	return nil, nil
}
