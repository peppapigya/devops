//go:build wireinject
// +build wireinject

package wireinfo

import (
	"k8s-platform-go/internal/config/db"
	"k8s-platform-go/internal/controller"
	"k8s-platform-go/internal/mapper"
	"k8s-platform-go/internal/service"

	"github.com/google/wire"
)

func InitializeUserController() *controller.UserController {
	wire.Build(db.NewDB, mapper.NewUserMapper, service.NewUserService, controller.NewUserController)
	return &controller.UserController{}
}
