package db

import (
	"gorm.io/gorm"
)

type BaseMapper[T any] struct {
	DB    *gorm.DB
	query *gorm.DB
}

// GetDB 获取DB，支持原始的条件查询
func (base *BaseMapper[T]) GetDB() *gorm.DB {
	return base.DB
}
func (base *BaseMapper[T]) GetQuery() *gorm.DB {
	return base.query
}
func NewBaseMapper[T any](DB *gorm.DB) *BaseMapper[T] {
	return &BaseMapper[T]{
		DB:    DB,
		query: db.Model(new(T)),
	}
}

func (base *BaseMapper[T]) Where(query string, args ...interface{}) *BaseMapper[T] {
	base.query = base.query.Where(query, args...)
	return base
}

// Select 查询
func (base *BaseMapper[T]) Select(query interface{}, args ...interface{}) *BaseMapper[T] {
	base.query = base.query.Select(query, args)
	return base
}

// WhereIf 条件查询
func (base *BaseMapper[T]) WhereIf(condition bool, query string, args ...interface{}) *BaseMapper[T] {
	if condition {
		base.query = base.query.Where(query, args...)
	}
	return base
}

// SelectPage 分页查询
func (base *BaseMapper[T]) SelectPage(pageNum int, pageSize int) *BaseMapper[T] {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageNum - 1) * pageSize
	base.query = base.query.Offset(offset).Limit(pageSize)
	return base
}

// Count 查询总数
func (base *BaseMapper[T]) Count(total *int64) error {
	return base.query.Count(total).Error
}

// SelectByIds 根据id批量查询
func (base *BaseMapper[T]) SelectByIds(result *[]T, ids []int64) error {
	return base.query.Where("id in ?", ids).Find(result).Error
}

// SelectById 根据id查询
func (base *BaseMapper[T]) SelectById(result *T, id int64) {
	base.query.Where("id = ?", id).Find(result)
}

func (base *BaseMapper[T]) SelectOne(result *T) error {
	return base.query.Take(result).Error
}

func (base *BaseMapper[T]) Find(dest interface{}, conds ...interface{}) error {
	return base.query.Find(&dest).Error
}
