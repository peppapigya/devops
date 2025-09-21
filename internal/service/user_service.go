package service

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/mapper"
	"k8s-platform-go/internal/module"
	"k8s-platform-go/internal/module/dto"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	userMapper *mapper.UserMapper
}

func NewUserService(userMapper *mapper.UserMapper) *UserService {
	return &UserService{userMapper: userMapper}
}

func (userService *UserService) GetUserById() *module.UserDO {
	return userService.userMapper.SelectUserById(1)
}

func (userService *UserService) Register() {

}

func (userService *UserService) Login(loginRequest *dto.LoginRequest, c *gin.Context) {
	userDO := userService.userMapper.SelectUserByName(loginRequest.Username)
	if userDO == nil {
		common.Fail(c, common.UserNotExist)
		return
	}
	if userDO.Password != loginRequest.Password {
		common.Fail(c, common.UserPasswordError)
		return
	}
	// 生成token信息
	accessToken, refreshToken, err := getAccessTokenAndRefreshToken(c, userDO)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	// 封装返回信息
	loginResponse := dto.LoginResponse{
		Id:           userDO.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     userDO.Username,
		Roles:        nil,
	}
	common.Success(c, loginResponse)
}

// RefreshToken 刷新token
func (userService *UserService) RefreshToken(c *gin.Context, refreshTokenRequest dto.RefreshTokenRequest) {
	token, err := util.ParseToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		common.Fail(c, common.UNAUTHORIZED)
		return
	}
	common.Fail(c, common.UNAUTHORIZED)
	// 根据用户id查看用户信息
	userInfo := userService.userMapper.SelectUserById(int64(token.ID))
	// 重新生成token信息
	accessToken, refreshToken, err := getAccessTokenAndRefreshToken(c, userInfo)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	// 返回给前端
	common.Success(c, dto.LoginResponse{
		Id:           userInfo.ID,
		Username:     userInfo.Username,
		Roles:        nil,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// GetUserPage 获取用户分页信息
func (userService *UserService) GetUserPage(request dto.UserPageRequest, c *gin.Context) {
	userDos, total, err := userService.userMapper.SelectPageByCondition(request)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, common.PageInfoResponse{
		Data:     userDos,
		PageNum:  request.PageNum,
		PageSize: request.PageSize,
		Total:    total,
	})
}

// UpdateUserStatus 更新用户状态
func (userService *UserService) UpdateUserStatus(context *gin.Context, status string) {
	userId := util.GetUserIdFromContext(context)
	userService.userMapper.UpdateStatusByUserId(userId, status)
	common.Success(context, true)
}

func getAccessTokenAndRefreshToken(c *gin.Context, userDO *module.UserDO) (string, string, error) {
	accessToken, err := util.GenerateJwtToken(userDO.ID, userDO.Username, []string{})
	if err != nil {
		log.Printf("jwt 生成失败: %v", err)
		common.Fail(c, common.ServerError)
		return "", "", err
	}
	refreshToken, err := util.GenerateRefreshToken(userDO.ID)
	if err != nil {
		log.Printf("refresh token 生成失败: %v", err)
		common.Fail(c, common.ServerError)
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
