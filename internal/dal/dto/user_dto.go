package dto

import (
	"time"
)

type LoginRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CaptchaId string `json:"captchaId" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

// LoginResponse 登录返回
type LoginResponse struct {
	Id           int64    `json:"id"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	Username     string   `json:"username"`
	Roles        []string `json:"roles"`
}

// RefreshTokenRequest 刷新token请求类
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UserPageRequest struct {
	PageNum  int    `json:"pageNum" default:"1" validator:"min=1"`
	PageSize int    `json:"pageSize" default:"10" validator:"max=100"`
	Username string `json:"userName"`
}

// UserPageResponse 用户分页返回
type UserPageResponse struct {
	Total int64 `json:"total"`
	Data  []interface{}
}

type UserSaveRequest struct {
	ID        int64      `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true;comment:用户ID" json:"id"` // 用户ID
	Username  string     `gorm:"column:username;type:varchar(30);not null;comment:用户账号" json:"username"`         // 用户账号
	Nickname  string     `gorm:"column:nickname;type:varchar(30);not null;comment:用户昵称" json:"nickname"`         // 用户昵称
	Remark    *string    `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`                       // 备注
	DeptID    *int64     `gorm:"column:dept_id;type:bigint(20);comment:部门ID" json:"dept_id"`                     // 部门ID
	PostIds   *string    `gorm:"column:post_ids;type:varchar(255);comment:岗位编号数组" json:"post_ids"`               // 岗位编号数组
	Email     *string    `gorm:"column:email;type:varchar(50);comment:用户邮箱" json:"email"`                        // 用户邮箱
	Mobile    *string    `gorm:"column:mobile;type:varchar(11);comment:手机号码" json:"mobile"`                      // 手机号码
	Sex       *int32     `gorm:"column:sex;type:tinyint(4);comment:用户性别" json:"sex"`                             // 用户性别
	Avatar    *string    `gorm:"column:avatar;type:varchar(512);comment:头像地址" json:"avatar"`                     // 头像地址
	Status    int32      `gorm:"column:status;type:tinyint(4);not null;comment:帐号状态（0正常 1停用）" json:"status"`     // 帐号状态（0正常 1停用）
	LoginIP   *string    `gorm:"column:login_ip;type:varchar(50);comment:最后登录IP" json:"login_ip"`                // 最后登录IP
	LoginDate *time.Time `gorm:"column:login_date;type:datetime;comment:最后登录时间" json:"login_date"`               // 最后登录时间
}

type CaptchaResponse struct {
	CaptchaId string `json:"captchaId"`
	Code      string `json:"code"`
}
