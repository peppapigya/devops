package dto

type ExecutorScript struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Content      string   `json:"content"`
	Parameters   string   `json:"parameters"`
	TimeOut      int      `json:"timeOut"`
	Env          []string `json:"env"`
	HostIds      []int    `json:"hostIds"`
	HostGroupIds []int    `json:"hostGroupIds"`
	WorkDir      string   `json:"workDir"`
}
