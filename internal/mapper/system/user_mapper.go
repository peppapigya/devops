package system

import (
	"context"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/dal/query"
	"k8s-platform-go/internal/util"

	"gorm.io/gorm"
)

type UserMapper struct {
	DB    *gorm.DB
	query *query.Query
}

func NewUserMapper(DB *gorm.DB) *UserMapper {
	return &UserMapper{
		DB:    DB,
		query: query.Use(DB),
	}
}

// SelectUserById 通过用户id获取用户信息
func (userMapper *UserMapper) SelectUserById(id int64) (*model.SystemUser, error) {
	user := userMapper.query.SystemUser
	return query.SystemUser.WithContext(context.Background()).
		Where(user.ID.Eq(id)).First()
}

// SelectUserByName 通过用户名获取用户信息
func (userMapper *UserMapper) SelectUserByName(username string) (*model.SystemUser, error) {
	user := userMapper.query.SystemUser
	return userMapper.query.SystemUser.WithContext(context.Background()).
		Select(user.ID, user.Username, user.Password).
		Where(user.Username.Eq(username)).First()
}

// SelectPageByCondition 分页查询用户信息
func (userMapper *UserMapper) SelectPageByCondition(request dto.UserPageRequest) (util.PageInfoResponse[model.SystemUser], error) {

	user := userMapper.query.SystemUser
	return util.FindPageResult[model.SystemUser](userMapper.query.SystemUser.DO, request.PageNum, request.PageSize, false,
		util.WhereIf(request.Username != "", user.Username.Like("%"+request.Username+"%")),
	)
}

// UpdateStatusByUserId 更新用户状态
func (userMapper *UserMapper) UpdateStatusByUserId(id int64, status string) error {
	user := userMapper.query.SystemUser
	_, err := userMapper.GetBaseMapper().Where(user.ID.Eq(id)).Update(user.Status, status)
	return err
}

func (userMapper *UserMapper) GetBaseMapper() query.ISystemUserDo {
	return userMapper.query.SystemUser.WithContext(context.Background())
}

// DeleteUserById 删除用户
func (userMapper *UserMapper) DeleteUserById(id int64) error {
	user := userMapper.query.SystemUser
	_, err := userMapper.GetBaseMapper().WithContext(context.Background()).Where(user.ID.Eq(id)).Delete()
	return err
}

// InsertUser 添加用户信息到数据库
func (userMapper *UserMapper) InsertUser(sysUser *model.SystemUser) error {
	err := userMapper.GetBaseMapper().WithContext(context.Background()).Create(sysUser)
	return err
}
