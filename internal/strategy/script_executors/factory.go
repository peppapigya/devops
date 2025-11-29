package script_executors

type ScriptExecutorFactory struct {
	executors map[string]ScriptExecutor
}

// NewExecutorFactory 初始化脚本执行工厂
func NewExecutorFactory() *ScriptExecutorFactory {
	factory := &ScriptExecutorFactory{
		executors: make(map[string]ScriptExecutor),
	}
	// 注册执行器
	factory.RegisterFactory(&ShellExecutor{})
	factory.RegisterFactory(&PowerShellExecutor{})
	factory.RegisterFactory(&PythonExecutor{})

	return factory

}

// RegisterFactory 注册脚本执行器
func (factory *ScriptExecutorFactory) RegisterFactory(executor ScriptExecutor) {
	factory.executors[executor.GetType()] = executor
}

var GlobalFactory *ScriptExecutorFactory

func InitGlobalFactory() {
	GlobalFactory = NewExecutorFactory()
}

func GetExecutor(typeName string) ScriptExecutor {
	if typeName == "" {
		return nil
	}
	if GlobalFactory == nil {
		InitGlobalFactory()
	}

	return GlobalFactory.executors[typeName]
}
