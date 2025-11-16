package router

import (
	"k8s-platform-go/internal/util/wireinfo"

	"github.com/gin-gonic/gin"
)

func InitMenuRouter(rg *gin.RouterGroup) {
	ctl := wireinfo.InitializeMenuController()
	ctl.Seed()
	g := rg.Group("/sysMenu")
	{
		// 基础路由
		g.GET("/routes", ctl.Routes)
		g.GET("/tree", ctl.Tree)
		g.GET("/options", ctl.Options)
		g.GET("/:id", ctl.GetById)

		// CRUD操作
		g.POST("", ctl.Create)
		g.PUT("/:id", ctl.Update)
		g.DELETE("/:id", ctl.Delete)
		g.GET("/list", ctl.List)
	}
}
