package dto

import (
	"time"

	"gorm.io/gorm"
)

type JobPlanPageResponse struct {
	ID          uint32                  `json:"id"`
	Name        string                  `json:"name"`
	GlobalVars  string                  `json:"globalVars"`
	HostIds     []string                `json:"hostIds"`
	HostGroupID *uint32                 `json:"hostGroupId"`
	Remark      *string                 `json:"remark"`
	Scripts     []JobPlanScriptResponse `json:"scripts"`
	CreatedAt   time.Time               `json:"createdAt"`
	UpdatedAt   time.Time               `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt          `json:"deletedAt"`
}
