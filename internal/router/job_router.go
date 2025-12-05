package router

import (
	"k8s-platform-go/internal/util/wireinfo"

	"github.com/gin-gonic/gin"
)

func InitJobRouter(rg *gin.RouterGroup) {
	jobScriptController := wireinfo.InitializeJobScriptController()
	jobPlanController := wireinfo.InitializeJobPlanController()
	jobScheduledTaskController := wireinfo.InitializeJobScheduledTaskController()
	jobPlanLogController := wireinfo.InitializeJobPlanLogController()

	jobGroup := rg.Group("/jobs")
	{
		// Script routes
		scriptGroup := jobGroup.Group("/script")
		{
			scriptGroup.POST("/", jobScriptController.CreateJobScript)
			scriptGroup.PUT("/", jobScriptController.UpdateJobScript)
			scriptGroup.DELETE("/:id", jobScriptController.DeleteJobScript)
			scriptGroup.POST("/page", jobScriptController.GetJobScriptPage)
			scriptGroup.GET("/:id", jobScriptController.GetJobScriptById)
			scriptGroup.GET("/select", jobScriptController.GetJobScriptSelect)
			scriptGroup.POST("/execute", jobScriptController.ExecuteJobScript)
		}

		// Plan routes
		planGroup := jobGroup.Group("/plan")
		{
			planGroup.POST("/", jobPlanController.CreateJobPlan)
			planGroup.PUT("/", jobPlanController.UpdateJobPlan)
			planGroup.DELETE("/:id", jobPlanController.DeleteJobPlan)
			planGroup.POST("/page", jobPlanController.GetJobPlanPage)
			planGroup.GET("/:id", jobPlanController.GetJobPlanById)
			planGroup.GET("/list", jobPlanController.GetJobPlanSelectList)
		}

		// Schedule routes
		scheduleGroup := jobGroup.Group("/schedule")
		{
			scheduleGroup.POST("/", jobScheduledTaskController.CreateJobScheduledTask)
			scheduleGroup.PUT("/", jobScheduledTaskController.UpdateJobScheduledTask)
			scheduleGroup.GET("/:id", jobScheduledTaskController.GetJobScheduledTaskById)
			scheduleGroup.DELETE("/:id", jobScheduledTaskController.DeleteJobScheduledTask)
			scheduleGroup.POST("/page", jobScheduledTaskController.GetJobScheduledTaskPage)
			scheduleGroup.PUT("/status", jobScheduledTaskController.UpdateJobScheduledTaskStatus)
		}

		// Log routes
		logGroup := jobGroup.Group("/log")
		{
			logGroup.POST("/page", jobPlanLogController.GetJobPlanLogPage)
		}
	}
}
