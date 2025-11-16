//go:build wireinject
// +build wireinject

package wireinfo

import (
	"k8s-platform-go/internal/config/db"
	"k8s-platform-go/internal/controller"
	"k8s-platform-go/internal/dal/redis"
	"k8s-platform-go/internal/mapper"
	"k8s-platform-go/internal/service"

	"github.com/google/wire"
)

// 初始化用户控制器
func InitializeUserController() *controller.UserController {
	wire.Build(db.NewDB, db.InitRedis, redis.NewClient, mapper.NewUserMapper, service.NewUserService, controller.NewUserController)
	return &controller.UserController{}
}

// 初始化主机控制器
func InitializeHostController() *controller.HostController {
	wire.Build(db.NewDB, mapper.NewHostMapper, service.NewHostService, controller.NewHostController)
	return &controller.HostController{}
}

// 初始化部门控制器
func InitializeDeptController() *controller.DeptController {
	wire.Build(db.NewDB, mapper.NewDeptMapper, service.NewDeptService, controller.NewDeptController)
	return &controller.DeptController{}
}

// 初始化字典类型控制器
func InitializeDictTypeController() *controller.DictTypeController {
	wire.Build(db.NewDB, mapper.NewDictTypeMapper, service.NewDictTypeService, controller.NewDictTypeController)
	return &controller.DictTypeController{}
}

// 初始化菜单控制器
func InitializeMenuController() *controller.MenuController {
	wire.Build(db.NewDB, mapper.NewMenuMapper, service.NewMenuService, controller.NewMenuController)
	return &controller.MenuController{}
}
