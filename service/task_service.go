package service

import (
	"GolandProjects/model"
	"GolandProjects/repository"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

// GetTaskPage 获取任务分页列表
func GetTaskPage(page, pageSize int, taskName string, createAt []string) (*model.PageResult, error) {
	// 1. 参数校验与默认值处理
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // 限制最大每页条数，防止恶意请求
	}

	// 2. 调用 Repository 层获取数据和总数
	tasks, total, err := repository.GetTaskList(page, pageSize, taskName, createAt)
	if err != nil {
		return nil, err
	}

	// 3. 组装返回结果
	return &model.PageResult{
		List:     tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CreateTask 创建任务
func CreateTask(task *model.TaskRecord) error {
	// 这里可以加入业务逻辑，比如检查标题是否重复
	if task.TaskName == "" {
		return errors.New("任务标题不能为空")
	}
	return repository.CreateTask(task)
}

// UpdateTask 更新任务
func UpdateTask(task *model.TaskRecord) error {
	// 可以加入权限校验等逻辑
	return repository.UpdateTask(task)
}

// DeleteTask 删除任务
func DeleteTask(id uint) error {
	return repository.DeleteTaskByID(id)
}

// 处理并发任务并写入数据库
func ProcessAndSaveTasks(taskCount int) ([]model.TaskRecord, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		records []model.TaskRecord
	)

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(2)) * time.Second)

			record := model.TaskRecord{
				TaskName:  fmt.Sprintf("并发任务-%d", id),
				Status:    "SUCCESS",
				CreatedAt: time.Now(),
			}

			// 调用 Repository 层写入数据库
			if err := repository.CreateTask(&record); err != nil {
				fmt.Printf("任务 %d 写入失败: %v\n", id, err)
				return
			}

			mu.Lock()
			records = append(records, record)
			mu.Unlock()
		}(i + 1)
	}

	wg.Wait()
	return records, nil
}

// GenerateImportTemplate 生成标准导入模板
func GenerateImportTemplate() (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "导入模板"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// 表头
	headers := []string{"任务名称*", "状态*", "优先级", "负责人", "截止日期"}
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 11, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// 示例数据
	examples := [][]string{
		{"完成用户模块开发", "TODO", "HIGH", "张三", "2026-08-15"},
		{"编写单元测试", "IN_PROGRESS", "MEDIUM", "李四", "2026-08-20"},
		{"部署上线", "DONE", "HIGH", "王五", "2026-09-01"},
	}
	for rowIdx, row := range examples {
		for colIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	// 说明 Sheet
	f.NewSheet("填写说明")
	instructions := []string{
		"字段说明：",
		"  任务名称*：必填，不超过 100 个字符",
		"  状态*：必填，可选值：TODO / IN_PROGRESS / DONE",
		"  优先级：可选，可选值：LOW / MEDIUM / HIGH（默认 MEDIUM）",
		"  负责人：可选，填写姓名",
		"  截止日期：可选，格式：2006-01-02",
		"",
		"注意事项：",
		"  1. 带 * 号的列为必填项",
		"  2. 请勿修改表头行",
		"  3. 单次导入不超过 5000 行",
		"  4. 文件仅支持 .xlsx 格式",
	}
	for i, line := range instructions {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		f.SetCellValue("填写说明", cell, line)
	}

	// 列宽
	colWidths := []float64{25, 15, 12, 12, 15}
	for i, w := range colWidths {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, col, col, w)
	}

	return f, nil
}
