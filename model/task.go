package model

import "time"

// 定义数据库表结构
type TaskRecord struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskName  string    `json:"task_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// PageResult 用于封装分页查询的返回结果
type PageResult struct {
	List     interface{} `json:"list"`     // 数据列表
	Total    int64       `json:"total"`    // 总记录数
	Page     int         `json:"page"`     // 当前页码
	PageSize int         `json:"pageSize"` // 每页条数
}
