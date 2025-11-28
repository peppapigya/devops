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

func (ctl *DictTypeController) Page(c *gin.Context) {
	var req dto.DictTypePageRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctl.s.Page(c, req)
}
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
func (ctl *DictTypeController) Create(c *gin.Context) {
	var req dto.DictTypeSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		common.Fail(c, common.BadRequest)
		return
	}
	ctl.s.Create(c, req)
}
func (ctl *DictTypeController) Update(c *gin.Context) {
	var req dto.DictTypeSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		common.Fail(c, common.BadRequest)
		return
	}
	ctl.s.Update(c, req)
}
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
