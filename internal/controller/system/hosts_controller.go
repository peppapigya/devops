// internal/controller/host_controller.go

package system

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/service"
	"k8s-platform-go/internal/util"
	"log"
	"strconv"

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
	var updateHostRequest dto.UpdateHostDTO
	if ok := util.BindAndValidate(c, &updateHostRequest); !ok {
		log.Printf("参数解析失败或验证失败\n")
		return
	}
	hostController.hostService.UpdateHost(updateHostRequest, c)
}

// @Tags 主机管理
// @Summary 删除主机
// @Param id path int true "主机ID"
// @Router /hosts/{id} [delete]
func (hostController *HostController) DeleteHost(c *gin.Context) {
	var id int
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.Fail(c, common.BadRequest)
		c.Abort()
		return
	}
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
	idStr := c.Param("id")

	id, err1 := strconv.Atoi(idStr)
	if err1 != nil {
		common.Fail(c, common.BadRequest)
		c.Abort()
		return
	}
	Response, err := hostController.hostService.TestConnection(id)
	if err != nil {
		common.Fail(c, err)
		return
	}
	c.String(200, Response) // 返回纯文本结果
}

// @Tags 主机管理
// @Summary 巡检主机状态
// @Param id path int true "主机ID"
// @Router /hosts/{id}/inspect [post]
func (hostController *HostController) InspectHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err1 := strconv.Atoi(idStr)
	if err1 != nil {
		common.Fail(c, common.BadRequest)
		c.Abort()
		return
	}

	Response, err := hostController.hostService.InspectHost(id)
	if err != nil {
		common.Fail(c, err)
		return
	}
	common.Success(c, Response) // 返回JSON格式的结果
}

// GetHostSelectList 获取主机下拉列表
// @Tags 主机管理
// @Summary 获取主机下拉列表
// @Description 获取所有主机的下拉选项列表，用于前端选择框展示
// @Accept json
// @Produce json
// @Success 200 {object} common.Response "成功返回主机下拉列表数据"
// @Failure 500 {object} common.Response "业务处理失败"
// @Router /hosts/select [get]
func (hostController *HostController) GetHostSelectList(c *gin.Context) {
	list, err := hostController.hostService.GetHostSelectList()
	if err != nil {
		common.BusinessFail(c, err.Error())
		return
	}
	common.Success(c, list)
}
