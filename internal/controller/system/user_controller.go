package system

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service/system"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *system.UserService
}

func NewUserController(userService *system.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Tags 用户管理
// @Summary 获取根据用户id获取用户信息
// @Param id path int true "用户id"
// @Router /sysUser/getUserById [get]
// @security Bearer
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

// @Tags 用户管理
// @Summary 刷新token
// @Param refreshTokenRequest body dto.RefreshTokenRequest true "刷新token请求参数"
// @Router /sysUser/refreshToken [post]
func (userController *UserController) RefreshToken(c *gin.Context) {
	var refreshRequest dto.RefreshTokenRequest
	if ok := util.BindAndValidate(c, &refreshRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	userController.userService.RefreshToken(c, refreshRequest)
}

// @Tags 用户管理
// @Summary 获取用户列表分页
// @Param userPageRequest body dto.UserPageRequest true "用户列表分页请求参数"
// @Router /sysUser/getUserPage [post]
func (userController *UserController) GetUserPage(c *gin.Context) {
	var userPageRequest dto.UserPageRequest
	if ok := util.BindAndValidate(c, &userPageRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	userController.userService.GetUserPage(userPageRequest, c)
}

// @Tags 用户管理
// @Summary 更新用户状态
// @Param status query string true "用户状态"
// @Router /sysUser/updateUserStatus [put]
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

// @Tags 用户管理
// @Summary 更新用户信息
// @Param updateUserRequest body dto.UserSaveRequest true "更新用户请求参数"
// @Success 200 {object} common.Response
// @Router /sysUser/ [put]
// @Security Bearer
func (userController *UserController) UpdateUser(context *gin.Context) {
	//var updateUserRequest dto.UserSaveRequest
}

// @Tags 用户管理
// @Summary 删除用户
// @Param id path int true "用户id"
// @Router /sysUser/deleteUserById [delete]
// @Success 200 {object} common.Response
// @Security Bearer
func (userController *UserController) DeleteUserById(context *gin.Context) {
	var id int64
	util.GetParam(context, "id", &id, func(param interface{}) {
		if id <= 0 {
			common.Fail(context, common.BadRequest)
			context.Abort()
			return
		}
	})
	userController.userService.DeleteUserById(context, id)
}

// @Tags 用户管理
// @Summary 添加用户
// @Param registerRequest body dto.UserSaveRequest true "添加用户请求参数"
// @Success 200 {object} common.Response
// @Router /sysUser/register [post]
// @Security Bearer
func (userController *UserController) CreateUser(context *gin.Context) {
	var userSaveRequest dto.UserSaveRequest
	if ok := util.BindAndValidate(context, &userSaveRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		common.Fail(context, common.BadRequest)
		context.Abort()
		return
	}
	userController.userService.AddUser(context, userSaveRequest)
}

// @Tags 用户管理
// @Summary 获取验证码
// @Router /sysUser/getCaptcha [get]
func (userController *UserController) GetCaptcha(c *gin.Context) {
	captcha, err := userController.userService.GetCaptcha()
	if err != nil {
		common.Fail(c, err)
		return
	}
	common.Success(c, captcha)
}
