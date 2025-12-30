package service

import (
	"erb/internal/model"
	"erb/internal/repository"
	"errors"
	"fmt"
	"time"
)

// ReportService 日报服务
type ReportService struct {
	reportRepo *repository.ReportRepository
	draftRepo  *repository.DraftRepository
	userRepo   *repository.UserRepository
}

// NewReportService 创建日报服务
func NewReportService(reportRepo *repository.ReportRepository, draftRepo *repository.DraftRepository, userRepo *repository.UserRepository) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		draftRepo:  draftRepo,
		userRepo:   userRepo,
	}
}

// CreateReportRequest 创建日报请求
type CreateReportRequest struct {
	ReportDate string `json:"report_date" binding:"required"` // 格式: YYYY-MM-DD
	TemplateID *uint  `json:"template_id"`
	Title      string `json:"title"`
	Content    string `json:"content" binding:"required"`
	IsLate     bool   `json:"is_late"`
	LateReason string `json:"late_reason"`
}

// UpdateReportRequest 更新日报请求
type UpdateReportRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ListReportRequest 日报列表请求
type ListReportRequest struct {
	Page       int        `form:"page" binding:"min=1"`
	PageSize   int        `form:"page_size" binding:"min=1,max=100"`
	UserID     *uint      `form:"user_id"`
	Status     string     `form:"status"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
	Department string     `form:"department"`
}

// SaveDraftRequest 保存草稿请求
type SaveDraftRequest struct {
	ReportDate string `json:"report_date" binding:"required"` // 格式: YYYY-MM-DD
	TemplateID *uint  `json:"template_id"`
	Content    string `json:"content"`
	Title      string `json:"title"`
}

// CreateReport 创建日报
func (s *ReportService) CreateReport(userID uint, req *CreateReportRequest) (*model.DailyReport, error) {
	// 解析日期字符串
	reportDate, err := time.Parse("2006-01-02", req.ReportDate)
	if err != nil {
		return nil, errors.New("日期格式错误，请使用YYYY-MM-DD格式")
	}

	// 检查该日期是否已有日报
	existing, _ := s.reportRepo.GetByUserAndDate(userID, reportDate)
	if existing != nil {
		return nil, errors.New("该日期已提交日报")
	}

	now := time.Now()
	report := &model.DailyReport{
		UserID:      userID,
		ReportDate:  reportDate,
		TemplateID:  req.TemplateID,
		Title:       req.Title,
		Content:     req.Content,
		Status:      "submitted",
		IsLate:      req.IsLate,
		LateReason:  req.LateReason,
		SubmittedAt: &now,
	}

	if err := s.reportRepo.Create(report); err != nil {
		return nil, fmt.Errorf("创建日报失败: %w", err)
	}

	// 删除对应日期的草稿
	_ = s.draftRepo.Delete(userID, reportDate)

	// 创建审批记录（推送给直属上级）
	user, err := s.userRepo.GetByID(userID)
	if err == nil && user.ManagerID != nil {
		approval := &model.ReportApproval{
			ReportID:   report.ID,
			ApproverID: *user.ManagerID,
			Status:     "pending",
		}
		_ = s.reportRepo.CreateApproval(approval)
	}

	// 重新获取日报信息
	report, err = s.reportRepo.GetByID(report.ID)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// UpdateReport 更新日报
func (s *ReportService) UpdateReport(reportID, userID uint, req *UpdateReportRequest) (*model.DailyReport, error) {
	report, err := s.reportRepo.GetByID(reportID)
	if err != nil {
		return nil, errors.New("日报不存在")
	}

	// 检查权限：只能修改自己的日报
	if report.UserID != userID {
		return nil, errors.New("无权修改此日报")
	}

	// 检查状态：已审批的日报不允许修改
	if report.Status == "approved" {
		return nil, errors.New("已审批的日报不允许修改")
	}

	if req.Title != "" {
		report.Title = req.Title
	}
	if req.Content != "" {
		report.Content = req.Content
	}

	if err := s.reportRepo.Update(report); err != nil {
		return nil, fmt.Errorf("更新日报失败: %w", err)
	}

	// 重新获取日报信息
	report, err = s.reportRepo.GetByID(reportID)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// DeleteReport 删除日报
func (s *ReportService) DeleteReport(reportID, userID uint) error {
	report, err := s.reportRepo.GetByID(reportID)
	if err != nil {
		return errors.New("日报不存在")
	}

	// 检查权限：只能删除自己的日报
	if report.UserID != userID {
		return errors.New("无权删除此日报")
	}

	return s.reportRepo.Delete(reportID)
}

// GetReport 获取日报详情
func (s *ReportService) GetReport(reportID uint) (*model.DailyReport, error) {
	return s.reportRepo.GetByID(reportID)
}

// ListReports 获取日报列表
func (s *ReportService) ListReports(req *ListReportRequest) ([]model.DailyReport, int64, error) {
	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if req.UserID != nil {
		conditions["user_id"] = *req.UserID
	}
	if req.Status != "" {
		conditions["status"] = req.Status
	}
	if req.StartDate != nil {
		conditions["report_date_start"] = req.StartDate
	}
	if req.EndDate != nil {
		conditions["report_date_end"] = req.EndDate
	}
	if req.Department != "" {
		conditions["department"] = req.Department
	}

	return s.reportRepo.List(req.Page, req.PageSize, conditions)
}

// SaveDraft 保存草稿
func (s *ReportService) SaveDraft(userID uint, req *SaveDraftRequest) error {
	// 解析日期字符串
	reportDate, err := time.Parse("2006-01-02", req.ReportDate)
	if err != nil {
		return errors.New("日期格式错误，请使用YYYY-MM-DD格式")
	}

	draft := &model.ReportDraft{
		UserID:     userID,
		ReportDate: reportDate,
		TemplateID: req.TemplateID,
		Title:      req.Title,
		Content:    req.Content,
	}

	return s.draftRepo.CreateOrUpdate(draft)
}

// GetDraft 获取草稿
func (s *ReportService) GetDraft(userID uint, date time.Time) (*model.ReportDraft, error) {
	return s.draftRepo.GetByUserAndDate(userID, date)
}

// ApproveReportRequest 审批日报请求
type ApproveReportRequest struct {
	Status  string `json:"status" binding:"required"` // approved, rejected
	Comment string `json:"comment"`
}

// ApproveReport 审批日报
func (s *ReportService) ApproveReport(reportID, approverID uint, req *ApproveReportRequest) error {
	report, err := s.reportRepo.GetByID(reportID)
	if err != nil {
		return errors.New("日报不存在")
	}

	// 查找审批记录
	approvals, err := s.reportRepo.GetApprovalsByReportID(reportID)
	if err != nil {
		return fmt.Errorf("获取审批记录失败: %w", err)
	}

	var approval *model.ReportApproval
	for _, a := range approvals {
		if a.ApproverID == approverID && a.Status == "pending" {
			approval = &a
			break
		}
	}

	if approval == nil {
		return errors.New("未找到待审批记录")
	}

	now := time.Now()
	approval.Status = req.Status
	approval.Comment = req.Comment
	approval.ApprovedAt = &now

	if err := s.reportRepo.UpdateApproval(approval); err != nil {
		return fmt.Errorf("更新审批记录失败: %w", err)
	}

	// 更新日报状态
	if req.Status == "approved" {
		report.Status = "approved"
	} else if req.Status == "rejected" {
		report.Status = "rejected"
	}

	return s.reportRepo.Update(report)
}

// AddCommentRequest 添加评论请求
type AddCommentRequest struct {
	Content     string `json:"content" binding:"required"`
	ParentID    *uint  `json:"parent_id"`
	IsImportant bool   `json:"is_important"`
}

// AddComment 添加评论
func (s *ReportService) AddComment(reportID, userID uint, req *AddCommentRequest) (*model.ReportComment, error) {
	comment := &model.ReportComment{
		ReportID:    reportID,
		UserID:      userID,
		Content:     req.Content,
		ParentID:    req.ParentID,
		IsImportant: req.IsImportant,
	}

	if err := s.reportRepo.CreateComment(comment); err != nil {
		return nil, fmt.Errorf("创建评论失败: %w", err)
	}

	return comment, nil
}

// GetComments 获取评论列表
func (s *ReportService) GetComments(reportID uint) ([]model.ReportComment, error) {
	return s.reportRepo.GetCommentsByReportID(reportID)
}
