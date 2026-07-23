package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是统一的成功响应结构
type Response struct {
	Code int         `json:"code"`           // 业务状态码，0 表示成功
	Msg  string      `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 业务数据，omitempty 表示如果为空则不显示该字段
}

// ErrorResponse 是统一的错误响应结构
type ErrorResponse struct {
	Code int    `json:"code"` // 业务错误码
	Msg  string `json:"msg"`  // 错误信息
}

// Success 封装成功响应
// 它会自动使用 HTTP 200 状态码，并将业务码设为 0
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

// Fail 封装失败响应
// 你可以自定义 HTTP 状态码和业务错误码
func Fail(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, ErrorResponse{
		Code: code,
		Msg:  msg,
	})
}

// PageData 通用的分页数据响应结构
type PageData struct {
	List     interface{} `json:"list"`     // 数据列表
	Total    int64       `json:"total"`    // 总记录数
	Page     int         `json:"page"`     // 当前页码
	PageSize int         `json:"pageSize"` // 每页条数
}
