package model

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password     string         `gorm:"size:255;not null" json:"-"` // 不返回给前端
	RealName     string         `gorm:"size:50;not null" json:"real_name"`
	Email        string         `gorm:"size:100" json:"email"`
	Phone        string         `gorm:"size:20" json:"phone"`
	Department   string         `gorm:"size:100" json:"department"`
	Position     string         `gorm:"size:50" json:"position"`
	EmployeeNo   string         `gorm:"size:50" json:"employee_no"`
	JoinDate     *time.Time     `json:"join_date"`
	ManagerID    *uint           `json:"manager_id"` // 直属上级ID
	Avatar       string          `gorm:"size:255" json:"avatar"`
	Status       string          `gorm:"size:20;default:'normal'" json:"status"` // normal, disabled, deleted
	LastLoginAt  *time.Time      `json:"last_login_at"`
	LastLoginIP  string          `gorm:"size:50" json:"last_login_ip"`
	PasswordChangedAt *time.Time  `json:"password_changed_at"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`

	// 关联关系
	Roles        []Role          `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	Permissions  []Permission    `gorm:"many2many:user_permissions;" json:"permissions,omitempty"`
	Manager      *User           `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Subordinates []User          `gorm:"foreignKey:ManagerID" json:"subordinates,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

