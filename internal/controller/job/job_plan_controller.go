package job

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type JobPlanController struct {
	jobPlanService *service.JobPlanService
}

func NewJobPlanController(jobPlanService *service.JobPlanService) *JobPlanController {
	return &JobPlanController{
		jobPlanService: jobPlanService,
	}
}

// @Tags 作业计划管理
// @Summary 创建计划
// @Param request body dto.JobPlanSaveRequest true "请求参数"
// @Router /jobs/plan/create [post]
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
// @Router /jobs/plan/update [post]
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
// @Router /jobs/plan/delete [delete]
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
	ctrl.jobPlanService.GetJobPlanPage(c, req)
}

// @Tags 作业计划管理
// @Summary 获取计划详情
// @Param id path int true "计划ID"
// @Router /jobs/plan/detail [get]
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
