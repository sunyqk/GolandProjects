package model

import "time"

func (TaskNewRecord) TableName() string {
	return "task_new_records"
}

type TaskNewRecord struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	TaskName  string     `json:"task_name" gorm:"size:200;not null"`
	Status    string     `json:"status" gorm:"size:20;default:TODO"`
	Priority  string     `json:"priority" gorm:"size:20;default:MEDIUM"`
	Assignee  string     `json:"assignee" gorm:"size:50"`
	DueDate   *time.Time `json:"due_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
