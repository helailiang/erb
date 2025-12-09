package model

import (
	"time"
	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"uniqueIndex;size:100;not null" json:"code"` // 权限编码，如：user:create
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Type        string         `gorm:"size:20;default:'function'" json:"type"` // function, data, operation
	Module      string         `gorm:"size:50" json:"module"` // 所属模块
	Status      string         `gorm:"size:20;default:'enabled'" json:"status"` // enabled, disabled
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Users []User `gorm:"many2many:user_permissions;" json:"users,omitempty"`
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID     uint      `gorm:"primaryKey" json:"role_id"`
	PermissionID uint    `gorm:"primaryKey" json:"permission_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserPermission 用户权限关联表（直接分配）
type UserPermission struct {
	UserID       uint      `gorm:"primaryKey" json:"user_id"`
	PermissionID uint      `gorm:"primaryKey" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserPermission) TableName() string {
	return "user_permissions"
}

