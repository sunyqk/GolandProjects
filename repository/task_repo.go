package repository

import (
	"GolandProjects/config"
	"GolandProjects/model"
	"GolandProjects/response"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CreateTask 创建任务
func CreateTask(task *model.TaskRecord) error {
	return config.DB.Create(task).Error
}

// GetTaskByID 根据ID查询单条任务
// 找不到数据返回 nil,nil；数据库异常返回 err
func GetTaskByID(id uint) (*model.TaskRecord, error) {
	var task model.TaskRecord
	err := config.DB.First(&task, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// TaskQuery 分页查询条件结构体，扩展筛选条件用
type TaskQuery struct {
	Page     int
	PageSize int
	Status   int      // 任务状态筛选
	Name     string   // 任务名称模糊搜索
	CreateAt []string `form:"createAt"` // 接收数组，例如: createAt=2026-01-01&createAt=2026-12-31
}

// GetTaskPage 分页查询，整合条件+分页，直接返回分页结果
func GetTaskPage(query TaskQuery) (*response.PageResult, error) {
	db := config.DB.Model(&model.TaskRecord{})

	// 条件过滤
	if query.Status != 0 {
		db = db.Where("status = ?", query.Status)
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}

	// 统计总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页默认值处理
	page := query.Page
	pageSize := query.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 查询列表
	var list []model.TaskRecord
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&list).Error
	if err != nil {
		return nil, err
	}

	return &response.PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// UpdateTask 更新任务
// 注意：Updates 忽略零值，如需全量更新改用 Save
func UpdateTask(task *model.TaskRecord) error {
	return config.DB.Model(task).Updates(task).Error
}

// UpdateTaskAll 全字段更新（包含零值）
func UpdateTaskAll(task *model.TaskRecord) error {
	return config.DB.Save(task).Error
}

// DeleteTaskByID 物理删除任务
func DeleteTaskByID(id uint) error {
	return config.DB.Delete(&model.TaskRecord{}, id).Error
}

// LogicDeleteTask 逻辑删除（推荐，需model自带DeleteAt字段）
func LogicDeleteTask(id uint) error {
	return config.DB.Model(&model.TaskRecord{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// GetTaskList 分页查询任务列表
// page: 当前页码
// pageSize: 每页显示数量
// 返回值: 任务列表, 总记录数, 错误信息
func GetTaskList(page, pageSize int, taskName string, createAt []string) ([]model.TaskRecord, int64, error) {
	var tasks []model.TaskRecord
	var total int64

	// 1. 获取数据库连接，并使用 NewSession 开启一个独立的查询会话（防止并发污染）
	query := config.DB.Model(&model.TaskRecord{})

	// 2. 如果传入了 taskName 且不为空，则追加模糊查询条件
	if taskName != "" {
		query = query.Where("task_name LIKE ?", "%"+taskName+"%")
	}

	// 处理时间范围数组
	switch len(createAt) {
	case 2:
		// 传了完整的开始和结束时间
		query = query.Where("create_at >= ? AND create_at <= ?", createAt[0], createAt[1])
	case 1:
		// 只传了一个时间，默认作为开始时间
		query = query.Where("create_at >= ?", createAt[0])
	}

	// 3. 统计总数（注意：这里使用的是带有条件的 query）
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. 分页查询（同样使用带有条件的 query）
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetTasksByCondition 根据条件查询任务列表（用于导出）
func GetTasksByConditionold(query *model.TaskExportQuery) ([]model.TaskRecord, error) {
	var tasks []model.TaskRecord
	db := config.DB.Model(&model.TaskRecord{})

	// 动态拼接查询条件
	if query.TaskName != "" {
		db = db.Where("task_name LIKE ?", "%"+query.TaskName+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Priority != "" {
		db = db.Where("priority = ?", query.Priority)
	}
	if query.Assignee != "" {
		db = db.Where("assignee LIKE ?", "%"+query.Assignee+"%")
	}
	if query.StartDate != "" {
		start, err := time.Parse("2006-01-02", query.StartDate)
		if err == nil {
			db = db.Where("created_at >= ?", start)
		}
	}
	if query.EndDate != "" {
		end, err := time.Parse("2006-01-02", query.EndDate)
		if err == nil {
			// 包含当天
			db = db.Where("created_at <= ?", end.Add(24*time.Hour-time.Second))
		}
	}

	if err := db.Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// BatchCreateTasks 批量创建任务（用于导入）
func BatchCreateTasksold(tasks []model.TaskRecord) error {
	if len(tasks) == 0 {
		return nil
	}
	return config.DB.CreateInBatches(tasks, 100).Error
}

func BatchCreateTasks2(tasks []model.TaskNewRecord) error {
	if len(tasks) == 0 {
		return nil
	}
	return config.DB.CreateInBatches(tasks, 100).Error
}

// repository —— 注意返回类型和 Find 的接收变量都要改
func GetTasksByCondition(query *model.TaskExportQuery) ([]model.TaskNewRecord, error) {
	//var tasks []model.TaskNewRecord
	//db := config.DB.Model(&model.TaskNewRecord{})
	//// ... 其余 Where 不变
	//if err := db.Order("created_at DESC").Find(&tasks).Error; err != nil {
	//	return nil, err
	//}
	//return tasks, nil

	var tasks []model.TaskNewRecord

	// 1. 初始化数据库连接并指定模型
	db := config.DB.Model(&model.TaskNewRecord{})

	// 2. 【核心修改】添加 task_name 模糊查询
	// 只有当 query.TaskName 不为空时才执行过滤
	if query.TaskName != "" {
		// 使用 LIKE %keyword% 实现模糊匹配
		db = db.Where("task_name LIKE ?", "%"+query.TaskName+"%")
	}

	// 3. 排序并查找
	// 注意：这里链式调用 db，确保上面的 Where 条件生效
	if err := db.Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func BatchCreateTasks(tasks []model.TaskNewRecord) error {
	if len(tasks) == 0 {
		return nil
	}
	return config.DB.CreateInBatches(tasks, 100).Error
}
