package repository

import (
	"erb/internal/model"
	"gorm.io/gorm"
	"time"
)

// ReportRepository 日报数据访问层
type ReportRepository struct {
	db *gorm.DB
}

// NewReportRepository 创建日报仓库
func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// Create 创建日报
func (r *ReportRepository) Create(report *model.DailyReport) error {
	return r.db.Create(report).Error
}

// GetByID 根据ID获取日报
func (r *ReportRepository) GetByID(id uint) (*model.DailyReport, error) {
	var report model.DailyReport
	err := r.db.Preload("User").Preload("Template").
		Preload("Approvals.Approver").Preload("Comments.User").
		Preload("Attachments").First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// Update 更新日报
func (r *ReportRepository) Update(report *model.DailyReport) error {
	return r.db.Save(report).Error
}

// Delete 删除日报（软删除）
func (r *ReportRepository) Delete(id uint) error {
	return r.db.Delete(&model.DailyReport{}, id).Error
}

// List 获取日报列表
func (r *ReportRepository) List(page, pageSize int, conditions map[string]interface{}) ([]model.DailyReport, int64, error) {
	var reports []model.DailyReport
	var total int64

	query := r.db.Model(&model.DailyReport{}).Preload("User")

	// 应用筛选条件
	for key, value := range conditions {
		if key == "user_id" {
			query = query.Where("user_id = ?", value)
		} else if key == "status" {
			query = query.Where("status = ?", value)
		} else if key == "report_date_start" {
			query = query.Where("report_date >= ?", value)
		} else if key == "report_date_end" {
			query = query.Where("report_date <= ?", value)
		} else if key == "department" {
			query = query.Joins("JOIN users ON users.id = daily_reports.user_id").
				Where("users.department = ?", value)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("report_date DESC, created_at DESC").
		Offset(offset).Limit(pageSize).Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// GetByUserAndDate 根据用户和日期获取日报
func (r *ReportRepository) GetByUserAndDate(userID uint, date time.Time) (*model.DailyReport, error) {
	var report model.DailyReport
	dateStr := date.Format("2006-01-02")
	err := r.db.Where("user_id = ? AND DATE(report_date) = ?", userID, dateStr).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// GetUserReports 获取用户的日报列表
func (r *ReportRepository) GetUserReports(userID uint, page, pageSize int, startDate, endDate *time.Time) ([]model.DailyReport, int64, error) {
	var reports []model.DailyReport
	var total int64

	query := r.db.Model(&model.DailyReport{}).Where("user_id = ?", userID)

	if startDate != nil {
		query = query.Where("report_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("report_date <= ?", endDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("report_date DESC").
		Offset(offset).Limit(pageSize).Find(&reports).Error
	if err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// CreateApproval 创建审批记录
func (r *ReportRepository) CreateApproval(approval *model.ReportApproval) error {
	return r.db.Create(approval).Error
}

// UpdateApproval 更新审批记录
func (r *ReportRepository) UpdateApproval(approval *model.ReportApproval) error {
	return r.db.Save(approval).Error
}

// GetApprovalsByReportID 获取日报的审批记录
func (r *ReportRepository) GetApprovalsByReportID(reportID uint) ([]model.ReportApproval, error) {
	var approvals []model.ReportApproval
	err := r.db.Where("report_id = ?", reportID).
		Preload("Approver").Order("created_at DESC").Find(&approvals).Error
	return approvals, err
}

// CreateComment 创建评论
func (r *ReportRepository) CreateComment(comment *model.ReportComment) error {
	return r.db.Create(comment).Error
}

// GetCommentsByReportID 获取日报的评论
func (r *ReportRepository) GetCommentsByReportID(reportID uint) ([]model.ReportComment, error) {
	var comments []model.ReportComment
	err := r.db.Where("report_id = ? AND parent_id IS NULL", reportID).
		Preload("User").Preload("Replies.User").
		Order("created_at ASC").Find(&comments).Error
	return comments, err
}



