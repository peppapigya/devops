package util

import (
	"k8s-platform-go/internal/common"

	"github.com/gin-gonic/gin"
)

// GetUserInfoFromContext 从上下文中获取用户信息
func GetUserInfoFromContext(c *gin.Context) *Claims {
	claims, _ := c.Get(common.UserInfoKey)
	if claims == nil {
		common.Fail(c, common.UserNotExist)
		return nil
	}
	return claims.(*Claims)
}

// GetUserIdFromContext 从上下文中获取用户ID
func GetUserIdFromContext(c *gin.Context) uint {
	return GetUserInfoFromContext(c).GetUserId()
}
