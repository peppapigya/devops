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
	"io"
	_ "k8s-platform-go/docs"
	"k8s-platform-go/internal/common"
	"os"
	"path/filepath"

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
	r := gin.Default()
	// 加载配置文件
	err := config.LoadConfig()
	if err != nil {
		log.Printf("config load faild: %v", err)
		return
	}
	globalConfig := config.GetGlobalConfig()

	configLog(globalConfig)

	// 初始化数据库连接
	db.NewDB()
	// 关闭数据库连接
	defer db.CloseDB()

	// Redis配置
	db.InitRedis()
	defer db.CloseRedis()

	binding.Validator = &common.DefaultValidator{}
	// 设置中间件
	setMiddleware(r, globalConfig)
	// 设置全局路由
	router.InitRouter(r)
	fmt.Printf("代理地址%v\n", globalConfig.Server.TrustedProxies)
	// 设置代理，避免gin启动告警
	if globalConfig.Server.TrustedProxies != nil {
		err = r.SetTrustedProxies(globalConfig.Server.TrustedProxies)
		if err != nil {
			_ = fmt.Errorf("set trusted proxies error: %v", err)
		}
	}

	fmt.Printf("启动环境: %s\n", gin.Mode())
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

// 配置全局日志
func configLog(globalConfig *config.GlobalConfig) {
	if !globalConfig.Log.Enable {
		log.Println("未开启系统日志...")
		return
	}
	// 禁用终端颜色
	gin.DisableConsoleColor()
	log.Println("已开启系统日志...")
	// 先创建目录
	dir := filepath.Dir(globalConfig.Log.Path)
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		log.Printf("create log dir error: %v", err)
		panic("create log dir error")
	}
	file, _ := os.Create(globalConfig.Log.Path)
	log.Printf("日志文件路径: %s\n", globalConfig.Log.Path)

	if globalConfig.Log.Stdout {
		gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	} else {
		// 只将日志输出到文件(正式环境使用),这个不是追加写文件，而是覆盖写文件
		gin.DefaultWriter = io.MultiWriter(file)
	}
}
