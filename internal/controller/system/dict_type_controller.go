package system

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type DictTypeController struct{ s *service.DictTypeService }

func NewDictTypeController(s *service.DictTypeService) *DictTypeController {
	return &DictTypeController{s: s}
}

// Page 获取字典类型分页列表
// @Summary 获取字典类型分页列表
// @Description 根据条件获取字典类型的分页数据
// @Tags 字典类型管理
// @Accept json
// @Produce json
// @Param dictTypePageRequest body dto.DictTypePageRequest true "字典类型分页请求参数"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dict/type/page [post]
func (ctl *DictTypeController) Page(c *gin.Context) {
	var req dto.DictTypePageRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctl.s.Page(c, req)
}

// Detail 获取字典类型详情
// @Summary 获取字典类型详情
// @Description 根据ID获取字典类型的详细信息
// @Tags 字典类型管理
// @Accept json
// @Produce json
// @Param id path int true "字典类型ID"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dict/type/{id} [get]
func (ctl *DictTypeController) Detail(c *gin.Context) {
	var id int64
	util.GetParam(c, "id", &id, func(param interface{}) {
		if id <= 0 {
			common.Fail(c, common.BadRequest)
			c.Abort()
			return
		}
	})
	if c.IsAborted() {
		return
	}
	ctl.s.Detail(c, id)
}

// Create 创建字典类型
// @Summary 创建字典类型
// @Description 创建新的字典类型
// @Tags 字典类型管理
// @Accept json
// @Produce json
// @Param dictTypeSaveRequest body dto.DictTypeSaveRequest true "字典类型保存请求参数"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dict/type [post]
func (ctl *DictTypeController) Create(c *gin.Context) {
	var req dto.DictTypeSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		common.Fail(c, common.BadRequest)
		return
	}
	ctl.s.Create(c, req)
}

// Update 更新字典类型
// @Summary 更新字典类型
// @Description 根据ID更新字典类型信息
// @Tags 字典类型管理
// @Accept json
// @Produce json
// @Param dictTypeSaveRequest body dto.DictTypeSaveRequest true "字典类型保存请求参数"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dict/type [put]
func (ctl *DictTypeController) Update(c *gin.Context) {
	var req dto.DictTypeSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		common.Fail(c, common.BadRequest)
		return
	}
	ctl.s.Update(c, req)
}

// Remove 删除字典类型
// @Summary 删除字典类型
// @Description 根据ID删除字典类型
// @Tags 字典类型管理
// @Accept json
// @Produce json
// @Param id path int true "字典类型ID"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dict/type/{id} [delete]
func (ctl *DictTypeController) Remove(c *gin.Context) {
	var id int64
	util.GetParam(c, "id", &id, func(param interface{}) {
		if id <= 0 {
			common.Fail(c, common.BadRequest)
			c.Abort()
			return
		}
	})
	if c.IsAborted() {
		return
	}
	ctl.s.Remove(c, id)
}

// UpdateStatus 更新字典类型状态
// @Summary 更新字典类型状态
// @Description 根据ID更新字典类型的状态
// @Tags 字典类型管理
// @Accept json
// @Produce json
// @Param id path int true "字典类型ID"
// @Param status path int true "状态值"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dict/type/{id}/{status} [put]
func (ctl *DictTypeController) UpdateStatus(c *gin.Context) {
	var id int64
	util.GetParam(c, "id", &id, func(param interface{}) {
		if id <= 0 {
			common.Fail(c, common.BadRequest)
			c.Abort()
			return
		}
	})
	var status int32
	util.GetParam(c, "status", &status, func(param interface{}) {})
	if c.IsAborted() {
		return
	}
	ctl.s.UpdateStatus(c, id, status)
}
