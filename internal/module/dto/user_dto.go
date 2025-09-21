package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录返回
type LoginResponse struct {
	Id           uint     `json:"id"`
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
