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
	Id int `json:"id"`
	// 文件
	File *multipart.FileHeader `json:"file"`
	// 上传到远端路径
	RemotePath string `json:"remotePath"`
	// 是否备份，如果为false，则直接覆盖文件
	Backup bool `json:"backup"`
}
