package job

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/mapper/job"

	"github.com/gin-gonic/gin"
)

type JobPlanLogService struct {
	jobPlanLogMapper *job.JobPlanLogMapper
}

func NewJobPlanLogService(jobPlanLogMapper *job.JobPlanLogMapper) *JobPlanLogService {
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
