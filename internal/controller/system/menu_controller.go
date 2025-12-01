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
// @Summary 获取菜单路由
// @Description 获取当前用户的菜单路由信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /menu/routes [get]
func (ctl *MenuController) Routes(c *gin.Context) {
	ctl.s.Routes(c)
}

// Create 创建菜单
// @Summary 创建菜单
// @Description 创建新的菜单项
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu body dto.MenuCreateDTO true "菜单创建信息"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数验证失败"
// @Failure 500 {object} common.Response "业务处理失败"
// @Router /menu [post]
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
// @Summary 更新菜单
// @Description 根据ID更新菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu body dto.MenuUpdateDTO true "菜单更新信息"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数验证失败"
// @Failure 500 {object} common.Response "业务处理失败"
// @Router /menu [put]
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
// @Summary 删除菜单
// @Description 根据ID删除指定菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path string true "菜单ID"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数验证失败"
// @Failure 500 {object} common.Response "业务处理失败"
// @Router /menu/{id} [delete]
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
// @Summary 获取菜单列表
// @Description 根据条件查询菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param name query string false "菜单名称"
// @Param status query string false "菜单状态"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数验证失败"
// @Failure 500 {object} common.Response "业务处理失败"
// @Router /menu/list [get]
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
// @Summary 获取菜单树
// @Description 获取完整的菜单树形结构
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response "业务处理失败"
// @Router /menu/tree [get]
func (ctl *MenuController) Tree(c *gin.Context) {
	if err := ctl.s.GetMenuTree(c); err != nil {
		common.BusinessFail(c, err.Error())
	}
}

// Options 获取菜单选项
// @Summary 获取菜单选项
// @Description 获取菜单下拉选项数据
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response "业务处理失败"
// @Router /menu/options [get]
func (ctl *MenuController) Options(c *gin.Context) {
	if err := ctl.s.GetMenuOptions(c); err != nil {
		common.BusinessFail(c, err.Error())
	}
}

// Seed 初始化默认菜单
// @Summary 初始化默认菜单
// @Description 初始化系统默认菜单数据
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /menu/seed [post]
func (ctl *MenuController) Seed() {
	ctl.s.SeedDefault()
}

// GetById 根据ID获取菜单详情
// @Summary 根据ID获取菜单详情
// @Description 根据菜单ID获取菜单详细信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id path string true "菜单ID"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数验证失败"
// @Failure 500 {object} common.Response "业务处理失败"
// @Router /menu/{id} [get]
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
