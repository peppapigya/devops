package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 配置所有路由
	api := r.Group("/api/v1")
	{
		InitUserRouter(api)
	}
}
