package job

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service/job"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type JobPlanController struct {
	jobPlanService *job.JobPlanService
}

func NewJobPlanController(jobPlanService *job.JobPlanService) *JobPlanController {
	return &JobPlanController{
		jobPlanService: jobPlanService,
	}
}

// @Tags 作业计划管理
// @Summary 创建计划
// @Param request body dto.JobPlanSaveRequest true "请求参数"
// @Router /jobs/plan/ [post]
func (ctrl *JobPlanController) CreateJobPlan(c *gin.Context) {
	var req dto.JobPlanSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobPlanService.CreateJobPlan(c, req)
}

// @Tags 作业计划管理
// @Summary 更新计划
// @Param request body dto.JobPlanSaveRequest true "请求参数"
// @Router /jobs/plan/ [post]
func (ctrl *JobPlanController) UpdateJobPlan(c *gin.Context) {
	var req dto.JobPlanSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobPlanService.UpdateJobPlan(c, req)
}

// @Tags 作业计划管理
// @Summary 删除计划
// @Param id path int true "计划ID"
// @Router /jobs/plan/ [delete]
func (ctrl *JobPlanController) DeleteJobPlan(c *gin.Context) {
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
	ctrl.jobPlanService.DeleteJobPlan(c, id)
}

// @Tags 作业计划管理
// @Summary 获取计划分页
// @Param request body dto.JobPlanPageRequest true "请求参数"
// @Router /jobs/plan/page [post]
func (ctrl *JobPlanController) GetJobPlanPage(c *gin.Context) {
	var req dto.JobPlanPageRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	res, err := ctrl.jobPlanService.GetJobPlanPage(req)
	if err != nil {
		common.BusinessFail(c, err.Error())
		c.Abort()
		return
	}
	common.Success(c, res)
}

// @Tags 作业计划管理
// @Summary 获取计划详情
// @Param id path int true "计划ID"
// @Router /jobs/plan/ [get]
func (ctrl *JobPlanController) GetJobPlanById(c *gin.Context) {
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
	ctrl.jobPlanService.GetJobPlanById(c, id)
}

// @Tags 作业计划管理
// @Summary 获取计划下拉列表
// @Router /jobs/plan/list [get]
func (ctrl *JobPlanController) GetJobPlanSelectList(c *gin.Context) {
	list, err := ctrl.jobPlanService.GetJobPlanSelectList()
	if err != nil {
		common.BusinessFail(c, err.Error())
		c.Abort()
		return
	}
	common.Success(c, list)
}
