package model

import "time"

// TaskExcelRow Excel 导入导出的行结构
type TaskExcelRow struct {
	TaskName  string    `json:"task_name"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	Assignee  string    `json:"assignee"`
	DueDate   string    `json:"due_date"`
	CreatedAt time.Time `json:"created_at"`
}

// TaskExportQuery 导出条件查询参数
type TaskExportQuery struct {
	TaskName  string `form:"task_name"`  // 任务名称（模糊）
	Status    string `form:"status"`     // 状态筛选
	Priority  string `form:"priority"`   // 优先级筛选
	StartDate string `form:"start_date"` // 开始日期 2026-01-01
	EndDate   string `form:"end_date"`   // 结束日期 2026-12-31
	Assignee  string `form:"assignee"`   // 负责人
}

// TaskImportResult 导入结果
type TaskImportResult struct {
	TotalCount   int      `json:"total_count"`   // 总行数
	SuccessCount int      `json:"success_count"` // 成功数
	FailCount    int      `json:"fail_count"`    // 失败数
	Errors       []string `json:"errors"`        // 错误详情
}
