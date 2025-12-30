package v1

import (
	"erb/internal/api/middleware"
	"erb/internal/config"
	"erb/internal/repository"
	"erb/internal/service"
	"erb/pkg/response"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StatisticsHandler 统计处理器
type StatisticsHandler struct {
	statisticsService *service.StatisticsService
	cfg               *config.Config
}

// NewStatisticsHandler 创建统计处理器
func NewStatisticsHandler(db *gorm.DB, cfg *config.Config) *StatisticsHandler {
	reportRepo := repository.NewReportRepository(db)
	userRepo := repository.NewUserRepository(db)
	statisticsService := service.NewStatisticsService(reportRepo, userRepo)
	return &StatisticsHandler{
		statisticsService: statisticsService,
		cfg:               cfg,
	}
}

// GetPersonalStatistics 获取个人统计
// @Summary 获取个人统计
// @Description 获取个人日报统计数据
// @Tags 统计分析
// @Security BearerAuth
// @Param user_id query int false "用户ID（不传则查询当前用户）"
// @Param start_date query string false "开始日期" example(2024-01-01)
// @Param end_date query string false "结束日期" example(2024-01-31)
// @Success 200 {object} response.Response{data=service.PersonalStatisticsResponse}
// @Router /api/v1/statistics/personal [get]
func (h *StatisticsHandler) GetPersonalStatistics(c *gin.Context) {
	var req service.PersonalStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 如果没有指定用户ID，使用当前登录用户
	if req.UserID == 0 {
		req.UserID = c.MustGet("user_id").(uint)
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

	stats, err := h.statisticsService.GetPersonalStatistics(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetTeamStatistics 获取团队统计
// @Summary 获取团队统计
// @Description 获取团队日报统计数据
// @Tags 统计分析
// @Security BearerAuth
// @Param department query string false "部门"
// @Param start_date query string false "开始日期" example(2024-01-01)
// @Param end_date query string false "结束日期" example(2024-01-31)
// @Success 200 {object} response.Response{data=service.TeamStatisticsResponse}
// @Router /api/v1/statistics/team [get]
func (h *StatisticsHandler) GetTeamStatistics(c *gin.Context) {
	var req service.TeamStatisticsRequest
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

	stats, err := h.statisticsService.GetTeamStatistics(&req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, stats)
}

// RegisterRoutes 注册路由
func (h *StatisticsHandler) RegisterRoutes(r *gin.RouterGroup, cfg *config.Config) {
	statistics := r.Group("/statistics")
	statistics.Use(middleware.AuthMiddleware(cfg))
	{
		statistics.GET("/personal", h.GetPersonalStatistics)
		statistics.GET("/team", h.GetTeamStatistics)
	}
}

