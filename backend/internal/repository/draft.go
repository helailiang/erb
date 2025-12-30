package repository

import (
	"erb/internal/model"
	"gorm.io/gorm"
	"time"
)

// DraftRepository 草稿数据访问层
type DraftRepository struct {
	db *gorm.DB
}

// NewDraftRepository 创建草稿仓库
func NewDraftRepository(db *gorm.DB) *DraftRepository {
	return &DraftRepository{db: db}
}

// CreateOrUpdate 创建或更新草稿
func (r *DraftRepository) CreateOrUpdate(draft *model.ReportDraft) error {
	var existing model.ReportDraft
	err := r.db.Where("user_id = ? AND DATE(report_date) = DATE(?)", draft.UserID, draft.ReportDate).
		First(&existing).Error

	if err == nil {
		// 更新现有草稿
		existing.Content = draft.Content
		existing.Title = draft.Title
		existing.TemplateID = draft.TemplateID
		existing.Version++
		existing.UpdatedAt = time.Now()
		return r.db.Save(&existing).Error
	} else if err == gorm.ErrRecordNotFound {
		// 创建新草稿
		return r.db.Create(draft).Error
	}
	return err
}

// GetByUserAndDate 根据用户和日期获取草稿
func (r *DraftRepository) GetByUserAndDate(userID uint, date time.Time) (*model.ReportDraft, error) {
	var draft model.ReportDraft
	dateStr := date.Format("2006-01-02")
	err := r.db.Where("user_id = ? AND DATE(report_date) = ?", userID, dateStr).First(&draft).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

// Delete 删除草稿
func (r *DraftRepository) Delete(userID uint, date time.Time) error {
	dateStr := date.Format("2006-01-02")
	return r.db.Where("user_id = ? AND DATE(report_date) = ?", userID, dateStr).Delete(&model.ReportDraft{}).Error
}

// List 获取用户的草稿列表
func (r *DraftRepository) List(userID uint, page, pageSize int) ([]model.ReportDraft, int64, error) {
	var drafts []model.ReportDraft
	var total int64

	query := r.db.Model(&model.ReportDraft{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("report_date DESC").
		Offset(offset).Limit(pageSize).Find(&drafts).Error
	return drafts, total, err
}



