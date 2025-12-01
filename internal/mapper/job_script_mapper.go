package mapper

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/dal/query"
	"k8s-platform-go/internal/util"

	"gorm.io/gorm"
)

type JobScriptMapper struct {
	DB    *gorm.DB
	query *query.Query
}

func NewJobScriptMapper(DB *gorm.DB) *JobScriptMapper {
	return &JobScriptMapper{
		DB:    DB,
		query: query.Use(DB),
	}
}

func (m *JobScriptMapper) InsertJobScript(script *model.JobScript) error {
	return m.DB.Create(script).Error
}

func (m *JobScriptMapper) UpdateJobScript(script *model.JobScript) error {
	return m.DB.Save(script).Error
}

func (m *JobScriptMapper) DeleteJobScript(id int64) error {
	return m.DB.Delete(&model.JobScript{}, id).Error
}

func (m *JobScriptMapper) GetJobScriptById(id int64) (*model.JobScript, error) {
	var script model.JobScript
	err := m.DB.First(&script, id).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (m *JobScriptMapper) GetJobScriptPage(request dto.JobScriptPageRequest) (util.PageInfoResponse[model.JobScript], error) {
	var scripts []model.JobScript
	var total int64
	db := m.DB.Model(&model.JobScript{})
	if request.Name != "" {
		db = db.Where("name LIKE ?", "%"+request.Name+"%")
	}
	if request.Type != "" {
		db = db.Where("type = ?", request.Type)
	}

	err := db.Count(&total).Error
	if err != nil {
		return util.PageInfoResponse[model.JobScript]{}, err
	}

	offset := (request.PageNum - 1) * request.PageSize
	err = db.Offset(offset).Limit(request.PageSize).Find(&scripts).Error
	if err != nil {
		return util.PageInfoResponse[model.JobScript]{}, err
	}

	return util.PageInfoResponse[model.JobScript]{
		Total:    total,
		Data:     scripts,
		PageNum:  request.PageNum,
		PageSize: request.PageSize,
	}, nil
}

// SelectListByCondition 根据脚本的名字查询脚本
func (m *JobScriptMapper) SelectListByCondition(condition string) ([]*model.JobScript, error) {
	script := m.query.JobScript
	if condition != "" {
		return script.Where(script.Name.Like("%" + condition + "%")).Find()
	}
	return script.Find()
}
