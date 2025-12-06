package system

import (
	"fmt"
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper/system"
	"k8s-platform-go/internal/util"
	"time"

	"github.com/gin-gonic/gin"
)

type MenuService struct{ m *system.MenuMapper }

func NewMenuService(m *system.MenuMapper) *MenuService { return &MenuService{m: m} }

type MenuRoute struct {
	Name          string      `json:"name"`
	Path          string      `json:"path"`
	Component     string      `json:"component,omitempty"`
	ComponentName string      `json:"componentName,omitempty"`
	Icon          string      `json:"icon,omitempty"`
	Permission    string      `json:"permission,omitempty"`
	KeepAlive     bool        `json:"keepAlive,omitempty"`
	Children      []MenuRoute `json:"children,omitempty"`
}

// Routes 获取菜单路由
func (s *MenuService) Routes(c *gin.Context) {
	list, err := s.m.ListAllEnabledVisible()
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	mp := map[int64]*MenuRoute{}
	var roots []*MenuRoute
	for _, m := range list {
		r := &MenuRoute{
			Name:          m.Name,
			Path:          getStr(m.Path),
			Component:     getStr(m.Component),
			ComponentName: getStr(m.ComponentName),
			Icon:          getStr(m.Icon),
			Permission:    m.Permission,
			KeepAlive:     bitToBool(m.KeepAlive),
		}
		mp[m.ID] = r
	}
	for _, m := range list {
		r := mp[m.ID]
		if m.ParentID == 0 {
			roots = append(roots, r)
		} else if p, ok := mp[m.ParentID]; ok {
			p.Children = append(p.Children, *r)
		} else {
			roots = append(roots, r)
		}
	}
	var out []MenuRoute
	for _, r := range roots {
		out = append(out, *r)
	}
	common.Success(c, out)
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(c *gin.Context, createDTO *dto.MenuCreateDTO) error {
	// 校验菜单名称是否已存在
	exists, err := s.m.CheckMenuNameExists(createDTO.Name, createDTO.ParentID)
	if err != nil {
		common.Fail(c, common.MenuNameCheckFailed)
		return err
	}
	if exists {
		common.Fail(c, common.MenuNameExist)
		return err
	}

	// 获取最大排序号
	maxSort, err := s.m.GetMaxSort(createDTO.ParentID)
	if err != nil {
		common.Fail(c, common.ServerError)
		return err
	}

	// 创建菜单
	menu := &model.SystemMenu{
		Name:          createDTO.Name,
		Permission:    createDTO.Permission,
		Type:          createDTO.Type,
		Sort:          maxSort + 1,
		ParentID:      createDTO.ParentID,
		Path:          createDTO.Path,
		Icon:          createDTO.Icon,
		Component:     createDTO.Component,
		ComponentName: createDTO.ComponentName,
		Status:        1,
		Visible:       toBit(createDTO.Visible),
		KeepAlive:     toBit(createDTO.KeepAlive),
		AlwaysShow:    toBit(createDTO.AlwaysShow),
		Creator:       &util.GetUserInfoFromContext(c).Username,
	}

	err = s.m.Create(menu)
	if err != nil {
		common.Fail(c, common.ServerError)
		return err
	}

	common.Success(c, menu)
	return nil
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(c *gin.Context, updateDTO *dto.MenuUpdateDTO) error {
	// 检查菜单是否存在
	existing, err := s.m.GetByID(updateDTO.ID)
	if err != nil {
		common.Fail(c, common.ServerError)
		return nil
	}

	// 校验菜单名称是否已存在（排除当前菜单）
	exists, err := s.m.CheckMenuNameExists(updateDTO.Name, updateDTO.ParentID, updateDTO.ID)
	if err != nil {
		common.Fail(c, common.MenuNameCheckFailed)
		return err
	}
	if exists {
		common.Fail(c, common.MenuNameExist)
	}

	// 检查是否会形成循环依赖
	if err := s.checkCircularDependency(updateDTO.ID, updateDTO.ParentID); err != nil {
		common.Fail(c, common.MenuNameExist)
		return err
	}

	// 更新菜单
	existing.Name = updateDTO.Name
	existing.Permission = updateDTO.Permission
	existing.Type = updateDTO.Type
	existing.Sort = updateDTO.Sort
	existing.ParentID = updateDTO.ParentID
	existing.Path = updateDTO.Path
	existing.Icon = updateDTO.Icon
	existing.Component = updateDTO.Component
	existing.ComponentName = updateDTO.ComponentName
	existing.Status = updateDTO.Status
	existing.Visible = toBit(updateDTO.Visible)
	existing.KeepAlive = toBit(updateDTO.KeepAlive)
	existing.AlwaysShow = toBit(updateDTO.AlwaysShow)
	existing.Updater = &util.GetUserInfoFromContext(c).Username
	existing.UpdateAt = time.Now()

	err = s.m.Update(existing)
	if err != nil {
		common.Fail(c, common.ServerError)
		return err
	}

	common.Success(c, existing)
	return nil
}

// DeleteMenu 删除菜单
func (s *MenuService) DeleteMenu(c *gin.Context, id int64) error {
	// 检查菜单是否存在
	_, err := s.m.GetByID(id)
	if err != nil {
		common.Fail(c, common.MenuNotExist)
		return err
	}

	// 检查是否有子菜单
	children, err := s.m.GetChildrenMenus(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return nil
	}
	if len(children) > 0 {
		common.Fail(c, common.MenuHasChildren)
		return err
	}

	err = s.m.Delete(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return err
	}

	common.Success(c, nil)
	return nil
}

// GetMenuList 获取菜单列表
func (s *MenuService) GetMenuList(c *gin.Context, queryDTO *dto.MenuQueryDTO) error {
	// 设置默认值
	if queryDTO.Page == 0 {
		queryDTO.Page = 1
	}
	if queryDTO.PageSize == 0 {
		queryDTO.PageSize = 10
	}

	list, total, err := s.m.GetMenuList(queryDTO)
	if err != nil {
		common.Fail(c, common.ServerError)
		return err
	}

	// 转换为VO
	menuList := make([]*dto.MenuListVO, 0, len(list))
	for _, menu := range list {
		menuList = append(menuList, s.convertToMenuListVO(menu))
	}

	result := &dto.PageResult{
		Total: total,
		List:  menuList,
	}

	common.Success(c, result)
	return nil
}

// GetMenuTree 获取菜单树
func (s *MenuService) GetMenuTree(c *gin.Context) error {
	list, err := s.m.GetMenuTree()
	if err != nil {
		common.Fail(c, common.ServerError)
		return nil
	}

	// 构建树形结构
	tree := s.buildMenuTree(list, 0)
	common.Success(c, tree)
	return nil
}

// GetMenuOptions 获取菜单选项
func (s *MenuService) GetMenuOptions(c *gin.Context) error {
	options, err := s.m.GetMenuOptions()
	if err != nil {
		common.Fail(c, common.ServerError)
		return err
	}
	res := make([]*dto.MenuOptionVO, len(options))
	for _, option := range options {
		res = append(res, &dto.MenuOptionVO{
			ID:   option.ID,
			Name: option.Name,
		})
	}
	common.Success(c, options)
	return nil
}

// GetMenuById 根据ID获取菜单详情
func (s *MenuService) GetMenuById(c *gin.Context, id int64) error {
	menu, err := s.m.GetByID(id)
	if err != nil {
		common.Fail(c, common.MenuNotExist)
		return err
	}

	// 转换为VO
	menuVO := s.convertToMenuListVO(menu)
	common.Success(c, menuVO)
	return nil
}

// SeedDefault 初始化默认菜单
func (s *MenuService) SeedDefault() {
	list, _ := s.m.ListAllEnabledVisible()
	if len(list) > 0 {
		return
	}
	menus := []*model.SystemMenu{
		{Name: "首页", Permission: "home:view", Type: 1, Sort: 1, ParentID: 0, Path: stringPtr("/home"), Icon: stringPtr("House"), Component: stringPtr("home/index"), ComponentName: stringPtr("home"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},
		{Name: "系统管理", Permission: "system:menu:list", Type: 1, Sort: 2, ParentID: 0, Path: stringPtr("/system"), Icon: stringPtr("Setting"), Component: stringPtr("system/index"), ComponentName: stringPtr("system"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},
		{Name: "字典管理", Permission: "system:dict:list", Type: 1, Sort: 3, ParentID: 0, Path: stringPtr("/dict"), Icon: stringPtr("CollectionTag"), Component: stringPtr("dict/index"), ComponentName: stringPtr("dict"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},
		{Name: "监控运维", Permission: "monitor:view", Type: 1, Sort: 4, ParentID: 0, Path: stringPtr("/monitor"), Icon: stringPtr("Monitor"), Component: stringPtr("monitor/index"), ComponentName: stringPtr("monitor"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},

		// 二级菜单 - 系统管理
		{Name: "用户管理", Permission: "system:user:list", Type: 2, Sort: 1, ParentID: 2, Path: stringPtr("/user"), Icon: stringPtr("User"), Component: stringPtr("user/index"), ComponentName: stringPtr("user"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},
		{Name: "部门管理", Permission: "system:dept:list", Type: 2, Sort: 2, ParentID: 2, Path: stringPtr("/dept"), Icon: stringPtr("Collection"), Component: stringPtr("dept/index"), ComponentName: stringPtr("dept"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},
		{Name: "角色管理", Permission: "system:role:list", Type: 2, Sort: 3, ParentID: 2, Path: stringPtr("/role"), Icon: stringPtr("Avatar"), Component: stringPtr("role/index"), ComponentName: stringPtr("role"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},
		{Name: "菜单管理", Permission: "system:menu:list", Type: 2, Sort: 4, ParentID: 2, Path: stringPtr("/menu"), Icon: stringPtr("List"), Component: stringPtr("system/menu/index"), ComponentName: stringPtr("menu"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},

		// 二级菜单 - 字典管理
		{Name: "字典类型", Permission: "system:dictType:list", Type: 2, Sort: 1, ParentID: 3, Path: stringPtr("/dict/type"), Icon: stringPtr("CollectionTag"), Component: stringPtr("dict/type/index"), ComponentName: stringPtr("dict-type"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},
		{Name: "字典数据", Permission: "system:dictData:list", Type: 2, Sort: 2, ParentID: 3, Path: stringPtr("/dict/data"), Icon: stringPtr("List"), Component: stringPtr("dict/data/index"), ComponentName: stringPtr("dict-data"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},

		// 二级菜单 - 监控运维
		{Name: "主机管理", Permission: "host:list", Type: 2, Sort: 1, ParentID: 4, Path: stringPtr("/hosts"), Icon: stringPtr("Monitor"), Component: stringPtr("hosts/index"), ComponentName: stringPtr("hosts"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},
		{Name: "AI助手", Permission: "ai:view", Type: 2, Sort: 2, ParentID: 4, Path: stringPtr("/ai"), Icon: stringPtr("MagicStick"), Component: stringPtr("ai/index"), ComponentName: stringPtr("ai"), Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(true), AlwaysShow: boolPtr(true), Creator: stringPtr("system")},

		// 按钮权限
		{Name: "新增用户", Permission: "system:user:add", Type: 3, Sort: 1, ParentID: 5, Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(false), AlwaysShow: boolPtr(false), Creator: stringPtr("system")},
		{Name: "编辑用户", Permission: "system:user:edit", Type: 3, Sort: 2, ParentID: 5, Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(false), AlwaysShow: boolPtr(false), Creator: stringPtr("system")},
		{Name: "删除用户", Permission: "system:user:delete", Type: 3, Sort: 3, ParentID: 5, Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(false), AlwaysShow: boolPtr(false), Creator: stringPtr("system")},
		{Name: "导出用户", Permission: "system:user:export", Type: 3, Sort: 4, ParentID: 5, Status: 1, Visible: boolPtr(true), KeepAlive: boolPtr(false), AlwaysShow: boolPtr(false), Creator: stringPtr("system")},
	}

	_ = s.m.InsertAll(menus)
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) []uint8 {
	if b {
		return []uint8{1}
	}
	return []uint8{0}
}

func getStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func bitToBool(b []uint8) bool {
	return len(b) > 0 && b[0] == 1
}
func toBit(b bool) []uint8 {
	if b {
		return []uint8{1}
	}
	return []uint8{0}
}

// convertToMenuListVO 转换为菜单列表VO
func (s *MenuService) convertToMenuListVO(menu *model.SystemMenu) *dto.MenuListVO {
	return &dto.MenuListVO{
		ID:            menu.ID,
		Name:          menu.Name,
		Permission:    menu.Permission,
		Type:          menu.Type,
		Sort:          menu.Sort,
		ParentID:      menu.ParentID,
		Path:          menu.Path,
		Icon:          menu.Icon,
		Component:     menu.Component,
		ComponentName: menu.ComponentName,
		Status:        menu.Status,
		Visible:       bitToBool(menu.Visible),
		KeepAlive:     bitToBool(menu.KeepAlive),
		AlwaysShow:    bitToBool(menu.AlwaysShow),
		Creator:       menu.Creator,
		CreateAt:      menu.CreateAt,
		Updater:       menu.Updater,
		UpdateAt:      menu.UpdateAt,
	}
}

// buildMenuTree 构建菜单树
func (s *MenuService) buildMenuTree(menus []*model.SystemMenu, parentID int64) []*dto.MenuTreeVO {
	tree := make([]*dto.MenuTreeVO, 0)
	for _, menu := range menus {
		if menu.ParentID == parentID {
			// 获取菜单路径，如果为空则使用默认路径
			path := ""
			if menu.Path != nil && *menu.Path != "" {
				path = *menu.Path
			} else {
				path = fmt.Sprintf("/menu-%d", menu.ID)
			}

			node := &dto.MenuTreeVO{
				ID:            menu.ID,
				Name:          menu.Name,
				Type:          menu.Type,
				Sort:          menu.Sort,
				Permission:    menu.Permission,
				Path:          menu.Path,
				Icon:          menu.Icon,
				Component:     menu.Component,
				ComponentName: menu.ComponentName,
				Visible:       menu.Visible,
				KeepAlive:     menu.KeepAlive,
				AlwaysShow:    menu.AlwaysShow,
				Status:        menu.Status,
				// 复制name到label字段，供前端使用
				Label:    menu.Name,
				ParentID: menu.ParentID,
			}

			// 添加路径字段，如果存在则使用，否则生成默认路径
			if path != "" {
				node.Path = &path
			}

			// 添加图标字段
			if menu.Icon != nil && *menu.Icon != "" {
				node.Icon = menu.Icon
			}
			children := s.buildMenuTree(menus, menu.ID)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// checkCircularDependency 检查循环依赖
func (s *MenuService) checkCircularDependency(menuID, newParentID int64) error {
	if newParentID == 0 {
		return nil // 根节点允许
	}

	if menuID == newParentID {
		return fmt.Errorf("不能设置自己为父菜单")
	}

	// 递归检查
	currentParentID := newParentID
	for currentParentID != 0 {
		if currentParentID == menuID {
			return fmt.Errorf("会形成循环依赖")
		}

		// 获取父菜单信息
		menu, err := s.m.GetByID(currentParentID)
		if err != nil {
			return nil // 如果父菜单不存在，跳过检查
		}

		currentParentID = menu.ParentID
	}

	return nil
}
