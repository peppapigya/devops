package job

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service/job"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type JobScheduledTaskController struct {
	jobScheduledTaskService *job.JobScheduledTaskService
}

func NewJobScheduledTaskController(jobScheduledTaskService *job.JobScheduledTaskService) *JobScheduledTaskController {
	return &JobScheduledTaskController{
		jobScheduledTaskService: jobScheduledTaskService,
	}
}

// @Tags 定时任务管理
// @Summary 创建任务
// @Param request body dto.JobScheduledTaskSaveRequest true "请求参数"
// @Router /jobs/schedule/create [post]
func (ctrl *JobScheduledTaskController) CreateJobScheduledTask(c *gin.Context) {
	var req dto.JobScheduledTaskSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobScheduledTaskService.CreateJobScheduledTask(c, req)
}

// @Tags 定时任务管理
// @Summary 更新任务
// @Param request body dto.JobScheduledTaskSaveRequest true "请求参数"
// @Router /jobs/schedule/update [post]
func (ctrl *JobScheduledTaskController) UpdateJobScheduledTask(c *gin.Context) {
	var req dto.JobScheduledTaskSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobScheduledTaskService.UpdateJobScheduledTask(c, req)
}

// @Tags 定时任务管理
// @Summary 删除任务
// @Param id path int true "任务ID"
// @Router /jobs/schedule/delete [delete]
func (ctrl *JobScheduledTaskController) DeleteJobScheduledTask(c *gin.Context) {
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
	ctrl.jobScheduledTaskService.DeleteJobScheduledTask(c, id)
}

// @Tags 定时任务管理
// @Summary 获取任务分页
// @Param request body dto.JobScheduledTaskPageRequest true "请求参数"
// @Router /jobs/schedule/page [post]
func (ctrl *JobScheduledTaskController) GetJobScheduledTaskPage(c *gin.Context) {
	var req dto.JobScheduledTaskPageRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobScheduledTaskService.GetJobScheduledTaskPage(c, req)
}

// @Tags 定时任务管理
// @Summary 获取任务详情
// @Param id path int true "任务ID"
// @Router /jobs/schedule/ [get]
func (ctrl *JobScheduledTaskController) GetJobScheduledTaskById(context *gin.Context) {
	var id int64
	util.GetParam(context, "id", &id, func(param interface{}) {
		if id <= 0 {
			common.Fail(context, common.BadRequest)
			context.Abort()
			return
		}
	})
	res, err := ctrl.jobScheduledTaskService.GetJobScheduledTaskById(id)
	if err != nil {
		common.BusinessFail(context, err.Error())
		context.Abort()
		return
	}
	common.Success(context, res)
}

// @Tags 定时任务管理
// @Summary 更新任务状态
// @Param request body dto.JobScheduledTaskStatusRequest true "请求参数"
// @Router /jobs/schedule/status [put]
func (ctrl *JobScheduledTaskController) UpdateJobScheduledTaskStatus(c *gin.Context) {
	var jobScheduledTaskStatusRequest *dto.JobScheduledTaskStatusRequest
	if ok := util.BindAndValidate(c, &jobScheduledTaskStatusRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		common.Fail(c, common.BadRequest)
		c.Abort()
		return
	}
	taskStatus, err := ctrl.jobScheduledTaskService.UpdateJobScheduledTaskStatus(jobScheduledTaskStatusRequest)
	if err != nil {
		common.BusinessFail(c, err.Error())
		c.Abort()
		return
	}
	common.Success(c, taskStatus)
}
