package controller

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type DeptController struct{ svc *service.DeptService }

func NewDeptController(s *service.DeptService) *DeptController { return &DeptController{svc: s} }

func (ctl *DeptController) Page(c *gin.Context) {
	var req dto.DeptPageRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctl.svc.Page(c, req)
}

func (ctl *DeptController) Tree(c *gin.Context) { ctl.svc.Tree(c) }

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

func (ctl *DeptController) Create(c *gin.Context) {
	var req dto.DeptSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		common.Fail(c, common.BadRequest)
		return
	}
	ctl.svc.Create(c, req)
}

func (ctl *DeptController) Update(c *gin.Context) {
	var req dto.DeptSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		common.Fail(c, common.BadRequest)
		return
	}
	ctl.svc.Update(c, req)
}

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
