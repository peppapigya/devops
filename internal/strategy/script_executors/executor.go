package script_executors

import (
	"fmt"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

type ScriptExecutor interface {

	// GetType 获取执行器类型
	GetType() string

	// GetFileSuffix 文件后缀的格式
	GetFileSuffix() string

	// Validate 校验脚本
	Validate(script *dto.ExecutorScript) error

	// Prepare 准备执行脚本
	Prepare(script *dto.ExecutorScript) ([]string, string)

	// Execute 执行脚本
	Execute(c *gin.Context, script *dto.ExecutorScript) (map[string][]*util.ExecutorResult, error)
}

// DefaultScriptExecutor 提供默认实现
type DefaultScriptExecutor struct{}

// Prepare 执行准备工作，将所有需要的前置指令拼接返回
func (defaultScript *DefaultScriptExecutor) Prepare(script *dto.ExecutorScript) ([]string, string) {
	commands := []string{"set -e"}
	// 1.设置工作目录,如果为指定，默认是/tmp
	if script.WorkDir != "" {
		commands = append(commands, fmt.Sprintf("cd %s", defaultScript.escapePath(script.WorkDir)))
	} else {
		commands = append(commands, "cd /tmp")
	}

	envs := util.SplitString(*script.Env, "\n")
	// 2. 设置环境变量
	for _, env := range envs {
		if strings.Contains(env, "=") {
			commands = append(commands, fmt.Sprintf("export %s", env))
		}
	}

	tempFileName := fmt.Sprintf("script_%s%s", script.Name, defaultScript.GetFileSuffix())
	// 3. 将脚本内容生成临时文件
	commands = append(commands, fmt.Sprintf("cat > %s <<'EOF'\n%s\nEOF", tempFileName, script.Content))
	return commands, tempFileName
}

// escapePath 转义路径中的特殊字符
func (defaultScript *DefaultScriptExecutor) escapePath(path string) string {
	return `"` + strings.ReplaceAll(path, `"`, `\"`) + `"`
}
func (defaultScript *DefaultScriptExecutor) GetFileSuffix() string {
	return ""
}
