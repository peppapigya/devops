package job

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper/job"

	"github.com/gin-gonic/gin"
)

type JobPlanService struct {
	jobPlanMapper       *job.JobPlanMapper
	jobPlanScriptMapper *job.JobPlanScriptMapper
}

func NewJobPlanService(jobPlanMapper *job.JobPlanMapper, jobPlanScriptMapper *job.JobPlanScriptMapper) *JobPlanService {
	return &JobPlanService{
		jobPlanMapper:       jobPlanMapper,
		jobPlanScriptMapper: jobPlanScriptMapper,
	}
}

func (s *JobPlanService) CreateJobPlan(c *gin.Context, req dto.JobPlanSaveRequest) {
	plan := &model.JobPlan{
		Name:        req.Name,
		GlobalVars:  req.GlobalVars,
		HostIds:     req.HostIds,
		HostGroupID: &req.HostGroupId,
		Remark:      &req.Remark,
	}

	err := s.jobPlanMapper.InsertJobPlan(plan)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}

	// Insert relations
	if len(req.ScriptIds) > 0 {
		var planScripts []model.JobPlanScript
		for i, scriptId := range req.ScriptIds {
			planScripts = append(planScripts, model.JobPlanScript{
				PlanID:   plan.ID,
				ScriptID: int32(scriptId),
				Sort:     uint32(i),
			})
		}
		s.jobPlanScriptMapper.InsertJobPlanScripts(planScripts)
	}

	common.Success(c, true)
}

func (s *JobPlanService) UpdateJobPlan(c *gin.Context, req dto.JobPlanSaveRequest) {
	plan, err := s.jobPlanMapper.GetJobPlanById(req.ID)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}

	plan.Name = req.Name
	plan.GlobalVars = req.GlobalVars
	plan.HostIds = req.HostIds
	plan.HostGroupID = &req.HostGroupId
	plan.Remark = &req.Remark

	err = s.jobPlanMapper.UpdateJobPlan(plan)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}

	// Update relations: Delete all and re-insert
	s.jobPlanScriptMapper.DeleteJobPlanScripts(int64(plan.ID))
	if len(req.ScriptIds) > 0 {
		var planScripts []model.JobPlanScript
		for i, scriptId := range req.ScriptIds {
			planScripts = append(planScripts, model.JobPlanScript{
				PlanID:   plan.ID,
				ScriptID: int32(scriptId),
				Sort:     uint32(i),
			})
		}
		s.jobPlanScriptMapper.InsertJobPlanScripts(planScripts)
	}

	common.Success(c, true)
}

func (s *JobPlanService) DeleteJobPlan(c *gin.Context, id int64) {
	err := s.jobPlanMapper.DeleteJobPlan(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (s *JobPlanService) GetJobPlanPage(c *gin.Context, req dto.JobPlanPageRequest) {
	pageResult, err := s.jobPlanMapper.GetJobPlanPage(req)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, pageResult)
}

func (s *JobPlanService) GetJobPlanById(c *gin.Context, id int64) {
	plan, err := s.jobPlanMapper.GetJobPlanById(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, plan)
}
