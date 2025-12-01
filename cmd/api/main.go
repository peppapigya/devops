package main

// @title K8s Platform API
// @version 1.0
// @description This is a k8s management platform API documentation
// @host localhost:8081
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
import (
	"fmt"
	_ "k8s-platform-go/docs"
	"k8s-platform-go/internal/common"

	"k8s-platform-go/internal/common/exception"
	"k8s-platform-go/internal/common/middleware"
	"k8s-platform-go/internal/config"
	"k8s-platform-go/internal/config/db"

	"k8s-platform-go/internal/router"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	// 加载配置文件
	err := config.LoadConfig()
	if err != nil {
		log.Printf("config load faild: %v", err)
		return
	}
	globalConfig := config.GetGlobalConfig()
	// 初始化数据库连接
	db.NewDB()
	// 关闭数据库连接
	defer db.CloseDB()

	// Redis配置
	db.InitRedis()
	defer db.CloseRedis()

	r := gin.Default()
	binding.Validator = &common.DefaultValidator{}
	// 设置中间件
	setMiddleware(r, globalConfig)
	// 设置全局路由
	router.InitRouter(r)

	err = r.Run(fmt.Sprintf(":%s", globalConfig.Server.Port))
	if err != nil {
		log.Printf("server run faild: %v", err)
		return
	}
}

// 生效的顺序是根据添加的顺序决定的
func setMiddleware(router *gin.Engine, globalConfig *config.GlobalConfig) {

	// 全局请求日志处理
	router.Use(middleware.RequestLogHandler())
	// 跨域配置
	router.Use(middleware.Cors())
	// 配置全局异常处理
	router.Use(exception.GlobalExceptionHandler())
	// 认证
	router.Use(middleware.Authenticate(globalConfig.Jwt.ExcludePaths...))
}
