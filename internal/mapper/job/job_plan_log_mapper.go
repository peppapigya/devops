package job

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/util"

	"gorm.io/gorm"
)

type JobPlanLogMapper struct {
	DB *gorm.DB
}

func NewJobPlanLogMapper(DB *gorm.DB) *JobPlanLogMapper {
	return &JobPlanLogMapper{
		DB: DB,
	}
}

func (m *JobPlanLogMapper) GetJobPlanLogPage(request dto.JobPlanLogPageRequest) (util.PageInfoResponse[model.JobPlanLog], error) {
	var logs []model.JobPlanLog
	var total int64
	db := m.DB.Model(&model.JobPlanLog{})
	if request.PlanID != 0 {
		db = db.Where("plan_id = ?", request.PlanID)
	}
	if request.HostID != 0 {
		db = db.Where("host_id = ?", request.HostID)
	}

	err := db.Count(&total).Error
	if err != nil {
		return util.PageInfoResponse[model.JobPlanLog]{}, err
	}

	offset := (request.PageNum - 1) * request.PageSize
	err = db.Offset(offset).Limit(request.PageSize).Order("created_at desc").Find(&logs).Error
	if err != nil {
		return util.PageInfoResponse[model.JobPlanLog]{}, err
	}

	return util.PageInfoResponse[model.JobPlanLog]{
		Total: total,
		Data:  logs,
	}, nil
}
