package util

import (
	"errors"
	"k8s-platform-go/internal/common"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindAndValidate 解析参数并将参数绑定到obj
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var errs validator.ValidationErrors
		if ok := errors.As(err, &errs); ok {
			log.Printf("参数校验失败: %v", errs)
			common.FailWithMsg(c, errs[0].Translate(trans))
			return false
		}
		log.Printf("解析参数失败: %v", err)
		common.Fail(c, common.ServerError)
		return false
	}
	return true
}

// GetParam 获取路径参数以及参数校验
func GetParam(c *gin.Context, key string, param interface{}, validate func(param interface{})) {
	var value string
	value = c.Query(key)
	if value == "" {
		value = c.Param(key)
	}
	if strParam, ok := param.(*string); ok {
		*strParam = value
	}
	if validate != nil {
		validate(param)
	}
	return
}
