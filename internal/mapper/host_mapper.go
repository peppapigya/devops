package mapper

import (
	"context"
	"k8s-platform-go/internal/dal/dto"
	"k8s-platform-go/internal/dal/model"
	"k8s-platform-go/internal/dal/query"
	"k8s-platform-go/internal/util"

	"gorm.io/gorm"
)

type HostMapper struct {
	DB    *gorm.DB
	query *query.Query
}

func NewHostMapper(DB *gorm.DB) *HostMapper {
	return &HostMapper{
		DB:    DB,
		query: query.Use(DB),
	}
}

// SelectHostById 通过主机id获取主机信息
func (hostMapper *HostMapper) SelectHostById(id int) (*model.Host, error) {
	host := hostMapper.query.Host
	return hostMapper.query.Host.WithContext(context.Background()).
		Where(host.ID.Eq(uint32(id))).First()
}

// SelectPageByCondition 分页查询主机信息
func (hostMapper *HostMapper) SelectPageByCondition(request dto.HostPageRequest) (util.PageInfoResponse[model.Host], error) {
	return util.FindPageResult[model.Host](hostMapper.query.Host.DO, request.PageNum, request.PageSize)
}

// UpdateHost 更新主机信息
func (hostMapper *HostMapper) UpdateHost(host *model.Host) error {
	_, err := hostMapper.query.Host.WithContext(context.Background()).
		Where(hostMapper.query.Host.ID.Eq(host.ID)).Updates(host)
	return err
}

// DeleteHostById 删除主机
func (hostMapper *HostMapper) DeleteHostById(id int) error {
	_, err := hostMapper.query.Host.WithContext(context.Background()).
		Where(hostMapper.query.Host.ID.Eq(uint32(id))).Delete()
	return err
}

// InsertHost 添加主机信息到数据库
func (hostMapper *HostMapper) InsertHost(host *model.Host) error {
	err := hostMapper.query.Host.WithContext(context.Background()).Create(host)
	return err
}
