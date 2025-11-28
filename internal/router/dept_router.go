package router

import (
	"k8s-platform-go/internal/util/wireinfo"

	"github.com/gin-gonic/gin"
)

func InitDeptRouter(rg *gin.RouterGroup) {
	ctl := wireinfo.InitializeDeptController()
	g := rg.Group("/sysDept")
	{
		g.POST("/page", ctl.Page)
		g.GET("/tree", ctl.Tree)
		g.GET("/:id", ctl.Detail)
		g.POST("/", ctl.Create)
		g.PUT("/", ctl.Update)
		g.DELETE("/:id", ctl.Remove)
		g.PUT("/:id/status", ctl.UpdateStatus)
	}
}
