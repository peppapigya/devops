package dto

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
