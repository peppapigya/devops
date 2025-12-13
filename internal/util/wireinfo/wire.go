//go:build wireinject
// +build wireinject

package wireinfo

import (
	"k8s-platform-go/internal/config/db"
	host3 "k8s-platform-go/internal/controller/host"
	"k8s-platform-go/internal/controller/job"
	"k8s-platform-go/internal/controller/system"
	"k8s-platform-go/internal/dal/redis"
	host2 "k8s-platform-go/internal/mapper/host"
	job2 "k8s-platform-go/internal/mapper/job"
	system2 "k8s-platform-go/internal/mapper/system"
	"k8s-platform-go/internal/service/host"
	job3 "k8s-platform-go/internal/service/job"
	system3 "k8s-platform-go/internal/service/system"
	strategy "k8s-platform-go/internal/strategy/script_executors"

	"github.com/google/wire"
)

// 初始化用户控制器
func InitializeUserController() *system.UserController {
	wire.Build(db.NewDB, db.InitRedis, redis.NewClient, system2.NewUserMapper, system3.NewUserService, system.NewUserController)
	return &system.UserController{}
}

// 初始化主机控制器
func InitializeHostController() *host3.HostController {
	wire.Build(db.NewDB, host2.NewHostMapper, host.NewHostService, host3.NewHostController)
	return &host3.HostController{}
}

// 初始化部门控制器
func InitializeDeptController() *system.DeptController {
	wire.Build(db.NewDB, system2.NewDeptMapper, system3.NewDeptService, system.NewDeptController)
	return &system.DeptController{}
}

// 初始化字典类型控制器
func InitializeDictTypeController() *system.DictTypeController {
	wire.Build(db.NewDB, system2.NewDictTypeMapper, system3.NewDictTypeService, system.NewDictTypeController)
	return &system.DictTypeController{}
}

// 初始化菜单控制器
func InitializeMenuController() *system.MenuController {
	wire.Build(db.NewDB, system2.NewMenuMapper, system3.NewMenuService, system.NewMenuController)
	return &system.MenuController{}
}

// 初始化作业脚本控制器
func InitializeJobScriptController() *job.JobScriptController {
	wire.Build(db.NewDB, host2.NewHostMapper, host.NewHostService, job2.NewJobPlanLogMapper, job3.NewJobPlanLogService, job2.NewJobScriptMapper, job3.NewJobScriptService, job.NewJobScriptController)
	return &job.JobScriptController{}
}

// 初始化作业计划控制器
func InitializeJobPlanController() *job.JobPlanController {
	wire.Build(db.NewDB, job2.NewJobPlanMapper, job2.NewJobPlanScriptMapper, job3.NewJobPlanService, job.NewJobPlanController)
	return &job.JobPlanController{}
}

// 初始化定时任务控制器
func InitializeJobScheduledTaskController() *job.JobScheduledTaskController {
	wire.Build(db.NewDB, job2.NewJobScheduledTaskMapper, job3.NewJobScheduledTaskService, job.NewJobScheduledTaskController)
	return &job.JobScheduledTaskController{}
}

// 初始化作业日志控制器
func InitializeJobPlanLogController() *job.JobPlanLogController {
	wire.Build(db.NewDB, job2.NewJobPlanLogMapper, job3.NewJobPlanLogService, job.NewJobPlanLogController)
	return &job.JobPlanLogController{}
}

// 初始化shell脚本控制器
func InitializeShellScriptExecutor() *strategy.ShellExecutor {
	wire.Build(db.NewDB, host2.NewHostMapper, strategy.NewShellExecutor)
	return &strategy.ShellExecutor{}
}
