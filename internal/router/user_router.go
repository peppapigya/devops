package router

import (
	"k8s-platform-go/internal/util/wireinfo"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(rg *gin.RouterGroup) {
	userController := wireinfo.InitializeUserController()
	userGroup := rg.Group("/sysUser")
	{
		userGroup.GET("/:id", userController.GetUserDOById)
		userGroup.POST("/login", userController.Login)
		userGroup.POST("/refresh", userController.RefreshToken)
		userGroup.POST("/page", userController.GetUserPage)
		userGroup.PUT("/status/", userController.UpdateUserStatus)
		userGroup.PUT("/", userController.UpdateUser)
		userGroup.GET("/captcha", userController.GetCaptcha)
	}
}
