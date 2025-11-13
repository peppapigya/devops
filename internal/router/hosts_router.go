package router

import (
	"k8s-platform-go/internal/util/wireinfo"

	"github.com/gin-gonic/gin"
)

func InitHostsRouter(rg *gin.RouterGroup) {
	hostController := wireinfo.InitializeHostController()
	hostGroup := rg.Group("/hosts")
	{
		hostGroup.POST("/page", hostController.GetHostPage)
		hostGroup.POST("/", hostController.CreateHost)
		hostGroup.PUT("/", hostController.UpdateHost)
		hostGroup.DELETE("/:id", hostController.DeleteHost)
		hostGroup.POST("/:id/test", hostController.TestConnection)
		hostGroup.POST("/:id/inspect", hostController.InspectHost)
	}
}
