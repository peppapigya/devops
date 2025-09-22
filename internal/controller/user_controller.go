package controller

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (userController *UserController) GetUserDOById(c *gin.Context) {
	user, err := userController.userService.GetUserById()
	if err != nil {
		common.Fail(c, err)
	}
	common.Success(c, user)
}

// @Tags 用户管理
// @Summary 用户登录
// @Param loginRequest body dto.LoginRequest true "登录请求参数"
// @Success 200 {object}  common.Response{data=dto.CaptchaResponse}
// @Router /sysUser/login [post]
// @security Bearer
func (userController *UserController) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if ok := util.BindAndValidate(c, &loginRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	userController.userService.Login(&loginRequest, c)
}

func (userController *UserController) RefreshToken(c *gin.Context) {
	var refreshRequest dto.RefreshTokenRequest
	if ok := util.BindAndValidate(c, &refreshRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	userController.userService.RefreshToken(c, refreshRequest)
}

func (userController *UserController) GetUserPage(c *gin.Context) {
	var userPageRequest dto.UserPageRequest
	if ok := util.BindAndValidate(c, &userPageRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	userController.userService.GetUserPage(userPageRequest, c)
}

func (userController *UserController) UpdateUserStatus(context *gin.Context) {
	var status string
	// 获取路径参数
	util.GetParam(context, "status", &status, func(param interface{}) {
		if status > "1" || status < "0" {
			common.Fail(context, common.BadRequest)
			context.Abort()
			return
		}
	})
	// 添加校验，如果上下文已被终止则直接返回
	if context.IsAborted() {
		return
	}
	userController.userService.UpdateUserStatus(context, status)
}

func (userController *UserController) UpdateUser(context *gin.Context) {
	//var updateUserRequest dto.UserSaveRequest
}

func (userController *UserController) GetCaptcha(c *gin.Context) {
	captcha, err := userController.userService.GetCaptcha()
	if err != nil {
		common.Fail(c, err)
		return
	}
	common.Success(c, captcha)
}
