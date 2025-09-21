package mapper

import (
	"k8s-platform-go/internal/module"
	"k8s-platform-go/internal/module/dto"

	"gorm.io/gorm"
)

type UserMapper struct {
	DB *gorm.DB
}

func NewUserMapper(DB *gorm.DB) *UserMapper {
	return &UserMapper{DB: DB}
}

// SelectUserById 通过用户id获取用户信息
func (userMapper *UserMapper) SelectUserById(id int64) *module.UserDO {
	var userDO module.UserDO
	userMapper.DB.Where("id = ?", id).First(&userDO)
	return &userDO
}

func (userMapper *UserMapper) SelectUserByName(username string) *module.UserDO {
	var userDO module.UserDO
	userMapper.DB.Select("id", "username", "password").Where("username = ?", username).Take(&userDO)
	return &userDO
}

// SelectPageByCondition 分页查询用户信息
func (userMapper *UserMapper) SelectPageByCondition(request dto.UserPageRequest) ([]module.UserDO, int64, error) {

	var (
		userDO []module.UserDO
		total  int64
	)
	query := userMapper.DB.Model(&userDO)
	if request.Username != "" {
		query.Where("username like ?", "%"+request.Username+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (request.PageNum - 1) * request.PageSize
	if ok := query.Offset(offset).Find(&userDO).Error; ok != nil {
		return nil, 0, ok
	}
	return userDO, total, nil
}

// UpdateStatusByUserId 更新用户状态
func (userMapper *UserMapper) UpdateStatusByUserId(id uint, status string) {
	userMapper.getQuery().Where("id = ?", id).Update("status", status)
}

func (userMapper *UserMapper) getQuery() *gorm.DB {
	return userMapper.DB.Model(&module.UserDO{})
}
