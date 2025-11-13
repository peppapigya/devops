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
