package service

import (
	"k8s-platform-go/internal/common"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/mapper"
	"log"

	"github.com/gin-gonic/gin"
)

type HostService struct {
	hostMapper *mapper.HostMapper
	context    *gin.Context
}

func NewHostService(hostMapper *mapper.HostMapper) *HostService {
	return &HostService{
		hostMapper: hostMapper,
	}
}

func (hostService *HostService) GetHostPage(request dto.HostPageRequest, c *gin.Context) {
	pageResult, err := hostService.hostMapper.SelectPageByCondition(request)
	if err != nil {
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, pageResult)
}

func (hostService *HostService) CreateHost(request dto.CreateHostDTO, c *gin.Context) {
	host := &model.Host{
		HostName:     request.HostName,
		Address:      request.Address,
		HostPort:     request.HostPort,
		Username:     request.Username,
		HostPassword: &request.HostPassword,
		Remark:       &request.Remark,
	}

	err := hostService.hostMapper.InsertHost(host)
	if err != nil {
		log.Printf("添加主机失败: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (hostService *HostService) UpdateHost(id int64, request dto.UpdateHostDTO, c *gin.Context) {
	host := &model.Host{
		HostName: request.HostName,
		Address:  request.Address,
		HostPort: request.HostPort,
		Username: request.Username,
		Remark:   &request.Remark,
	}

	// 如果密码不为空，则更新密码
	if request.HostPassword != "" {
		host.HostPassword = &request.HostPassword
	}

	err := hostService.hostMapper.UpdateHost(id, host)
	if err != nil {
		log.Printf("更新主机失败: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (hostService *HostService) DeleteHost(id int64, c *gin.Context) {
	err := hostService.hostMapper.DeleteHostById(id)
	if err != nil {
		log.Printf("删除主机失败: %v", err)
		common.Fail(c, common.ServerError)
		return
	}
	common.Success(c, true)
}

func (hostService *HostService) TestConnection(id int64) (string, *common.ErrorCode) {
	// 这里应该实现实际的连接测试逻辑
	// 目前返回模拟结果
	host, err := hostService.hostMapper.SelectHostById(id)
	if err != nil {
		return "", nil
	}

	if host == nil {
		return "主机不存在", nil
	}

	// 模拟连接测试结果
	return "连接成功", nil
}

func (hostService *HostService) InspectHost(id int64) (interface{}, *common.ErrorCode) {
	// 这里应该实现实际的巡检逻辑
	// 目前返回模拟结果
	host, err := hostService.hostMapper.SelectHostById(id)
	if err != nil {
		return nil, common.ServerError
	}

	if host == nil {
		return nil, common.ServerError
	}

	// 模拟巡检结果
	result := map[string]interface{}{
		"id":           host.ID,
		"host_name":    host.HostName,
		"address":      host.Address,
		"status":       "online",
		"cpu_usage":    "45%",
		"memory_usage": "60%",
		"disk_usage":   "70%",
	}

	return result, nil
}
