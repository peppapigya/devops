package system

import (
	"context"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/dal/query"
	"k8s-platform-go/internal/util"
	"strconv"

	"gorm.io/gen"
	"gorm.io/gorm"
)

type DeptMapper struct {
	DB    *gorm.DB
	query *query.Query
}

func NewDeptMapper(DB *gorm.DB) *DeptMapper {
	return &DeptMapper{DB: DB, query: query.Use(DB)}
}

func (m *DeptMapper) GetBaseMapper() query.ISystemDeptDo {
	return m.query.SystemDept.WithContext(context.Background())
}

func (m *DeptMapper) SelectByID(id int64) (*model.SystemDept, error) {
	d := m.query.SystemDept
	return m.GetBaseMapper().Where(d.ID.Eq(id)).First()
}

func (m *DeptMapper) SelectPageByCondition(req dto.DeptPageRequest) (util.PageInfoResponse[model.SystemDept], error) {
	d := m.query.SystemDept
	var conds []gen.Condition
	conds = append(conds, util.WhereIf(req.Name != "", d.Name.Like("%"+req.Name+"%")))
	if req.Status != "" {
		if v, err := strconv.Atoi(req.Status); err == nil {
			conds = append(conds, d.Status.Eq(int32(v)))
		}
	}
	return util.FindPageResult[model.SystemDept](m.query.SystemDept.DO, req.PageNum, req.PageSize, false, conds...)
}

func (m *DeptMapper) Insert(dept *model.SystemDept) error {
	return m.GetBaseMapper().Create(dept)
}

func (m *DeptMapper) Update(dept *model.SystemDept) error {
	_, err := m.GetBaseMapper().Where(m.query.SystemDept.ID.Eq(dept.ID)).Updates(dept)
	return err
}

func (m *DeptMapper) DeleteByID(id int64) error {
	_, err := m.GetBaseMapper().Where(m.query.SystemDept.ID.Eq(id)).Delete()
	return err
}

func (m *DeptMapper) UpdateStatus(id int64, status int32) error {
	_, err := m.GetBaseMapper().Where(m.query.SystemDept.ID.Eq(id)).Update(m.query.SystemDept.Status, status)
	return err
}

func (m *DeptMapper) ListAllEnabled() ([]*model.SystemDept, error) {
	d := m.query.SystemDept
	return m.GetBaseMapper().Where(d.Status.Eq(1)).Order(d.Sort).Find()
}
