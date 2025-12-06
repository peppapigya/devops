package dto

// JobScriptSaveRequest 保存脚本请求
type JobScriptSaveRequest struct {
	ID         int64  `json:"id"`
	Name       string `json:"name" validate:"required,min=1" label:"脚本名字"`
	Type       string `json:"type" validate:"required"`
	Category   string `json:"category"`
	Content    string `json:"content"`
	Parameters string `json:"parameters"`
	Timeout    int32  `json:"timeout"`
	WorkDir    string `json:"workDir"`
	Env        string `json:"env"`
}

// JobScriptPageRequest 脚本分页请求
type JobScriptPageRequest struct {
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Name     string `json:"name" form:"name"`
	Type     string `json:"type" form:"type"`
}

// JobPlanSaveRequest 保存计划请求
type JobPlanSaveRequest struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name" validate:"required"`
	GlobalVars  string          `json:"globalVars"`
	HostIds     []int64         `json:"hostIds"`
	HostGroupId uint32          `json:"hostGroupId"`
	Remark      string          `json:"remark"`
	Scripts     []JobPlanScript `json:"scripts"` // 关联的脚本ID列表
}

// JobPlanPageRequest 计划分页请求
type JobPlanPageRequest struct {
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Name     string `json:"name" form:"name"`
}

// JobScheduledTaskSaveRequest 保存定时任务请求
type JobScheduledTaskSaveRequest struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" validate:"required"`
	PlanID   uint32 `json:"planId" validate:"required"`
	Strategy string `json:"strategy" validate:"required"`
	Status   int32  `json:"status"`
}

// JobScheduledTaskPageRequest 定时任务分页请求
type JobScheduledTaskPageRequest struct {
	PageNum  int    `json:"pageNum" form:"pageNum"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Name     string `json:"name" form:"name"`
	Status   int32  `json:"status" form:"status"`
}

// JobPlanLogPageRequest 日志分页请求
type JobPlanLogPageRequest struct {
	PageNum  int   `json:"pageNum" form:"pageNum"`
	PageSize int   `json:"pageSize" form:"pageSize"`
	PlanID   int64 `json:"planId" form:"planId"`
	HostID   int64 `json:"hostId" form:"hostId"`
}
