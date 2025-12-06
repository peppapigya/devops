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

type DictTypeMapper struct {
	DB    *gorm.DB
	query *query.Query
}

func NewDictTypeMapper(DB *gorm.DB) *DictTypeMapper {
	return &DictTypeMapper{DB: DB, query: query.Use(DB)}
}

func (m *DictTypeMapper) base() query.ISystemDictTypeDo {
	return m.query.SystemDictType.WithContext(context.Background())
}

func (m *DictTypeMapper) SelectByID(id int64) (*model.SystemDictType, error) {
	d := m.query.SystemDictType
	return m.base().Where(d.ID.Eq(id)).First()
}

func (m *DictTypeMapper) SelectPage(req dto.DictTypePageRequest) (util.PageInfoResponse[model.SystemDictType], error) {
	d := m.query.SystemDictType
	var conds []gen.Condition
	conds = append(conds, util.WhereIf(req.Name != "", d.Name.Like("%"+req.Name+"%")))
	conds = append(conds, util.WhereIf(req.Type != "", d.Type.Like("%"+req.Type+"%")))
	if req.Status != "" {
		if v, err := strconv.Atoi(req.Status); err == nil {
			conds = append(conds, d.Status.Eq(int32(v)))
		}
	}
	return util.FindPageResult[model.SystemDictType](m.query.SystemDictType.DO, req.PageNum, req.PageSize, false, conds...)
}

func (m *DictTypeMapper) Insert(v *model.SystemDictType) error { return m.base().Create(v) }
func (m *DictTypeMapper) Update(v *model.SystemDictType) error {
	_, err := m.base().Where(m.query.SystemDictType.ID.Eq(v.ID)).Updates(v)
	return err
}
func (m *DictTypeMapper) DeleteByID(id int64) error {
	_, err := m.base().Where(m.query.SystemDictType.ID.Eq(id)).Delete()
	return err
}
func (m *DictTypeMapper) UpdateStatus(id int64, status int32) error {
	_, err := m.base().Where(m.query.SystemDictType.ID.Eq(id)).Update(m.query.SystemDictType.Status, status)
	return err
}
