package service

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper"
	"log"

	"github.com/gin-gonic/gin"
)

type DictTypeService struct{ m *mapper.DictTypeMapper }

func NewDictTypeService(m *mapper.DictTypeMapper) *DictTypeService { return &DictTypeService{m: m} }

func (s *DictTypeService) Page(c *gin.Context, req dto.DictTypePageRequest) {
	res, err := s.m.SelectPage(req)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, res)
}
func (s *DictTypeService) Detail(c *gin.Context, id int64) {
	d, err := s.m.SelectByID(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, d)
}
func (s *DictTypeService) Create(c *gin.Context, req dto.DictTypeSaveRequest) {
	v := &model.SystemDictType{ID: req.ID, Name: req.Name, Type: req.Type, Status: req.Status, Remark: &req.Remark}
	if err := s.m.Insert(v); err != nil {
		log.Printf("dictType create err: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}
func (s *DictTypeService) Update(c *gin.Context, req dto.DictTypeSaveRequest) {
	v := &model.SystemDictType{ID: req.ID, Name: req.Name, Type: req.Type, Status: req.Status, Remark: &req.Remark}
	if err := s.m.Update(v); err != nil {
		log.Printf("dictType update err: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}
func (s *DictTypeService) Remove(c *gin.Context, id int64) {
	if err := s.m.DeleteByID(id); err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}
func (s *DictTypeService) UpdateStatus(c *gin.Context, id int64, status int32) {
	if err := s.m.UpdateStatus(id, status); err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}
