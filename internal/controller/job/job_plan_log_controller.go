package job

import (
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service/job"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type JobPlanLogController struct {
	jobPlanLogService *job.JobPlanLogService
}

func NewJobPlanLogController(jobPlanLogService *job.JobPlanLogService) *JobPlanLogController {
	return &JobPlanLogController{
		jobPlanLogService: jobPlanLogService,
	}
}

// @Tags 作业日志管理
// @Summary 获取日志分页
// @Param request body dto.JobPlanLogPageRequest true "请求参数"
// @Router /jobs/log/page [post]
func (ctrl *JobPlanLogController) GetJobPlanLogPage(c *gin.Context) {
	var req dto.JobPlanLogPageRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobPlanLogService.GetJobPlanLogPage(c, req)
}
