package job

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/util"

	"gorm.io/gorm"
)

type JobPlanMapper struct {
	DB *gorm.DB
}

func NewJobPlanMapper(DB *gorm.DB) *JobPlanMapper {
	return &JobPlanMapper{
		DB: DB,
	}
}

func (m *JobPlanMapper) InsertJobPlan(plan *model.JobPlan) error {
	return m.DB.Create(plan).Error
}

func (m *JobPlanMapper) UpdateJobPlan(plan *model.JobPlan) error {
	return m.DB.Save(plan).Error
}

func (m *JobPlanMapper) DeleteJobPlan(id int64) error {
	return m.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.JobPlan{}, id).Error; err != nil {
			return err
		}
		// Also delete associated script relations
		if err := tx.Where("plan_id = ?", id).Delete(&model.JobPlanScript{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (m *JobPlanMapper) GetJobPlanById(id int64) (*model.JobPlan, error) {
	var plan model.JobPlan
	err := m.DB.First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (m *JobPlanMapper) GetJobPlanPage(request dto.JobPlanPageRequest) (util.PageInfoResponse[model.JobPlan], error) {
	var plans []model.JobPlan
	var total int64
	db := m.DB.Model(&model.JobPlan{})
	if request.Name != "" {
		db = db.Where("name LIKE ?", "%"+request.Name+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return util.PageInfoResponse[model.JobPlan]{}, err
	}

	offset := (request.PageNum - 1) * request.PageSize
	err = db.Offset(offset).Limit(request.PageSize).Find(&plans).Error
	if err != nil {
		return util.PageInfoResponse[model.JobPlan]{}, err
	}

	return util.PageInfoResponse[model.JobPlan]{
		Total: total,
		Data:  plans,
	}, nil
}
