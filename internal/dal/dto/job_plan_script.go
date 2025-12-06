package dto

type JobPlanScript struct {
	ScriptId uint32 `json:"scriptId"`
	Sort     uint32 `json:"sort"`
	Name     string `json:"scriptName"`
}

type JobPlanScriptResponse struct {
	ID         uint32 `json:"id"`
	PlanId     uint32 `json:"planId"`
	ScriptId   uint32 `json:"scriptId"`
	Sort       uint32 `json:"sort"`
	ScriptName string `json:"scriptName"`
}
