package host

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
	host := hostMapper.query.Host
	return util.FindPageResult[model.Host](hostMapper.query.Host.DO, request.PageNum, request.PageSize, false,
		util.WhereIf(request.Keyword != "", host.HostName.Like("%"+request.Keyword+"%")),
		util.WhereIf(request.Keyword != "", host.Address.Like("%"+request.Keyword+"%")),
		util.WhereIf(request.Keyword != "", host.Username.Like("%"+request.Keyword+"%")),
	)
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

// SelectHostSelectList 获取主机下拉列表
func (hostMapper *HostMapper) SelectHostSelectList() ([]*model.Host, error) {
	host := hostMapper.query.Host
	return hostMapper.query.Host.WithContext(context.Background()).Select(host.ID, host.HostName, host.Address).Find()
}

// SelectByIds 根据id列表查询主机信息
func (hostMapper *HostMapper) SelectByIds(ids []uint32) ([]*model.Host, error) {
	host := hostMapper.query.Host
	return hostMapper.query.Host.Where(host.ID.In(ids...)).Find()
}

// SelectHostByIds 根据主机id列表查询主机信息
func (hostMapper *HostMapper) SelectHostByIds(ids []uint32) ([]*model.Host, error) {
	host := hostMapper.query.Host
	return host.Where(host.ID.In(ids...)).Find()
}
