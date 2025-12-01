package mapper

import (
	"k8s-platform-go/internal/dal/model"

	"gorm.io/gorm"
)

type JobPlanScriptMapper struct {
	DB *gorm.DB
}

func NewJobPlanScriptMapper(DB *gorm.DB) *JobPlanScriptMapper {
	return &JobPlanScriptMapper{
		DB: DB,
	}
}

func (m *JobPlanScriptMapper) DeleteJobPlanScripts(planId int64) error {
	return m.DB.Where("plan_id = ?", planId).Delete(&model.JobPlanScript{}).Error
}

func (m *JobPlanScriptMapper) InsertJobPlanScripts(scripts []model.JobPlanScript) error {
	if len(scripts) == 0 {
		return nil
	}
	return m.DB.Create(&scripts).Error
}

func (m *JobPlanScriptMapper) GetJobPlanScripts(planId int64) ([]model.JobPlanScript, error) {
	var scripts []model.JobPlanScript
	err := m.DB.Where("plan_id = ?", planId).Order("sort asc").Find(&scripts).Error
	return scripts, err
}
