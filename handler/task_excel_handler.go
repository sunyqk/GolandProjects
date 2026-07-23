package handler

import (
	"GolandProjects/model"
	"GolandProjects/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ExportTasks 导出任务 Excel
// GET /api/v1/tasks/export?task_name=&status=&start_date=&end_date=
func ExportTasks(c *gin.Context) {
	// 1. 绑定查询参数
	var query model.TaskExportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 2. 调用 Service 生成 Excel
	f, err := service.ExportTasksToExcel(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "导出失败: " + err.Error(),
		})
		return
	}
	defer f.Close()

	// 3. 设置响应头，返回文件流
	fileName := fmt.Sprintf("任务导出_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "写入文件失败",
		})
	}
}

// ImportTasks 导入任务 Excel
// POST /api/v1/tasks/import
func ImportTasks(c *gin.Context) {
	// 1. 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请上传文件（字段名: file）",
		})
		return
	}
	defer file.Close()

	// 2. 校验文件类型
	if header.Size > 10*1024*1024 { // 10MB 限制
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "文件大小不能超过 10MB",
		})
		return
	}

	ext := header.Filename[len(header.Filename)-5:]
	if ext != ".xlsx" && ext != ".xls" {
		// 更严谨的校验
		if len(header.Filename) < 5 ||
			(header.Filename[len(header.Filename)-5:] != ".xlsx" &&
				header.Filename[len(header.Filename)-4:] != ".xls") {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "仅支持 .xlsx 或 .xls 格式",
			})
			return
		}
	}

	// 3. 调用 Service 解析并导入
	result, err := service.ImportTasksFromExcel(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "导入失败: " + err.Error(),
		})
		return
	}

	// 4. 返回导入结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "导入完成",
		"data":    result,
	})
}

// DownloadTemplate 下载导入模板
// GET /api/v1/tasks/template
func DownloadTemplate(c *gin.Context) {
	f, err := service.GenerateImportTemplate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成模板失败",
		})
		return
	}
	defer f.Close()

	fileName := "任务导入模板.xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	f.Write(c.Writer)
}
