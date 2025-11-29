package middleware

import (
	"fmt"
	"io"
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/util"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogHandler 请求日志处理
func RequestLogHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		startTime := time.Now()

		url := context.Request.RequestURI
		var params string
		if context.Request.Method == "PUT" || context.Request.Method == "POST" {
			// 读取原始请求体
			body, err := io.ReadAll(context.Request.Body)
			if err != nil {
				log.Printf("请求参数读取失败：%v", err)
				context.Abort()
				return
			}
			params = string(body)

			// 重要：重新设置请求体，以便后续处理可以读取
			context.Request.Body = io.NopCloser(strings.NewReader(string(body)))
		} else {
			params = context.Request.URL.RawQuery
		}
		context.Next()

		endTime := time.Now()
		duration := float64(endTime.Sub(startTime)) / 1000000.0
		var durationStr string
		if duration > 1000 {
			durationStr = fmt.Sprintf("%.3f秒", duration/1000.0)
		} else {
			durationStr = fmt.Sprintf("%.3f毫秒", duration)
		}
		log.Printf("请求开始：Url:[%s] 参数:[%s] 耗时:%s\n ", url, params, durationStr)
	}

}

func Authenticate(excludePaths ...string) gin.HandlerFunc {
	excludePathRegex := make([]*regexp.Regexp, 0)
	for _, path := range excludePaths {
		str := strings.ReplaceAll(path, "*", ".*")
		excludePathRegex = append(excludePathRegex, regexp.MustCompile(str))
	}
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 如果当前请求路径在排除列表中，则不进行权限验证
		for _, regexPath := range excludePathRegex {
			if regexPath.MatchString(path) {
				c.Next()
				return
			}
		}
		token := c.GetHeader(common.TokenKey)
		token, found := strings.CutPrefix(token, "Bearer ")
		if token == "" && !found {
			common.Fail(c, common.UNAUTHORIZED)
			// 终端整个请求链
			c.Abort()
			return
		}
		claims, err := util.ParseToken(token)
		if err != nil {
			common.Fail(c, common.UNAUTHORIZED)
			c.Abort()
			return
		}
		// 将解析的用户信息设置到上下文中
		c.Set(common.UserInfoKey, claims)
		c.Next()
	}
}

// Cors 跨域信息
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		// 1. [必须]接受指定域的请求，可以使用*不加以限制，但不安全
		//context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", context.GetHeader("Origin"))
		fmt.Println(context.GetHeader("Origin"))
		// 2. [必须]设置服务器支持的所有跨域请求的方法
		context.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		// 3. [可选]服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
		context.Header("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Token, Authorization, X-Requested-With, Accept, Origin, Cache-Control, Pragma")
		// 4. [可选]设置XMLHttpRequest的响应对象能拿到的额外字段
		context.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token,Authorization,X-Requested-With, Accept, Origin, Cache-Control")
		// 5. [可选]是否允许后续请求携带认证信息Cookie，该值只能是true，不需要则不设置
		context.Header("Access-Control-Allow-Credentials", "true")
		// 6. 放行所有OPTIONS方法
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}
		context.Next()
	}
}
