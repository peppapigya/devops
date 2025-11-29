//go:build wireinject
// +build wireinject

package wireinfo

import (
	"k8s-platform-go/internal/config/db"
	"k8s-platform-go/internal/controller/job"
	"k8s-platform-go/internal/controller/system"
	"k8s-platform-go/internal/dal/redis"
	"k8s-platform-go/internal/mapper"
	"k8s-platform-go/internal/service"

	"github.com/google/wire"
)

// 初始化用户控制器
func InitializeUserController() *system.UserController {
	wire.Build(db.NewDB, db.InitRedis, redis.NewClient, mapper.NewUserMapper, service.NewUserService, system.NewUserController)
	return &system.UserController{}
}

// 初始化主机控制器
func InitializeHostController() *system.HostController {
	wire.Build(db.NewDB, mapper.NewHostMapper, service.NewHostService, system.NewHostController)
	return &system.HostController{}
}

// 初始化部门控制器
func InitializeDeptController() *system.DeptController {
	wire.Build(db.NewDB, mapper.NewDeptMapper, service.NewDeptService, system.NewDeptController)
	return &system.DeptController{}
}

// 初始化字典类型控制器
func InitializeDictTypeController() *system.DictTypeController {
	wire.Build(db.NewDB, mapper.NewDictTypeMapper, service.NewDictTypeService, system.NewDictTypeController)
	return &system.DictTypeController{}
}

// 初始化菜单控制器
func InitializeMenuController() *system.MenuController {
	wire.Build(db.NewDB, mapper.NewMenuMapper, service.NewMenuService, system.NewMenuController)
	return &system.MenuController{}
}

// 初始化作业脚本控制器
func InitializeJobScriptController() *job.JobScriptController {
	wire.Build(db.NewDB, mapper.NewJobScriptMapper, service.NewJobScriptService, job.NewJobScriptController)
	return &job.JobScriptController{}
}

// 初始化作业计划控制器
func InitializeJobPlanController() *job.JobPlanController {
	wire.Build(db.NewDB, mapper.NewJobPlanMapper, mapper.NewJobPlanScriptMapper, service.NewJobPlanService, job.NewJobPlanController)
	return &job.JobPlanController{}
}

// 初始化定时任务控制器
func InitializeJobScheduledTaskController() *job.JobScheduledTaskController {
	wire.Build(db.NewDB, mapper.NewJobScheduledTaskMapper, service.NewJobScheduledTaskService, job.NewJobScheduledTaskController)
	return &job.JobScheduledTaskController{}
}

// 初始化作业日志控制器
func InitializeJobPlanLogController() *job.JobPlanLogController {
	wire.Build(db.NewDB, mapper.NewJobPlanLogMapper, service.NewJobPlanLogService, job.NewJobPlanLogController)
	return &job.JobPlanLogController{}
}
