package router

import (
	"GolandProjects/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 允许跨域（前端下载文件需要）
	r.Use(corsMiddleware())

	api := r.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			//// 原有 CRUD
			//tasks.GET("", handler.GetTaskPage)
			//tasks.POST("", handler.CreateTask)
			//tasks.PUT("/:id", handler.UpdateTask)
			//tasks.DELETE("/:id", handler.DeleteTask)

			// Excel 导入导出
			tasks.GET("/export", handler.ExportTasks)        // 条件导出
			tasks.POST("/import", handler.ImportTasks)       // 文件导入
			tasks.GET("/template", handler.DownloadTemplate) // 下载模板
		}
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
