package dto

// HostPageRequest 主机分页请求参数
type HostPageRequest struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Keyword  string `json:"keyword"`
}

// CreateHostDTO 创建主机请求参数
type CreateHostDTO struct {
	HostName     string `json:"host_name" validate:"required"`
	Address      string `json:"address" validate:"required"`
	HostPort     int32  `json:"host_port" validate:"required"`
	Username     string `json:"username" validate:"required"`
	HostPassword string `json:"host_password" validate:"required"`
	Remark       string `json:"remark"`
}

// UpdateHostDTO 更新主机请求参数
type UpdateHostDTO struct {
	HostName     string `json:"host_name" validate:"required"`
	Address      string `json:"address" validate:"required"`
	HostPort     int32  `json:"host_port" validate:"required"`
	Username     string `json:"username" validate:"required"`
	HostPassword string `json:"host_password"`
	Remark       string `json:"remark"`
}
