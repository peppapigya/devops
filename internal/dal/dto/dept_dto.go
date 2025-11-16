package dto

type DeptPageRequest struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

type DeptSaveRequest struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	ParentID     *int64  `json:"parentId"`
	Sort         int32   `json:"sort"`
	LeaderUserID *int64  `json:"leaderUserId"`
	Phone        *string `json:"phone"`
	Email        *string `json:"email"`
	Status       int32   `json:"status"`
}
