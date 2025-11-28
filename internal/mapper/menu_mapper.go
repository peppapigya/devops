package mapper

import (
	"context"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/dal/query"

	"gorm.io/gorm"
)

type MenuMapper struct {
	DB *gorm.DB
	q  *query.Query
}

func NewMenuMapper(DB *gorm.DB) *MenuMapper {
	return &MenuMapper{DB: DB, q: query.Use(DB)}
}

func (m *MenuMapper) base() query.ISystemMenuDo {
	return m.q.SystemMenu.WithContext(context.Background())
}

// ListAllEnabledVisible 获取所有启用的可见菜单
func (m *MenuMapper) ListAllEnabledVisible() ([]*model.SystemMenu, error) {
	sm := m.q.SystemMenu
	return m.base().Where(sm.Status.Eq(1)).Find()
}

// InsertAll 批量插入菜单
func (m *MenuMapper) InsertAll(list []*model.SystemMenu) error {
	return m.base().CreateInBatches(list, 100)
}

// GetByID 根据 ID 获取菜单
func (m *MenuMapper) GetByID(id int64) (*model.SystemMenu, error) {
	sm := m.q.SystemMenu
	return m.base().Where(sm.ID.Eq(id)).First()
}

// Create 创建菜单
func (m *MenuMapper) Create(menu *model.SystemMenu) error {
	return m.base().Create(menu)
}

// Update 更新菜单
func (m *MenuMapper) Update(menu *model.SystemMenu) error {
	return m.base().Save(menu)
}

// Delete 删除菜单（软删除）
func (m *MenuMapper) Delete(id int64) error {
	sm := m.q.SystemMenu
	_, err := m.base().Where(sm.ID.Eq(id)).Delete(&model.SystemMenu{})
	return err

}

// GetMenuList 获取菜单列表（分页）
func (m *MenuMapper) GetMenuList(queryDTO *dto.MenuQueryDTO) ([]*model.SystemMenu, int64, error) {
	sm := m.q.SystemMenu
	q := m.base()

	// 构建查询条件
	if queryDTO.Name != "" {
		q = q.Where(sm.Name.Like("%" + queryDTO.Name + "%"))
	}
	if queryDTO.Type != 0 {
		q = q.Where(sm.Type.Eq(queryDTO.Type))
	}
	if queryDTO.Status != 0 {
		q = q.Where(sm.Status.Eq(queryDTO.Status))
	}
	if queryDTO.ParentID != 0 {
		q = q.Where(sm.ParentID.Eq(queryDTO.ParentID))
	}

	// 获取总数
	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (queryDTO.Page - 1) * queryDTO.PageSize
	list, err := q.Order(sm.Sort.Asc()).Offset(offset).Limit(queryDTO.PageSize).Find()
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetMenuTree 获取菜单树形结构
func (m *MenuMapper) GetMenuTree() ([]*model.SystemMenu, error) {
	sm := m.q.SystemMenu
	return m.base().Where(sm.Status.Eq(1)).Order(sm.Sort.Asc()).Find()
}

// GetMenuOptions 获取菜单选项列表（用于下拉选择）
func (m *MenuMapper) GetMenuOptions() ([]*model.SystemMenu, error) {
	sm := m.q.SystemMenu
	return m.base().Where(sm.Status.Eq(1)).Order(sm.Sort.Asc()).Find()
}

// CheckMenuNameExists 检查菜单名称是否存在
func (m *MenuMapper) CheckMenuNameExists(name string, parentID int64, excludeID ...int64) (bool, error) {
	sm := m.q.SystemMenu
	q := m.base().Where(sm.Name.Eq(name), sm.ParentID.Eq(parentID))

	if len(excludeID) > 0 {
		q = q.Where(sm.ID.Neq(excludeID[0]))
	}

	count, err := q.Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetChildrenMenus 获取子菜单
func (m *MenuMapper) GetChildrenMenus(parentID int64) ([]*model.SystemMenu, error) {
	sm := m.q.SystemMenu
	return m.base().Where(sm.ParentID.Eq(parentID)).Order(sm.Sort.Asc()).Find()
}

// GetMaxSort 获取指定父菜单下的最大排序号
func (m *MenuMapper) GetMaxSort(parentID int64) (int32, error) {
	sm := m.q.SystemMenu
	result, err := m.base().Where(sm.ParentID.Eq(parentID)).Order(sm.Sort.Desc()).First()
	if err != nil {
		if err.Error() == "record not found" {
			return 0, nil
		}
		return 0, err
	}

	return result.Sort, nil
}

// BatchDelete 批量删除菜单
func (m *MenuMapper) BatchDelete(ids []int64) error {
	sm := m.q.SystemMenu
	_, err := m.base().Where(sm.ID.In(ids...)).Delete(&model.SystemMenu{})
	return err
}
