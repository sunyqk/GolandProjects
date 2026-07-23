package response

// PageResult 通用分页返回结构体
type PageResult struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页条数
}
