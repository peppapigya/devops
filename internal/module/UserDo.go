package module

import "gorm.io/gorm"

type UserDO struct {
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
	// 角色编码
	RoleCode string `json:"role_code"`
	// 昵称
	DisplayName string `json:"displayName"`
	// 状态 0:禁用 1:正常
	Status int `json:"status"`
	// 手机号
	Phone string `json:"phone"`
	// 邮箱
	Email string `json:"email"`
	// 公共字段
	gorm.Model
}

// TableName 实现gorm提供的TableName方法指定表名
func (user *UserDO) TableName() string {
	return "sys_user"
}

func (user *UserDO) IsEmpty() bool {
	if user == nil {
		return true
	}
	return user.ID == 0
}
