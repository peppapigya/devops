package service

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper"

	"github.com/gin-gonic/gin"
)

type JobScriptService struct {
	jobScriptMapper *mapper.JobScriptMapper
}

func NewJobScriptService(jobScriptMapper *mapper.JobScriptMapper) *JobScriptService {
	return &JobScriptService{
		jobScriptMapper: jobScriptMapper,
	}
}

func (s *JobScriptService) CreateJobScript(c *gin.Context, req dto.JobScriptSaveRequest) {
	script := &model.JobScript{
		Name:       req.Name,
		Type:       req.Type,
		Category:   req.Category,
		Content:    req.Content,
		Parameters: req.Parameters,
		Timeout:    uint32(req.Timeout),
		WorkDir:    &req.WorkDir,
		Env:        &req.Env,
	}
	err := s.jobScriptMapper.InsertJobScript(script)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScriptService) UpdateJobScript(c *gin.Context, req dto.JobScriptSaveRequest) {
	script, err := s.jobScriptMapper.GetJobScriptById(req.ID)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	if script == nil {
		common.Fail(c, common.BadRequest)
		return
	}

	script.Name = req.Name
	script.Type = req.Type
	script.Category = req.Category
	script.Content = req.Content
	script.Parameters = req.Parameters
	script.Timeout = uint32(req.Timeout)
	script.WorkDir = &req.WorkDir
	script.Env = &req.Env

	err = s.jobScriptMapper.UpdateJobScript(script)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScriptService) DeleteJobScript(c *gin.Context, id int64) {
	err := s.jobScriptMapper.DeleteJobScript(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScriptService) GetJobScriptPage(c *gin.Context, req dto.JobScriptPageRequest) {
	pageResult, err := s.jobScriptMapper.GetJobScriptPage(req)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, pageResult)
}

func (s *JobScriptService) GetJobScriptById(c *gin.Context, id int64) {
	script, err := s.jobScriptMapper.GetJobScriptById(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, script)
}
