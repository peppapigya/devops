package service

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/mapper"

	"github.com/gin-gonic/gin"
)

type JobPlanLogService struct {
	jobPlanLogMapper *mapper.JobPlanLogMapper
}

func NewJobPlanLogService(jobPlanLogMapper *mapper.JobPlanLogMapper) *JobPlanLogService {
	return &JobPlanLogService{
		jobPlanLogMapper: jobPlanLogMapper,
	}
}

func (s *JobPlanLogService) GetJobPlanLogPage(c *gin.Context, req dto.JobPlanLogPageRequest) {
	pageResult, err := s.jobPlanLogMapper.GetJobPlanLogPage(req)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, pageResult)
}
