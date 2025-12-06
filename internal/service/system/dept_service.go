package system

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper/system"
	"log"

	"github.com/gin-gonic/gin"
)

type DeptService struct {
	mapper *system.DeptMapper
}

func NewDeptService(m *system.DeptMapper) *DeptService { return &DeptService{mapper: m} }

func (s *DeptService) Page(c *gin.Context, req dto.DeptPageRequest) {
	res, err := s.mapper.SelectPageByCondition(req)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, res)
}

func (s *DeptService) Tree(c *gin.Context) {
	list, err := s.mapper.ListAllEnabled()
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	nodes := buildTree(list)
	common.Success(c, nodes)
}

func (s *DeptService) Detail(c *gin.Context, id int64) {
	d, err := s.mapper.SelectByID(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, d)
}

func (s *DeptService) Create(c *gin.Context, req dto.DeptSaveRequest) {
	d := &model.SystemDept{ID: req.ID, Name: req.Name, Sort: req.Sort, Status: req.Status}
	if req.ParentID != nil {
		d.ParentID = *req.ParentID
	} else {
		d.ParentID = 0
	}
	d.LeaderUserID = req.LeaderUserID
	d.Phone = req.Phone
	d.Email = req.Email
	if err := s.mapper.Insert(d); err != nil {
		log.Printf("dept create err: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *DeptService) Update(c *gin.Context, req dto.DeptSaveRequest) {
	d := &model.SystemDept{ID: req.ID, Name: req.Name, Sort: req.Sort, Status: req.Status}
	if req.ParentID != nil {
		d.ParentID = *req.ParentID
	} else {
		d.ParentID = 0
	}
	d.LeaderUserID = req.LeaderUserID
	d.Phone = req.Phone
	d.Email = req.Email
	if err := s.mapper.Update(d); err != nil {
		log.Printf("dept update err: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *DeptService) Remove(c *gin.Context, id int64) {
	if err := s.mapper.DeleteByID(id); err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *DeptService) UpdateStatus(c *gin.Context, id int64, status int32) {
	if err := s.mapper.UpdateStatus(id, status); err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

type TreeNode struct {
	ID       int64       `json:"id"`
	Name     string      `json:"name"`
	ParentID int64       `json:"parentId"`
	Children []*TreeNode `json:"children,omitempty"`
}

func buildTree(list []*model.SystemDept) []*TreeNode {
	mp := map[int64]*TreeNode{}
	var roots []*TreeNode
	for _, d := range list {
		tn := &TreeNode{ID: d.ID, Name: d.Name, ParentID: d.ParentID}
		mp[d.ID] = tn
	}
	for _, d := range list {
		tn := mp[d.ID]
		if d.ParentID == 0 {
			roots = append(roots, tn)
		} else if parent, ok := mp[d.ParentID]; ok {
			parent.Children = append(parent.Children, tn)
		} else {
			roots = append(roots, tn)
		}
	}
	return roots
}
