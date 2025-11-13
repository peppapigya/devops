package dto

// HostPageRequest 主机分页请求参数
type HostPageRequest struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Keyword  string `json:"keyword"`
}

// CreateHostDTO 创建主机请求参数
type CreateHostDTO struct {
	HostName     string `json:"hostName" validate:"required"`
	Address      string `json:"address" validate:"required"`
	HostPort     int32  `json:"hostPort" validate:"required"`
	Username     string `json:"username" validate:"required"`
	HostPassword string `json:"hostPassword" validate:"required"`
	Remark       string `json:"remark"`
}

// UpdateHostDTO 更新主机请求参数
type UpdateHostDTO struct {
	ID           uint32 `json:"id" validate:"required"`
	HostName     string `json:"hostName" validate:"required"`
	Address      string `json:"address" validate:"required"`
	HostPort     int32  `json:"hostPort" validate:"required"`
	Username     string `json:"username" validate:"required"`
	HostPassword string `json:"hostPassword"`
	Remark       string `json:"remark"`
}
