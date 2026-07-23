package service

import (
	"GolandProjects/model"
	"GolandProjects/repository"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ==================== 导出 ====================

// ExportTasksToExcel 根据条件导出任务到 Excel
func ExportTasksToExcel(query *model.TaskExportQuery) (*excelize.File, error) {
	// 1. 根据条件查询数据
	tasks, err := repository.GetTasksByCondition(query)
	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %w", err)
	}

	// 2. 创建 Excel 文件
	f := excelize.NewFile()
	sheetName := "任务列表"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1") // 删除默认 Sheet

	// 3. 定义表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// 4. 写入表头
	headers := []string{"序号", "任务名称", "状态", "优先级", "负责人", "截止日期", "创建时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// 5. 写入数据行
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D9D9D9", Style: 1},
			{Type: "top", Color: "D9D9D9", Style: 1},
			{Type: "bottom", Color: "D9D9D9", Style: 1},
			{Type: "right", Color: "D9D9D9", Style: 1},
		},
	})

	for rowIdx, task := range tasks {
		row := rowIdx + 2 // 第 1 行是表头
		values := []interface{}{
			rowIdx + 1,
			task.TaskName,
			task.Status,
			task.Priority,
			task.Assignee,
			formatDate(task.DueDate),
			task.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheetName, cell, val)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// 6. 设置列宽
	colWidths := []float64{8, 30, 12, 10, 12, 15, 22}
	for i, width := range colWidths {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, colName, colName, width)
	}

	// 7. 冻结首行
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	return f, nil
}

// ==================== 导入 ====================

// ImportTasksFromExcel 从 Excel 文件导入任务
func ImportTasksFromExcel(file io.Reader) (*model.TaskImportResult, error) {
	// 1. 打开 Excel 文件
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("无法解析 Excel 文件: %w", err)
	}
	defer f.Close()

	// 2. 获取第一个 Sheet
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("Excel 文件中没有工作表")
	}

	// 3. 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel 文件中没有数据（至少需要表头 + 1 行数据）")
	}

	// 4. 解析表头，建立列名 → 索引映射
	headerMap := make(map[string]int)
	for colIdx, header := range rows[0] {
		headerMap[strings.TrimSpace(header)] = colIdx
	}

	// 验证必要列是否存在
	requiredCols := []string{"任务名称", "状态"}
	for _, col := range requiredCols {
		if _, ok := headerMap[col]; !ok {
			return nil, fmt.Errorf("缺少必要列: %s", col)
		}
	}

	// 5. 逐行解析数据
	result := &model.TaskImportResult{
		TotalCount: len(rows) - 1, // 减去表头
		Errors:     make([]string, 0),
	}

	var validTasks []model.TaskNewRecord

	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		lineNum := rowIdx + 1 // Excel 行号（从 1 开始）

		// 跳过空行
		if isEmptyRow(row) {
			result.TotalCount--
			continue
		}

		// 解析并校验
		task, errMsg := parseTaskRow(row, headerMap, lineNum)
		if errMsg != "" {
			result.FailCount++
			result.Errors = append(result.Errors, errMsg)
			continue
		}

		validTasks = append(validTasks, *task)
		result.SuccessCount++
	}

	// 6. 批量写入数据库
	if len(validTasks) > 0 {
		if err := repository.BatchCreateTasks2(validTasks); err != nil {
			return nil, fmt.Errorf("批量写入数据库失败: %w", err)
		}
	}

	return result, nil
}

// ==================== 辅助函数 ====================

// parseTaskRow 解析单行数据为 TaskNewRecord
func parseTaskRow(row []string, headerMap map[string]int, lineNum int) (*model.TaskNewRecord, string) {
	getCell := func(colName string) string {
		if idx, ok := headerMap[colName]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	taskName := getCell("任务名称")
	if taskName == "" {
		return nil, fmt.Sprintf("第 %d 行: 任务名称不能为空", lineNum)
	}

	status := getCell("状态")
	if status == "" {
		status = "TODO" // 默认值
	}
	// 校验状态合法性
	validStatuses := map[string]bool{
		"TODO": true, "IN_PROGRESS": true, "DONE": true, "SUCCESS": true,
	}
	if !validStatuses[strings.ToUpper(status)] {
		return nil, fmt.Sprintf("第 %d 行: 无效状态 '%s'（可选: TODO/IN_PROGRESS/DONE）", lineNum, status)
	}

	priority := getCell("优先级")
	if priority == "" {
		priority = "MEDIUM"
	}

	// 解析截止日期
	var dueDate *time.Time
	if dateStr := getCell("截止日期"); dateStr != "" {
		parsed, err := parseDate(dateStr)
		if err != nil {
			return nil, fmt.Sprintf("第 %d 行: 日期格式错误 '%s'（应为 2006-01-02）", lineNum, dateStr)
		}
		dueDate = &parsed
	}

	task := &model.TaskNewRecord{
		TaskName:  taskName,
		Status:    strings.ToUpper(status),
		Priority:  strings.ToUpper(priority),
		Assignee:  getCell("负责人"),
		DueDate:   dueDate,
		CreatedAt: time.Now(),
	}

	return task, ""
}

// parseDate 支持多种日期格式
func parseDate(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"2006-01-02 15:04:05",
		"20060102",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", s)
}

// formatDate 格式化日期指针
func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// isEmptyRow 判断是否为空行
func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
