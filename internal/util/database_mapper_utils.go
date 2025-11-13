package util

import (
	"log"

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
	Data interface{} `json:"data"`
}

// WhereIf 条件查询
func WhereIf(condition bool, cond gen.Condition) gen.Condition {
	if condition {
		return cond
	}
	return nil
}

// FindPageResult 查询分页结果
func FindPageResult[T any](dao gen.DO, pageNum, pageSize int, isAnd bool, conditions ...gen.Condition) (PageInfoResponse[T], error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	countDAO := dao.UnderlyingDB().Session(&gorm.Session{NewDB: true}).Table(dao.TableName())
	queryDao := dao.Session(&gorm.Session{NewDB: true})
	// 去除条件为空的元素
	conditions = RemoveEmptyConditions(conditions...)
	// 拼接条件
	for _, cond := range conditions {
		if isAnd {
			countDAO = countDAO.Where(cond)
			queryDao = queryDao.Where(cond)
		} else {
			countDAO = countDAO.Or(cond)
			queryDao = queryDao.Or(cond)
		}
	}

	var result PageInfoResponse[T]
	// 查询总数
	var total int64
	if err := countDAO.Count(&total).Error; err != nil {
		log.Printf("查询总数失败: %v", err)
		return result, err
	}
	offset := (pageNum - 1) * pageSize
	list, err := queryDao.Offset(offset).Limit(pageSize).Find()
	if err != nil {
		return result, err
	}

	// 拼接结果
	result = PageInfoResponse[T]{
		PageNum:  pageNum,
		PageSize: pageSize,
		Total:    total,
		Data:     list,
	}

	return result, nil

}

func RemoveEmptyConditions(conditions ...gen.Condition) (result []gen.Condition) {
	for _, condition := range conditions {
		if condition != nil {
			result = append(result, condition)
		}
	}
	return result
}
