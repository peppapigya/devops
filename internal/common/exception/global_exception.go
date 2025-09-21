package exception

import (
	"log"
	"runtime/debug"

	"k8s-platform-go/internal/common"

	"github.com/gin-gonic/gin"
)

// GlobalExceptionHandler 全局异常捕获
func GlobalExceptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印堆栈信息
				log.Printf("panic: %v", debug.Stack())
				if err, ok := err.(*common.ErrorCode); ok {
					log.Printf("解析参数失败: %v", err)
					common.Fail(c, err)
				} else {
					log.Printf("解析参数失败: %v", err)
					common.Fail(c, common.NewErrorCode(500, "服务器异常"))
				}
				// 阻止程序继续往下执行，直接返回，避免异常再次被处理
				c.Abort()
			}
		}()
		c.Next()
	}
}
