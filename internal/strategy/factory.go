package strategy

import (
	"k8s-platform-go/internal/config/db"
	"k8s-platform-go/internal/mapper"
	strategy "k8s-platform-go/internal/strategy/script_executors"
)

type ScriptExecutorFactory struct {
	Executors map[string]strategy.ScriptExecutor
}

func NewExecutorFactory() *ScriptExecutorFactory {
	factory := &ScriptExecutorFactory{
		Executors: make(map[string]strategy.ScriptExecutor),
	}
	factory.Register(strategy.NewShellExecutor(mapper.NewHostMapper(db.NewDB())))
	factory.Register(strategy.PythonExecutor{})
	factory.Register(strategy.PowerShellExecutor{})

	return factory
}

func (factory *ScriptExecutorFactory) Register(executor strategy.ScriptExecutor) {
	factory.Executors[executor.GetType()] = executor
}

var GlobalFactory *ScriptExecutorFactory

func InitGlobalFactory() {
	GlobalFactory = NewExecutorFactory()
}

func GetExecutor(typeName string) strategy.ScriptExecutor {
	if typeName == "" {
		return nil
	}
	if GlobalFactory == nil {
		InitGlobalFactory()
	}
	return GlobalFactory.Executors[typeName]
}
