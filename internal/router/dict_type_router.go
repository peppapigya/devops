package router

import (
	"k8s-platform-go/internal/util/wireinfo"

	"github.com/gin-gonic/gin"
)

func InitDictTypeRouter(rg *gin.RouterGroup) {
	ctl := wireinfo.InitializeDictTypeController()
	g := rg.Group("/sysDictType")
	{
		g.POST("/page", ctl.Page)
		g.GET("/:id", ctl.Detail)
		g.POST("/", ctl.Create)
		g.PUT("/:id", ctl.Update)
		g.DELETE("/:id", ctl.Remove)
		g.PUT("/:id/status", ctl.UpdateStatus)
	}
}
