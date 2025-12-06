package job

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/dal/query"
	"k8s-platform-go/internal/util"

	"gorm.io/gorm"
)

type JobScheduledTaskMapper struct {
	DB    *gorm.DB
	query *query.Query
}

func NewJobScheduledTaskMapper(DB *gorm.DB) *JobScheduledTaskMapper {
	return &JobScheduledTaskMapper{
		DB:    DB,
		query: query.Use(DB),
	}
}

func (m *JobScheduledTaskMapper) InsertJobScheduledTask(task *model.JobScheduledTask) error {
	return m.DB.Create(task).Error
}

func (m *JobScheduledTaskMapper) UpdateJobScheduledTask(task *model.JobScheduledTask) error {
	return m.DB.Save(task).Error
}

func (m *JobScheduledTaskMapper) DeleteJobScheduledTask(id int64) error {
	return m.DB.Delete(&model.JobScheduledTask{}, id).Error
}

func (m *JobScheduledTaskMapper) GetJobScheduledTaskById(id int64) (*model.JobScheduledTask, error) {
	var task model.JobScheduledTask
	err := m.DB.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (m *JobScheduledTaskMapper) GetJobScheduledTaskPage(request dto.JobScheduledTaskPageRequest) (util.PageInfoResponse[model.JobScheduledTask], error) {
	var tasks []model.JobScheduledTask
	var total int64
	db := m.DB.Model(&model.JobScheduledTask{})
	if request.Name != "" {
		db = db.Where("name LIKE ?", "%"+request.Name+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return util.PageInfoResponse[model.JobScheduledTask]{}, err
	}

	offset := (request.PageNum - 1) * request.PageSize
	err = db.Offset(offset).Limit(request.PageSize).Find(&tasks).Error
	if err != nil {
		return util.PageInfoResponse[model.JobScheduledTask]{}, err
	}

	return util.PageInfoResponse[model.JobScheduledTask]{
		Total: total,
		Data:  tasks,
	}, nil
}

func (m *JobScheduledTaskMapper) UpdateJobScheduledTaskStatus(id int64, status uint32) (bool, error) {
	jobScheduledTask := m.query.JobScheduledTask
	info, err := jobScheduledTask.Where(jobScheduledTask.ID.Eq(uint32(id))).UpdateColumn(jobScheduledTask.Status, status)
	return info.RowsAffected > 0, err

}
