package job

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service/job"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type JobScriptController struct {
	jobScriptService *job.JobScriptService
}

func NewJobScriptController(jobScriptService *job.JobScriptService) *JobScriptController {
	return &JobScriptController{
		jobScriptService: jobScriptService,
	}
}

// @Tags 作业脚本管理
// @Summary 创建脚本
// @Param request body dto.JobScriptSaveRequest true "请求参数"
// @Router /jobs/script/create [post]
func (ctrl *JobScriptController) CreateJobScript(c *gin.Context) {
	var req dto.JobScriptSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobScriptService.CreateJobScript(c, req)
}

// @Tags 作业脚本管理
// @Summary 更新脚本
// @Param request body dto.JobScriptSaveRequest true "请求参数"
// @Router /jobs/script/update [post]
func (ctrl *JobScriptController) UpdateJobScript(c *gin.Context) {
	var req dto.JobScriptSaveRequest
	if ok := util.BindAndValidate(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobScriptService.UpdateJobScript(c, req)
}

// @Tags 作业脚本管理
// @Summary 删除脚本
// @Param id path int true "脚本ID"
// @Router /jobs/script/delete [delete]
func (ctrl *JobScriptController) DeleteJobScript(c *gin.Context) {
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
	ctrl.jobScriptService.DeleteJobScript(c, id)
}

// @Tags 作业脚本管理
// @Summary 获取脚本分页
// @Param request body dto.JobScriptPageRequest true "请求参数"
// @Router /jobs/script/page [post]
func (ctrl *JobScriptController) GetJobScriptPage(c *gin.Context) {
	var req dto.JobScriptPageRequest
	if ok := util.BindQueryParam(c, &req); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	ctrl.jobScriptService.GetJobScriptPage(c, req)
}

// @Tags 作业脚本管理
// @Summary 获取脚本详情
// @Param id path int true "脚本ID"
// @Router /jobs/script/detail [get]
func (ctrl *JobScriptController) GetJobScriptById(c *gin.Context) {
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
	script, err := ctrl.jobScriptService.GetJobScriptById(id)
	if err != nil {
		common.FailWithError(c, err)
		return
	}
	common.Success(c, script)
}

// @Tags 作业脚本管理
// @Summary 获取脚本下拉框
// @Param condition query string true "条件"
// @Router /jobs/script/select [get]
func (ctrl *JobScriptController) GetJobScriptSelect(c *gin.Context) {
	var condition string
	util.GetParam(c, "condition", &condition, nil)
	scriptSelect, err := ctrl.jobScriptService.GetJobScriptSelect(condition)
	if err != nil {
		common.FailWithError(c, err)
		return
	}
	common.Success(c, scriptSelect)
}

// @Tags 作业脚本管理
// @Summary 执行脚本
// @Param request body dto.ExecutorScript true "请求参数"
func (ctrl *JobScriptController) ExecuteJobScript(c *gin.Context) {
	var script dto.ExecutorScript
	if ok := util.BindAndValidate(c, &script); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	result, err := ctrl.jobScriptService.ExecuteJobScript(c, script)
	if err != nil {
		common.FailWithError(c, err)
		return
	}
	common.Success(c, result)
}

// @Tags 作业脚本管理
// @Summary 分发脚本或脚本文件
// @Param request body dto.DistributeJobScript true "请求参数"
// @Router /jobs/script/distribute [post]
func (ctrl *JobScriptController) DistributeJobScript(c *gin.Context) {
	var distribute dto.DistributeJobScript
	if ok := util.BindAndValidate(c, &distribute); !ok {
		log.Printf("参数解析失败或验证失败\n")
		c.Abort()
		return
	}
	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err == nil {
		defer func() {
			_ = file.Close()
		}()
		distribute.File = header
	}
	result, err := ctrl.jobScriptService.DistributeJobScript(c, distribute)
	if err != nil {
		common.FailWithError(c, err)
		return
	}
	common.Success(c, result)
}
