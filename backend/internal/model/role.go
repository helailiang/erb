package model

import (
	"time"
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Status      string         `gorm:"size:20;default:'enabled'" json:"status"` // enabled, disabled
	IsSystem    bool           `gorm:"default:false" json:"is_system"` // 是否为系统预置角色
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

