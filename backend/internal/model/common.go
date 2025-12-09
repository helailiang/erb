package model

import (
	"time"
	"gorm.io/gorm"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	Action      string    `gorm:"size:50;not null" json:"action"` // 操作类型
	Module      string    `gorm:"size:50" json:"module"` // 操作模块
	Resource    string    `gorm:"size:100" json:"resource"` // 操作对象
	Result      string    `gorm:"size:20" json:"result"` // success, failure
	IP          string    `gorm:"size:50" json:"ip"`
	UserAgent   string    `gorm:"size:500" json:"user_agent"`
	Details     string    `gorm:"type:text" json:"details"` // 操作详情（JSON格式）
	CreatedAt   time.Time `json:"created_at"`

	// 关联关系
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// Notification 通知模型
type Notification struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	Type        string         `gorm:"size:20;not null" json:"type"` // system, work, security
	Title       string         `gorm:"size:200;not null" json:"title"`
	Content     string         `gorm:"type:text" json:"content"`
	IsRead      bool           `gorm:"default:false" json:"is_read"`
	ReadAt      *time.Time     `json:"read_at"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Key         string         `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value       string         `gorm:"type:text" json:"value"`
	Description string         `gorm:"size:255" json:"description"`
	Type        string         `gorm:"size:20;default:'string'" json:"type"` // string, number, boolean, json
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// FileUpload 文件上传模型
type FileUpload struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"`
	FileName    string         `gorm:"size:255;not null" json:"file_name"`
	FilePath    string         `gorm:"size:500;not null" json:"file_path"`
	FileSize    int64          `json:"file_size"`
	FileType    string         `gorm:"size:50" json:"file_type"`
	MimeType    string         `gorm:"size:100" json:"mime_type"`
	Hash        string         `gorm:"size:64" json:"hash"` // 文件MD5哈希
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (FileUpload) TableName() string {
	return "file_uploads"
}

