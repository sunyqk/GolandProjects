package main

import (
	"GolandProjects/config"
	"GolandProjects/handler"
	"GolandProjects/job"
	"GolandProjects/pkg/logger"
	"GolandProjects/router"
	"fmt"
)

func main() {

	// 1. 初始化数据库
	config.InitDB()

	// 2. 初始化路由
	//r := gin.Default()
	// 自定义路由,和下面的接口不冲突
	r := router.SetupRouter()

	// 配置项目链接本地或其他ip地址
	// 设置可信代理列表
	// 请将 "192.168.1.2" 和 "10.0.0.0/8" 替换为你实际的代理服务器IP或网段
	//err := r.SetTrustedProxies([]string{"192.168.4.122", "10.0.0.0/8"})
	err := r.SetTrustedProxies([]string{"192.168.4.122", "127.0.0.1", "::1"})
	if err != nil {
		// 处理错误，例如IP格式不正确
		panic(err)
	}

	// 初始化日志开发环境 dev localmaster master 这三种环境
	logger.Init("localmaster")
	defer logger.Log.Sync()

	// 2. 启动后台定时任务
	// 因为它内部开启了 goroutine，所以这里会立即返回，不会卡住程序
	job.StartTaskMonitor()
	fmt.Println("系统服务已启动，后台监控运行中...")

	// 3. 注册接口
	// 分页查询列表
	r.GET("/api/tasks/page", handler.GetTaskPageAPI)
	// 根据ID查单条
	r.GET("/api/tasks/detail/:id", handler.GetTaskDetailAPI)
	// 创建任务 POST
	r.POST("/api/tasks/create", handler.CreateTaskAPI)
	// 更新 PUT
	r.PUT("/api/tasks/update/:id", handler.UpdateTaskAPI)
	// 删除 DELETE
	r.DELETE("/api/tasks/delete/:id", handler.DeleteTaskAPI)

	// 文件上传、下载、预览路由
	r.POST("/api/files/upload", handler.UploadFileAPI)              // 上传文件
	r.GET("/api/files/download/:filename", handler.DownloadFileAPI) // 下载文件
	r.GET("/api/files/preview/:filename", handler.PreviewFileAPI)   // 预览文件
	// 【可选】如果你希望直接通过 URL 访问 uploads 目录下的所有文件（例如 http://localhost:8080/uploads/xxx.png）
	// 可以加上这一行静态资源映射，非常方便前端展示图片
	r.Static("/uploads", "./uploads")

	// 4. 启动服务
	r.Run(":8080")
}
