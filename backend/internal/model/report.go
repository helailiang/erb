package model

import (
	"time"

	"gorm.io/gorm"
)

// DailyReport 日报模型
type DailyReport struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	ReportDate  time.Time      `gorm:"index;not null" json:"report_date"` // 日报日期
	TemplateID  *uint          `json:"template_id"`                       // 使用的模板ID
	Title       string         `gorm:"size:200" json:"title"`
	Content     string         `gorm:"type:text" json:"content"`              // 日报内容（富文本）
	Status      string         `gorm:"size:20;default:'draft'" json:"status"` // draft, submitted, approved, rejected
	IsLate      bool           `gorm:"default:false" json:"is_late"`          // 是否补交
	LateReason  string         `gorm:"size:255" json:"late_reason"`           // 补交原因
	SubmittedAt *time.Time     `json:"submitted_at"`                          // 提交时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User        User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Template    *ReportTemplate    `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Approvals   []ReportApproval   `gorm:"foreignKey:ReportID" json:"approvals,omitempty"`
	Comments    []ReportComment    `gorm:"foreignKey:ReportID" json:"comments,omitempty"`
	Attachments []ReportAttachment `gorm:"foreignKey:ReportID" json:"attachments,omitempty"`
}

// TableName 指定表名
func (DailyReport) TableName() string {
	return "daily_reports"
}

// ReportTemplate 日报模板模型
type ReportTemplate struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	Type        string         `gorm:"size:20;default:'standard'" json:"type"`  // standard, project, sales, custom
	Fields      string         `gorm:"type:text" json:"fields"`                 // JSON格式的字段定义
	Status      string         `gorm:"size:20;default:'enabled'" json:"status"` // enabled, disabled
	IsSystem    bool           `gorm:"default:false" json:"is_system"`          // 是否为系统模板
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Reports []DailyReport `gorm:"foreignKey:TemplateID" json:"reports,omitempty"`
}

// TableName 指定表名
func (ReportTemplate) TableName() string {
	return "report_templates"
}

// ReportApproval 日报审批模型
type ReportApproval struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ReportID      uint       `gorm:"index;not null" json:"report_id"`
	ApproverID    uint       `gorm:"index;not null" json:"approver_id"` // 审批人ID
	Status        string     `gorm:"size:20;not null" json:"status"`    // pending, approved, rejected, transferred
	Comment       string     `gorm:"type:text" json:"comment"`          // 审批意见
	TransferredTo *uint      `json:"transferred_to"`                    // 转交给谁
	ApprovedAt    *time.Time `json:"approved_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// 关联关系
	Report   DailyReport `gorm:"foreignKey:ReportID" json:"report,omitempty"`
	Approver User        `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

// TableName 指定表名
func (ReportApproval) TableName() string {
	return "report_approvals"
}

// ReportComment 日报评论模型
type ReportComment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ReportID    uint           `gorm:"index;not null" json:"report_id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	ParentID    *uint          `json:"parent_id"` // 父评论ID，用于回复
	IsImportant bool           `gorm:"default:false" json:"is_important"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Report  DailyReport     `gorm:"foreignKey:ReportID" json:"report,omitempty"`
	User    User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Parent  *ReportComment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies []ReportComment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// TableName 指定表名
func (ReportComment) TableName() string {
	return "report_comments"
}

// ReportAttachment 日报附件模型
type ReportAttachment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ReportID  uint      `gorm:"index;not null" json:"report_id"`
	FileName  string    `gorm:"size:255;not null" json:"file_name"`
	FilePath  string    `gorm:"size:500;not null" json:"file_path"`
	FileSize  int64     `json:"file_size"`                // 文件大小（字节）
	FileType  string    `gorm:"size:50" json:"file_type"` // 文件类型
	CreatedAt time.Time `json:"created_at"`

	// 关联关系
	Report DailyReport `gorm:"foreignKey:ReportID" json:"report,omitempty"`
}

// TableName 指定表名
func (ReportAttachment) TableName() string {
	return "report_attachments"
}

// ReportDraft 日报草稿模型
type ReportDraft struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	ReportDate time.Time `gorm:"index;not null" json:"report_date"`
	TemplateID *uint     `json:"template_id"`
	Title      string    `gorm:"size:200" json:"title"`
	Content    string    `gorm:"type:text" json:"content"`
	Version    int       `gorm:"default:1" json:"version"` // 草稿版本
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联关系
	User     User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Template *ReportTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
}

// TableName 指定表名
func (ReportDraft) TableName() string {
	return "report_drafts"
}
