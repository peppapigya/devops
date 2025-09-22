package util

import (
	"gorm.io/gen"
	"gorm.io/gorm"
)

// PageInfoResponse 分页信息
type PageInfoResponse[T any] struct {
	// 当前页码
	PageNum int `json:"pageNum"`
	// 页面数量
	PageSize int `json:"pageSize"`
	// 数据总数
	Total int64 `json:"total"`
	// 数据
	Data []T `json:"data"`
}

// WhereIf 条件查询
func WhereIf(condition bool, cond gen.Condition) gen.Condition {
	if condition {
		return cond
	}
	return nil
}

// FindPageResult 查询分页结果
func FindPageResult[T any](dao gen.DO, pageNum, pageSize int, conditions ...gen.Condition) (PageInfoResponse[T], error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	countDAO := dao.UnderlyingDB().Session(&gorm.Session{NewDB: true}).Model(dao.TableName)
	// 拼接条件
	for _, cond := range conditions {
		if cond != nil {
			countDAO = countDAO.Where(cond)
		}
	}

	var result PageInfoResponse[T]
	// 查询总数
	total, err := dao.Count()
	if err != nil {
		return result, err
	}
	offset := (pageNum - 1) * pageSize
	list, err := dao.Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return result, err
	}
	// 拼接结果
	result = PageInfoResponse[T]{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    total,
		Data:     list.([]T),
	}
	return result, nil

}
