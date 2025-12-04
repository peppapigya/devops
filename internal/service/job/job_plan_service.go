package job

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper/job"
	"k8s-platform-go/internal/util"
	"log"

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
		HostIds:     util.ArrayToString(req.HostIds),
		HostGroupID: &req.HostGroupId,
		Remark:      &req.Remark,
	}

	err := s.jobPlanMapper.InsertJobPlan(plan)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}

	if len(req.Scripts) > 0 {
		var planScripts []model.JobPlanScript
		for i, script := range req.Scripts {
			planScripts = append(planScripts, model.JobPlanScript{
				PlanID:   plan.ID,
				ScriptID: script.ScriptId,
				Sort:     uint32(i),
				Name:     script.Name,
			})
		}
		_ = s.jobPlanScriptMapper.InsertJobPlanScripts(planScripts)
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
	plan.HostIds = util.ArrayToString(req.HostIds)
	plan.HostGroupID = &req.HostGroupId
	plan.Remark = &req.Remark

	err = s.jobPlanMapper.UpdateJobPlan(plan)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}

	_ = s.jobPlanScriptMapper.DeleteJobPlanScripts(int64(plan.ID))
	if len(req.Scripts) > 0 {
		var planScripts []model.JobPlanScript
		for i, scripts := range req.Scripts {
			planScripts = append(planScripts, model.JobPlanScript{
				PlanID:   plan.ID,
				ScriptID: scripts.ScriptId,
				Sort:     uint32(i),
				Name:     scripts.Name,
			})
		}
		_ = s.jobPlanScriptMapper.InsertJobPlanScripts(planScripts)
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

func (s *JobPlanService) GetJobPlanPage(req dto.JobPlanPageRequest) (*util.PageInfoResponse[dto.JobPlanPageResponse], error) {
	pageResult, err := s.jobPlanMapper.GetJobPlanPage(req)
	if err != nil {
		return nil, common.ServerError
	}
	jobPlans, ok := pageResult.Data.([]*model.JobPlan)
	if !ok {
		log.Fatalln("断言失败")
		return nil, common.ServerError
	}
	// 数据转化
	list, err := common.ConvertList(jobPlans, func(jobPlan *model.JobPlan) (*dto.JobPlanPageResponse, error) {
		jobPlanPage := buildJobPlanPageResponse(jobPlan)
		planScript, err := s.jobPlanScriptMapper.GetJobPlanScriptsByPlanId(jobPlan.ID)
		if err != nil {
			return nil, err
		}
		jobPlanPage.Scripts = buildJobPlanScript(planScript)
		return jobPlanPage, nil
	})
	if err != nil {
		return nil, common.ServerError
	}

	res := &util.PageInfoResponse[dto.JobPlanPageResponse]{
		Total: pageResult.Total,
		Data:  list,
	}
	return res, nil
}

func (s *JobPlanService) GetJobPlanById(c *gin.Context, id int64) {
	plan, err := s.jobPlanMapper.GetJobPlanById(id)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	planScript, err := s.jobPlanScriptMapper.GetJobPlanScriptsByPlanId(plan.ID)
	res := buildJobPlanPageResponse(plan)
	res.Scripts = buildJobPlanScript(planScript)
	common.Success(c, res)
}

// 构建计划脚本
func buildJobPlanScript(jobPlanScripts []*model.JobPlanScript) []dto.JobPlanScriptResponse {
	res := make([]dto.JobPlanScriptResponse, 0)
	for _, jobPlanScript := range jobPlanScripts {
		res = append(res, dto.JobPlanScriptResponse{
			ID:         jobPlanScript.ID,
			PlanId:     jobPlanScript.PlanID,
			ScriptId:   jobPlanScript.ScriptID,
			Sort:       jobPlanScript.Sort,
			ScriptName: jobPlanScript.Name,
		})
	}
	return res
}

// 构建计划分页响应
func buildJobPlanPageResponse(plan *model.JobPlan) *dto.JobPlanPageResponse {
	return &dto.JobPlanPageResponse{
		ID:          plan.ID,
		Name:        plan.Name,
		GlobalVars:  plan.GlobalVars,
		HostIds:     util.SplitString(plan.HostIds, ","),
		HostGroupID: plan.HostGroupID,
		Remark:      plan.Remark,
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
		DeletedAt:   plan.DeletedAt,
	}
}
