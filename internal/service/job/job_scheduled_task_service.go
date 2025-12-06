package job

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper/job"

	"github.com/gin-gonic/gin"
)

type JobScheduledTaskService struct {
	jobScheduledTaskMapper *job.JobScheduledTaskMapper
}

func NewJobScheduledTaskService(jobScheduledTaskMapper *job.JobScheduledTaskMapper) *JobScheduledTaskService {
	return &JobScheduledTaskService{
		jobScheduledTaskMapper: jobScheduledTaskMapper,
	}
}

func (s *JobScheduledTaskService) CreateJobScheduledTask(c *gin.Context, req dto.JobScheduledTaskSaveRequest) {
	task := &model.JobScheduledTask{
		Name:     req.Name,
		PlanID:   req.PlanID,
		Strategy: req.Strategy,
		Status:   uint32(req.Status),
	}
	err := s.jobScheduledTaskMapper.InsertJobScheduledTask(task)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScheduledTaskService) UpdateJobScheduledTask(c *gin.Context, req dto.JobScheduledTaskSaveRequest) {
	task, err := s.jobScheduledTaskMapper.GetJobScheduledTaskById(req.ID)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}

	task.Name = req.Name
	task.PlanID = req.PlanID
	task.Strategy = req.Strategy
	task.Status = uint32(req.Status)

	err = s.jobScheduledTaskMapper.UpdateJobScheduledTask(task)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScheduledTaskService) DeleteJobScheduledTask(c *gin.Context, id int64) {
	err := s.jobScheduledTaskMapper.DeleteJobScheduledTask(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobScheduledTaskService) GetJobScheduledTaskPage(c *gin.Context, req dto.JobScheduledTaskPageRequest) {
	pageResult, err := s.jobScheduledTaskMapper.GetJobScheduledTaskPage(req)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, pageResult)
}

func (s *JobScheduledTaskService) GetJobScheduledTaskById(id int64) (*model.JobScheduledTask, error) {
	return s.jobScheduledTaskMapper.GetJobScheduledTaskById(id)
}

func (s *JobScheduledTaskService) UpdateJobScheduledTaskStatus(jobScheduledTaskRequest *dto.JobScheduledTaskStatusRequest) (bool, error) {
	task, err := s.jobScheduledTaskMapper.GetJobScheduledTaskById(jobScheduledTaskRequest.Id)
	if err != nil || task == nil {
		return false, common.TaskNotExist
	}
	return s.jobScheduledTaskMapper.UpdateJobScheduledTaskStatus(jobScheduledTaskRequest.Id, jobScheduledTaskRequest.Status)
}
