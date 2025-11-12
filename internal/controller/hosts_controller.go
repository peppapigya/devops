// internal/controller/host_controller.go

package controller

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service"
	"k8s-platform-go/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

type HostController struct {
	hostService *service.HostService
}

func NewHostController(hostService *service.HostService) *HostController {
	return &HostController{
		hostService: hostService,
	}
}

// @Tags 主机管理
// @Summary 获取主机列表分页
// @Param hostPageRequest body dto.HostPageRequest true "主机列表分页请求参数"
// @Router /hosts/page [post]
func (hostController *HostController) GetHostPage(c *gin.Context) {
	var hostPageRequest dto.HostPageRequest
	if ok := util.BindAndValidate(c, &hostPageRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	hostController.hostService.GetHostPage(hostPageRequest, c)
}

// @Tags 主机管理
// @Summary 创建主机
// @Param createHostRequest body dto.CreateHostDTO true "创建主机请求参数"
// @Router /hosts [post]
func (hostController *HostController) CreateHost(c *gin.Context) {
	var createHostRequest dto.CreateHostDTO
	if ok := util.BindAndValidate(c, &createHostRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	hostController.hostService.CreateHost(createHostRequest, c)
}

// @Tags 主机管理
// @Summary 更新主机信息
// @Param id path int true "主机ID"
// @Param updateHostRequest body dto.UpdateHostDTO true "更新主机请求参数"
// @Router /hosts/{id} [put]
func (hostController *HostController) UpdateHost(c *gin.Context) {
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

	var updateHostRequest dto.UpdateHostDTO
	if ok := util.BindAndValidate(c, &updateHostRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	hostController.hostService.UpdateHost(id, updateHostRequest, c)
}

// @Tags 主机管理
// @Summary 删除主机
// @Param id path int true "主机ID"
// @Router /hosts/{id} [delete]
func (hostController *HostController) DeleteHost(c *gin.Context) {
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

	hostController.hostService.DeleteHost(id, c)
}

// @Tags 主机管理
// @Summary 测试主机连接
// @Param id path int true "主机ID"
// @Router /hosts/{id}/test [post]
func (hostController *HostController) TestConnection(c *gin.Context) {
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

	result, err := hostController.hostService.TestConnection(id)
	if err != nil {
		common.Fail(c, err)
		return
	}
	c.String(200, result) // 返回纯文本结果
}

// @Tags 主机管理
// @Summary 巡检主机状态
// @Param id path int true "主机ID"
// @Router /hosts/{id}/inspect [post]
func (hostController *HostController) InspectHost(c *gin.Context) {
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

	result, err := hostController.hostService.InspectHost(id)
	if err != nil {
		common.Fail(c, err)
		return
	}
	common.Success(c, result) // 返回JSON格式的结果
}
