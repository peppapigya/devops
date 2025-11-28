package dto

type DictTypePageRequest struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

type DictTypeSaveRequest struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status int32  `json:"status"`
	Remark string `json:"remark"`
}
