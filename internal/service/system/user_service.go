package system

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/convert"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/dal/redis"
	"k8s-platform-go/internal/mapper/system"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

// 用户相关业务

type UserService struct {
	userMapper  *system.UserMapper
	context     *gin.Context
	redisClient *redis.Client
}

func NewUserService(userMapper *system.UserMapper, redisClient *redis.Client) *UserService {
	return &UserService{
		userMapper:  userMapper,
		redisClient: redisClient,
	}
}

func (userService *UserService) GetUserById() (*model.SystemUser, *common.ErrorCode) {
	user, err := userService.userMapper.SelectUserById(1)
	if err != nil {
		return nil, common.ServerError
	}
	if user == nil {
		return nil, common.UserNotExist
	}
	return user, nil
}

func (userService *UserService) Register() {

}

func (userService *UserService) Login(loginRequest *dto.LoginRequest, c *gin.Context) {
	// 校验验证码
	errCode := userService.validateCaptCha(loginRequest.CaptchaId, loginRequest.Code)
	if errCode != nil {
		common.Fail(c, errCode)
		return
	}
	defer func(redisClient *redis.Client, key string) {
		err := redisClient.Delete(key)
		if err != nil {
			log.Printf("删除验证码失败: %v", err)
		}
	}(userService.redisClient, loginRequest.CaptchaId)
	userDO, err := userService.userMapper.SelectUserByName(loginRequest.Username)
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

func (userService *UserService) validateCaptCha(id string, code string) *common.ErrorCode {
	// 校验验证码是否过期
	if userService.redisClient.IsExpired(id) {
		return common.CaptchaNotExist
	}
	// 校验验证码是否正确
	if !util.VerifyCaptcha(userService.redisClient, id, code, false) {
		return common.CaptchaError
	}
	return nil
}

func (userService *UserService) RefreshToken(c *gin.Context, refreshTokenRequest dto.RefreshTokenRequest) {
	token, err := util.ParseToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		common.Fail(c, common.UNAUTHORIZED)
		return
	}
	common.Fail(c, common.UNAUTHORIZED)
	// 根据用户id查看用户信息
	userInfo, err := userService.userMapper.SelectUserById(token.ID)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
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

func (userService *UserService) GetUserPage(request dto.UserPageRequest, c *gin.Context) {
	pageResult, err := userService.userMapper.SelectPageByCondition(request)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, pageResult)
}

func (userService *UserService) UpdateUserStatus(context *gin.Context, status string) {
	userId := util.GetUserIdFromContext(context)
	err := userService.userMapper.UpdateStatusByUserId(userId, status)
	if err != nil {
		log.Printf("更新用户状态失败: %v", err)
		common.Fail(context, common.ServerError)
		return
	}
	common.Success(context, true)
}

func (userService *UserService) GetCaptcha() (dto.CaptchaResponse, *common.ErrorCode) {
	id, code, err := util.GetCaptcha(userService.redisClient)
	if err != nil {
		log.Printf("验证码生成失败: %v", err)
		return dto.CaptchaResponse{}, common.GenerateCaptchaError
	}
	return dto.CaptchaResponse{
		CaptchaId: id,
		Code:      code,
	}, nil
}

func (userService *UserService) DeleteUserById(context *gin.Context, id int64) {
	err := userService.userMapper.DeleteUserById(id)
	if err != nil {
		log.Printf("删除用户失败: %v", err)
		common.Fail(context, common.ServerError)
	}
	common.Success(context, true)
}

func (userService *UserService) AddUser(context *gin.Context, request dto.UserSaveRequest) {
	sysUser := convert.UserDtoToUserDO(&request)
	err := userService.userMapper.InsertUser(sysUser)

	if err != nil {
		log.Printf("添加用户失败: %v", err)
		common.Fail(context, common.ServerError)
	}
	common.Success(context, true)
}

func getAccessTokenAndRefreshToken(c *gin.Context, userDO *model.SystemUser) (string, string, error) {
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
