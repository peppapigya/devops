package dto

import "mime/multipart"

type ExecutorScript struct {
	Id           int      `json:"id"`
	ScriptId     int64    `json:"scriptId"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Content      string   `json:"content"`
	Parameters   string   `json:"parameters"`
	TimeOut      int      `json:"timeOut"`
	Env          *string  `json:"env"`
	HostIds      []uint32 `json:"hostIds"`
	HostGroupIds []int    `json:"hostGroupIds"`
	WorkDir      string   `json:"workDir"`
}

// DistributeJobScript 分发脚本或本地文件请求类
type DistributeJobScript struct {
	// 脚本ID,如果为空的话则代表是传输的是文件流，不为空则是脚本内容
	Id int64 `form:"id"`
	// 文件
	File       *multipart.FileHeader `form:"-"`
	HostIdsStr string                `form:"targetHosts" validate:"required"`
	// 主机id列表
	HostIds []uint32
	// 上传到远端路径
	RemotePath string `form:"targetPath"`
	// 是否备份
	Backup bool `form:"backup"`
	// 是否覆盖
	Overwrite bool `form:"overwrite"`
	// 文件权限
	Permission string `form:"filePermission"`
	// 传输使用的用户
	User string
}

// DistributeResult 脚本或文件分发结果
type DistributeResult struct {
	HostID   uint32 `json:"hostId"`
	Address  string `json:"address"`
	TaskName string `json:"taskName"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Duration string `json:"duration"`
}
