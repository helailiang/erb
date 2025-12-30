package v1

import (
	"erb/internal/api/middleware"
	"erb/internal/config"
	"erb/internal/repository"
	"erb/internal/service"
	"erb/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReportHandler 日报处理器
type ReportHandler struct {
	reportService *service.ReportService
	cfg           *config.Config
}

// NewReportHandler 创建日报处理器
func NewReportHandler(db *gorm.DB, cfg *config.Config) *ReportHandler {
	reportRepo := repository.NewReportRepository(db)
	draftRepo := repository.NewDraftRepository(db)
	userRepo := repository.NewUserRepository(db)
	reportService := service.NewReportService(reportRepo, draftRepo, userRepo)
	return &ReportHandler{
		reportService: reportService,
		cfg:           cfg,
	}
}

// CreateReport 创建日报
// @Summary 创建日报
// @Description 创建新日报
// @Tags 日报管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param report body service.CreateReportRequest true "日报信息"
// @Success 200 {object} response.Response{data=model.DailyReport}
// @Failure 400 {object} response.Response
// @Router /api/v1/reports [post]
func (h *ReportHandler) CreateReport(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req service.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", req.ReportDate); err != nil {
		response.BadRequest(c, "日期格式错误，请使用YYYY-MM-DD格式")
		return
	}

	report, err := h.reportService.CreateReport(userID, &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, report)
}

// UpdateReport 更新日报
// @Summary 更新日报
// @Description 更新日报信息
// @Tags 日报管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "日报ID"
// @Param report body service.UpdateReportRequest true "日报信息"
// @Success 200 {object} response.Response{data=model.DailyReport}
// @Failure 400 {object} response.Response
// @Router /api/v1/reports/{id} [put]
func (h *ReportHandler) UpdateReport(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的日报ID")
		return
	}

	var req service.UpdateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	report, err := h.reportService.UpdateReport(uint(id), userID, &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, report)
}

// DeleteReport 删除日报
// @Summary 删除日报
// @Description 删除日报
// @Tags 日报管理
// @Security BearerAuth
// @Param id path int true "日报ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/reports/{id} [delete]
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的日报ID")
		return
	}

	if err := h.reportService.DeleteReport(uint(id), userID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetReport 获取日报详情
// @Summary 获取日报详情
// @Description 根据ID获取日报详细信息
// @Tags 日报管理
// @Security BearerAuth
// @Param id path int true "日报ID"
// @Success 200 {object} response.Response{data=model.DailyReport}
// @Failure 404 {object} response.Response
// @Router /api/v1/reports/{id} [get]
func (h *ReportHandler) GetReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的日报ID")
		return
	}

	report, err := h.reportService.GetReport(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, report)
}

// ListReports 获取日报列表
// @Summary 获取日报列表
// @Description 分页获取日报列表
// @Tags 日报管理
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param user_id query int false "用户ID"
// @Param status query string false "状态"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param department query string false "部门"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /api/v1/reports [get]
func (h *ReportHandler) ListReports(c *gin.Context) {
	var req service.ListReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 解析日期参数
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			req.StartDate = &t
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			req.EndDate = &t
		}
	}

	reports, total, err := h.reportService.ListReports(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      reports,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// SaveDraft 保存草稿
// @Summary 保存草稿
// @Description 保存日报草稿
// @Tags 日报管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param draft body service.SaveDraftRequest true "草稿信息"
// @Success 200 {object} response.Response
// @Router /api/v1/reports/draft [post]
func (h *ReportHandler) SaveDraft(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req service.SaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", req.ReportDate); err != nil {
		response.BadRequest(c, "日期格式错误，请使用YYYY-MM-DD格式")
		return
	}

	if err := h.reportService.SaveDraft(userID, &req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetDraft 获取草稿
// @Summary 获取草稿
// @Description 根据日期获取草稿
// @Tags 日报管理
// @Security BearerAuth
// @Param date query string true "日期" example(2024-01-01)
// @Success 200 {object} response.Response{data=model.ReportDraft}
// @Router /api/v1/reports/draft [get]
func (h *ReportHandler) GetDraft(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	dateStr := c.Query("date")
	if dateStr == "" {
		response.BadRequest(c, "日期参数不能为空")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "日期格式错误，请使用YYYY-MM-DD格式")
		return
	}

	draft, err := h.reportService.GetDraft(userID, date)
	if err != nil {
		response.NotFound(c, "草稿不存在")
		return
	}

	response.Success(c, draft)
}

// ApproveReport 审批日报
// @Summary 审批日报
// @Description 审批日报
// @Tags 日报管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "日报ID"
// @Param approval body service.ApproveReportRequest true "审批信息"
// @Success 200 {object} response.Response
// @Router /api/v1/reports/{id}/approve [post]
func (h *ReportHandler) ApproveReport(c *gin.Context) {
	approverID := c.MustGet("user_id").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的日报ID")
		return
	}

	var req service.ApproveReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.reportService.ApproveReport(uint(id), approverID, &req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}

// AddComment 添加评论
// @Summary 添加评论
// @Description 为日报添加评论
// @Tags 日报管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "日报ID"
// @Param comment body service.AddCommentRequest true "评论信息"
// @Success 200 {object} response.Response{data=model.ReportComment}
// @Router /api/v1/reports/{id}/comments [post]
func (h *ReportHandler) AddComment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的日报ID")
		return
	}

	var req service.AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	comment, err := h.reportService.AddComment(uint(id), userID, &req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, comment)
}

// GetComments 获取评论列表
// @Summary 获取评论列表
// @Description 获取日报的评论列表
// @Tags 日报管理
// @Security BearerAuth
// @Param id path int true "日报ID"
// @Success 200 {object} response.Response{data=[]model.ReportComment}
// @Router /api/v1/reports/{id}/comments [get]
func (h *ReportHandler) GetComments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的日报ID")
		return
	}

	comments, err := h.reportService.GetComments(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, comments)
}

// RegisterRoutes 注册路由
func (h *ReportHandler) RegisterRoutes(r *gin.RouterGroup, cfg *config.Config) {
	reports := r.Group("/reports")
	reports.Use(middleware.AuthMiddleware(cfg))
	{
		reports.POST("", h.CreateReport)
		reports.GET("", h.ListReports)
		reports.GET("/:id", h.GetReport)
		reports.PUT("/:id", h.UpdateReport)
		reports.DELETE("/:id", h.DeleteReport)
		reports.POST("/draft", h.SaveDraft) //草稿
		reports.GET("/draft", h.GetDraft)
		reports.POST("/:id/approve", h.ApproveReport)
		reports.POST("/:id/comments", h.AddComment)
		reports.GET("/:id/comments", h.GetComments)
	}
}
