package user

import (
	"time"
)

// User 用户领域模型（对应数据库 user_base 表）
type User struct {
	// ID 主键，自增
	ID uint64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`

	// Username 用户名，唯一索引，长度 1-32
	Username string `gorm:"column:username;uniqueIndex;not null;size:32" json:"username" binding:"required,min=1,max=32"`

	// Password 密码（数据库中必须存储加密后的哈希值，绝不能存明文）
	Password string `gorm:"column:password;not null;size:128" json:"-" binding:"required,min=6,max=128"`

	// Email 邮箱，唯一索引
	Email string `gorm:"column:email;uniqueIndex;size:100" json:"email" binding:"required,email"`

	// Phone 手机号（可选）
	Phone string `gorm:"column:phone;size:20" json:"phone,omitempty"`

	// Avatar 头像 URL（可选）
	Avatar string `gorm:"column:avatar;size:255" json:"avatar,omitempty"`

	// Status 账号状态：1-正常, 0-禁用
	Status int8 `gorm:"column:status;default:1;not null" json:"status"`

	// CreatedAt 创建时间（由 GORM 自动填充）
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	// UpdatedAt 更新时间（由 GORM 自动填充）
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定数据库表名（遵循大型项目规范，使用统一前缀）
func (u *User) TableName() string {
	return "user_base"
}
