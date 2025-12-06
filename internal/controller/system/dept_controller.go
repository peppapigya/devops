package system

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service/system"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type DeptController struct{ svc *system.DeptService }

func NewDeptController(s *system.DeptService) *DeptController { return &DeptController{svc: s} }

// Page 获取部门分页列表
// @Summary 获取部门分页列表
// @Description 根据条件获取部门的分页数据
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param deptPageRequest body dto.DeptPageRequest true "部门分页请求参数"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dept/page [post]
func (ctl *DeptController) Page(c *gin.Context) {
	var req dto.DeptPageRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctl.svc.Page(c, req)
}

// Tree 获取部门树形结构
// @Summary 获取部门树形结构
// @Description 获取所有部门的树形结构数据
// @Tags 部门管理
// @Accept json
// @Produce json
// @Success 200 {object} common.Response
// @Router /dept/tree [get]
func (ctl *DeptController) Tree(c *gin.Context) { ctl.svc.Tree(c) }

// Detail 获取部门详情
// @Summary 获取部门详情
// @Description 根据ID获取部门的详细信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dept/{id} [get]
func (ctl *DeptController) Detail(c *gin.Context) {
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
	ctl.svc.Detail(c, id)
}

// Create 创建部门
// @Summary 创建部门
// @Description 创建新的部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param deptSaveRequest body dto.DeptSaveRequest true "部门保存请求参数"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dept [post]
func (ctl *DeptController) Create(c *gin.Context) {
	var req dto.DeptSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		common.Fail(c, common.BadRequest)
		return
	}
	ctl.svc.Create(c, req)
}

// Update 更新部门
// @Summary 更新部门
// @Description 根据ID更新部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param deptSaveRequest body dto.DeptSaveRequest true "部门保存请求参数"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dept [put]
func (ctl *DeptController) Update(c *gin.Context) {
	var req dto.DeptSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		common.Fail(c, common.BadRequest)
		return
	}
	ctl.svc.Update(c, req)
}

// Remove 删除部门
// @Summary 删除部门
// @Description 根据ID删除部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dept/{id} [delete]
func (ctl *DeptController) Remove(c *gin.Context) {
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
	ctl.svc.Remove(c, id)
}

// UpdateStatus 更新部门状态
// @Summary 更新部门状态
// @Description 根据ID更新部门状态
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id path int true "部门ID"
// @Param status path int true "状态值"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.Response "参数校验失败"
// @Router /dept/{id}/{status} [put]
func (ctl *DeptController) UpdateStatus(c *gin.Context) {
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
	ctl.svc.UpdateStatus(c, id, status)
}
