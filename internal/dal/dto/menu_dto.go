package dto

import "time"

// MenuCreateDTO 创建菜单 DTO
type MenuCreateDTO struct {
	Name          string  `json:"name" binding:"required,min=1,max=50" label:"菜单名称"`
	Permission    string  `json:"permission" binding:"required,min=1,max=100" label:"权限标识"`
	Type          int32   `json:"type" binding:"required,oneof=1 2 3" label:"菜单类型"` // 1:目录 2:菜单 3:按钮
	Sort          int32   `json:"sort" binding:"required,min=0" label:"显示顺序"`
	ParentID      int64   `json:"parent_id" binding:"required,min=0" label:"父菜单ID"`
	Path          *string `json:"path" binding:"max=200" label:"路由地址"`
	Icon          *string `json:"icon" binding:"max=100" label:"菜单图标"`
	Component     *string `json:"component" binding:"max=255" label:"组件路径"`
	ComponentName *string `json:"component_name" binding:"max=255" label:"组件名"`
	Visible       bool    `json:"visible" label:"是否可见"`
	KeepAlive     bool    `json:"keep_alive" label:"是否缓存"`
	AlwaysShow    bool    `json:"always_show" label:"是否总是显示"`
}

// MenuUpdateDTO 更新菜单 DTO
type MenuUpdateDTO struct {
	ID            int64   `json:"id" binding:"required,min=1" label:"菜单ID"`
	Name          string  `json:"name" binding:"required,min=1,max=50" label:"菜单名称"`
	Permission    string  `json:"permission" binding:"required,min=0,max=100" label:"权限标识"`
	Type          int32   `json:"type" binding:"required,oneof=1 2 3" label:"菜单类型"`
	Sort          int32   `json:"sort" binding:"required,min=0" label:"显示顺序"`
	ParentID      int64   `json:"parentId" label:"父菜单ID"`
	Path          *string `json:"path" binding:"max=200" label:"路由地址"`
	Icon          *string `json:"icon" binding:"max=100" label:"菜单图标"`
	Component     *string `json:"component" binding:"max=255" label:"组件路径"`
	ComponentName *string `json:"componentName" label:"组件名"`
	Visible       bool    `json:"visible" label:"是否可见"`
	KeepAlive     bool    `json:"keepAlive" label:"是否缓存"`
	AlwaysShow    bool    `json:"alwaysShow" label:"是否总是显示"`
	Status        int32   `json:"status"  label:"菜单状态"`
}

// MenuQueryDTO 菜单查询 DTO
type MenuQueryDTO struct {
	Name     string `json:"name" form:"name" label:"菜单名称"`
	Type     int32  `json:"type" form:"type" label:"菜单类型"`
	Status   int32  `json:"status" form:"status" label:"菜单状态"`
	ParentID int64  `json:"parent_id" form:"parent_id" label:"父菜单ID"`
	Page     int    `json:"page" form:"page" binding:"min=1" label:"页码"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100" label:"每页数量"`
}

// MenuListVO 菜单列表 VO
type MenuListVO struct {
	ID            int64         `json:"id"`
	Name          string        `json:"name"`
	Permission    string        `json:"permission"`
	Type          int32         `json:"type"`
	Sort          int32         `json:"sort"`
	ParentID      int64         `json:"parent_id"`
	Path          *string       `json:"path"`
	Icon          *string       `json:"icon"`
	Component     *string       `json:"component"`
	ComponentName *string       `json:"component_name"`
	Status        int32         `json:"status"`
	Visible       bool          `json:"visible"`
	KeepAlive     bool          `json:"keep_alive"`
	AlwaysShow    bool          `json:"always_show"`
	Creator       *string       `json:"creator"`
	CreateAt      time.Time     `json:"create_at"`
	Updater       *string       `json:"updater"`
	UpdateAt      time.Time     `json:"update_at"`
	Children      []*MenuListVO `json:"children,omitempty"`
}

// MenuTreeVO 菜单树形 VO
type MenuTreeVO struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	Label    string        `json:"label"`          // label字段，供前端使用
	Path     *string       `json:"path,omitempty"` // 路由路径
	Icon     *string       `json:"icon,omitempty"` // 菜单图标
	ParentID int64         `json:"parent_id"`
	Children []*MenuTreeVO `json:"children,omitempty"`
}

// MenuOptionVO 菜单选项 VO (用于下拉选择)
type MenuOptionVO struct {
	ID   int64  `json:"value"`
	Name string `json:"label"`
}

// PageResult 分页结果
type PageResult struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

// CommonResponse 通用响应
type CommonResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
