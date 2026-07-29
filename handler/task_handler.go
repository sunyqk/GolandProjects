package handler

import (
	"GolandProjects/model"
	"GolandProjects/pkg/response"
	"GolandProjects/repository"
	"GolandProjects/service"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
)

// 处理前端发来的 HTTP 请求
func HandleTaskAPI(c *gin.Context) {
	count := 5
	if q := c.Query("count"); q != "" {
		fmt.Sscanf(q, "%d", &count)
	}

	// 调用 Service 层处理业务
	results, err := service.ProcessAndSaveTasks(count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "处理失败"})
		return
	}

	// 将数据以 JSON 格式返回给前端
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": results,
	})
}

// CreateTaskAPI 创建任务 (POST)
func CreateTaskAPI(c *gin.Context) {
	var task model.TaskRecord
	// 1. 解析 JSON Body
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误", "error": err.Error()})
		return
	}
	// 2. 调用 Service 或直接调用 Repository
	if err := repository.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建成功", "data": task})
}

// UpdateTaskAPI 更新任务 (PUT)
func UpdateTaskAPI(c *gin.Context) {
	// 1. 获取路径中的 ID
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	// 2. 解析 JSON Body
	var task model.TaskRecord
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	task.ID = uint(id) // 将路径 ID 赋给结构体

	// 3. 更新数据库
	if err := repository.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// GetTaskPageAPI 对应 main.go 第19行
func GetTaskPageAPI(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("pageSize", "10")

	taskName := c.Query("taskName")
	createAt := c.QueryArray("createAt")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(sizeStr)

	// 假设 repository 中有 GetTaskList 方法
	tasks, total, err := repository.GetTaskList(page, pageSize, taskName, createAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":     tasks,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// DeleteTaskAPI 对应 main.go 第22行
func DeleteTaskAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID格式错误"})
		return
	}

	// 假设 repository 中有 DeleteTask 方法
	err = repository.DeleteTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// 测试统一封装返回数据结构接口
func GetTaskResponsePageAPI(c *gin.Context) {
	// 1. 解析分页参数
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("pageSize", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(sizeStr)

	// 2. 解析查询条件
	taskName := c.Query("taskName")
	createAt := c.QueryArray("createAt")

	// 3. 调用 Repository 获取数据
	tasks, total, err := repository.GetTaskList(page, pageSize, taskName, createAt)
	if err != nil {
		// 失败时使用统一的错误返回
		response.Fail(c, http.StatusInternalServerError, 50001, "查询任务列表失败")
		return
	}

	// 4. 组装分页数据并使用统一格式返回
	pageData := response.PageData{
		List:     tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	response.Success(c, pageData)
}

// 定义一个全局的错误变量，用于表示记录未找到
var ErrRecordNotFound = errors.New("record not found")

// GetTaskDetailAPI 获取单条任务详情 (GET)
func GetTaskDetailAPI(c *gin.Context) {
	// 1. 获取并校验路径参数 ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	// 校验：转换失败 或 ID <= 0
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID，ID必须为正整数",
		})
		return
	}

	// 2. 调用 Repository 层查询数据
	task, err := repository.GetTaskByID(uint(id))
	if err != nil {
		// 【推荐】使用 errors.Is 进行比较，兼容性更好
		if errors.Is(err, ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "任务不存在",
			})
		} else {
			// 【安全】生产环境不要暴露 err.Error()
			// 建议在控制台打印真实错误，给前端返回通用提示
			fmt.Printf("查询任务详情失败 [ID:%d]: %v\n", id, err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "服务器内部错误，请稍后重试",
			})
		}
		return
	}

	// 3. 查询成功，返回数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"data": task,
	})
}

// UploadFileAPI 处理文件上传
func UploadFileAPI(c *gin.Context) {
	// 1. 获取表单中的文件，"file" 是前端 input 的 name 属性
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "获取文件失败"})
		return
	}

	// 2. 定义保存路径 (建议在生产环境中按日期或用户ID分目录)
	dst := "./uploads/" + file.Filename

	// 3. 保存文件到指定目录
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "文件保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "上传成功",
		"data": gin.H{"filename": file.Filename},
	})
}

// DownloadFileAPI 处理文件下载
func DownloadFileAPI(c *gin.Context) {
	// 1. 从路径参数获取文件名
	filename := c.Param("filename")
	filePath := "./uploads/" + filename

	// 2. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件不存在"})
		return
	}

	// 3. 使用 Gin 的 FileAttachment 方法，强制浏览器下载
	// 第二个参数是浏览器下载时显示的默认文件名
	c.FileAttachment(filePath, filename)
}

// //PreviewFileAPI 处理文件在线预览
//func PreviewFileAPI(c *gin.Context) {
//	filename := c.Param("filename")
//	filePath := "./uploads/" + filename
//
//	// 检查文件是否存在
//	if _, err := os.Stat(filePath); os.IsNotExist(err) {
//		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件不存在"})
//		return
//	}
//
//	// 1. 设置响应头，告诉浏览器内联显示（inline）而不是下载（attachment）
//	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filename))
//
//	// 2. 根据文件类型设置 Content-Type (以图片为例)
//	// 如果是 PDF，则设置为 "application/pdf"
//	c.Header("Content-Type", "image/png")
//
//	// 3. 返回文件内容
//	c.File(filePath)
//}

// PreviewFileAPI 处理文件在线预览
func PreviewFileAPI(c *gin.Context) {
	filename := c.Param("filename")
	filePath := "./uploads/" + filename

	// 1. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "文件不存在"})
		return
	}

	// 2. 【关键】动态获取 Content-Type
	// mime.TypeByExtension 会根据 .png, .pdf, .txt 等后缀返回对应的类型
	contentType := mime.TypeByExtension(filepath.Ext(filename))

	// 兜底策略：如果无法识别后缀，尝试通过读取文件内容的前几个字节来判断
	if contentType == "" {
		file, _ := os.Open(filePath)
		defer file.Close()
		buffer := make([]byte, 512)
		file.Read(buffer)
		contentType = http.DetectContentType(buffer)
	}

	// 3. 设置响应头并输出文件
	// inline 告诉浏览器/Postman：请在当前页面内嵌显示，不要下载
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", filename))
	c.Header("Content-Type", contentType)
	c.File(filePath)
}
