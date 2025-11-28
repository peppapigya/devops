package system

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service"

	"github.com/gin-gonic/gin"
)

type MenuController struct{ s *service.MenuService }

func NewMenuController(s *service.MenuService) *MenuController {
	return &MenuController{s: s}
}

// Routes 获取菜单路由
func (ctl *MenuController) Routes(c *gin.Context) {
	ctl.s.Routes(c)
}

// Create 创建菜单
func (ctl *MenuController) Create(c *gin.Context) {
	var createDTO dto.MenuCreateDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		common.ValidateFail(c, err.Error())
		return
	}

	if err := ctl.s.CreateMenu(c, &createDTO); err != nil {
		common.BusinessFail(c, err.Error())
	}
}

// Update 更新菜单
func (ctl *MenuController) Update(c *gin.Context) {
	var updateDTO dto.MenuUpdateDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		common.ValidateFail(c, err.Error())
		return
	}

	if err := ctl.s.UpdateMenu(c, &updateDTO); err != nil {
		common.BusinessFail(c, err.Error())
	}
}

// Delete 删除菜单
func (ctl *MenuController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ValidateFail(c, "菜单ID不能为空")
		return
	}

	menuID := common.StrToInt64(id)
	if menuID <= 0 {
		common.ValidateFail(c, "菜单ID格式不正确")
		return
	}

	if err := ctl.s.DeleteMenu(c, menuID); err != nil {
		common.BusinessFail(c, err.Error())
	}
}

// List 获取菜单列表
func (ctl *MenuController) List(c *gin.Context) {
	var queryDTO dto.MenuQueryDTO
	if err := c.ShouldBindQuery(&queryDTO); err != nil {
		common.ValidateFail(c, err.Error())
		return
	}

	if err := ctl.s.GetMenuList(c, &queryDTO); err != nil {
		common.BusinessFail(c, err.Error())
	}
}

// Tree 获取菜单树
func (ctl *MenuController) Tree(c *gin.Context) {
	if err := ctl.s.GetMenuTree(c); err != nil {
		common.BusinessFail(c, err.Error())
	}
}

// Options 获取菜单选项
func (ctl *MenuController) Options(c *gin.Context) {
	if err := ctl.s.GetMenuOptions(c); err != nil {
		common.BusinessFail(c, err.Error())
	}
}

// Seed 初始化默认菜单
func (ctl *MenuController) Seed() {
	ctl.s.SeedDefault()
}

// GetById 根据ID获取菜单详情
func (ctl *MenuController) GetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		common.ValidateFail(c, "菜单ID不能为空")
		return
	}

	menuID := common.StrToInt64(id)
	if menuID <= 0 {
		common.ValidateFail(c, "菜单ID格式不正确")
		return
	}

	if err := ctl.s.GetMenuById(c, menuID); err != nil {
		common.BusinessFail(c, err.Error())
	}
}
