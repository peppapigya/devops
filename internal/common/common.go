package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一响应类

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(
		http.StatusOK,
		Response{
			Code: 200,
			Msg:  "success",
			Data: data,
		},
	)
}

func Fail(c *gin.Context, code *ErrorCode) {
	c.JSON(
		http.StatusOK,
		Response{
			Code: code.Code,
			Msg:  code.Msg,
			Data: nil,
		},
	)
}

func FailWithMsg(c *gin.Context, msg string) {
	c.JSON(
		http.StatusOK,
		Response{
			Code: 500,
			Msg:  msg,
			Data: nil,
		},
	)
}

// PageInfoResponse 分页信息
type PageInfoResponse[T any] struct {
	// 当前页码
	PageNum int `json:"pageNum"`
	// 页面数量
	PageSize int `json:"pageSize"`
	// 数据总数
	Total int64 `json:"total"`
	// 数据
	Data []T `json:"data"`
}
