package dto

// JobScheduledTaskStatusRequest 定时任务状态请求参数
type JobScheduledTaskStatusRequest struct {
	Id     int64  `json:"id" validate:"required,gt=0"`
	Status uint32 `json:"status" validate:"required,oneof=0 1"`
}
