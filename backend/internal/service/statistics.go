package service

import (
	"erb/internal/model"
	"erb/internal/repository"
	"errors"
	"fmt"
	"time"
)

// StatisticsService 统计服务
type StatisticsService struct {
	reportRepo *repository.ReportRepository
	userRepo   *repository.UserRepository
}

// NewStatisticsService 创建统计服务
func NewStatisticsService(reportRepo *repository.ReportRepository, userRepo *repository.UserRepository) *StatisticsService {
	return &StatisticsService{
		reportRepo: reportRepo,
		userRepo:   userRepo,
	}
}

// PersonalStatisticsRequest 个人统计请求
type PersonalStatisticsRequest struct {
	UserID    uint       `form:"user_id"`
	StartDate *time.Time `form:"start_date"`
	EndDate   *time.Time `form:"end_date"`
}

// PersonalStatisticsResponse 个人统计响应
type PersonalStatisticsResponse struct {
	TotalReports      int     `json:"total_reports"`       // 总日报数
	SubmittedReports  int     `json:"submitted_reports"`  // 已提交数
	ApprovedReports   int     `json:"approved_reports"`    // 已审批数
	RejectedReports   int     `json:"rejected_reports"`   // 已驳回数
	LateReports       int     `json:"late_reports"`        // 补交数
	SubmitRate        float64 `json:"submit_rate"`         // 提交率
	AverageSubmitTime string  `json:"average_submit_time"` // 平均提交时间
	StatusDistribution map[string]int `json:"status_distribution"` // 状态分布
	DailySubmitCount  []DailySubmitCount `json:"daily_submit_count"` // 每日提交统计
}

// DailySubmitCount 每日提交统计
type DailySubmitCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// TeamStatisticsRequest 团队统计请求
type TeamStatisticsRequest struct {
	Department string     `form:"department"`
	StartDate  *time.Time `form:"start_date"`
	EndDate    *time.Time `form:"end_date"`
}

// TeamStatisticsResponse 团队统计响应
type TeamStatisticsResponse struct {
	TotalMembers      int                      `json:"total_members"`       // 总成员数
	ActiveMembers     int                      `json:"active_members"`     // 活跃成员数
	TotalReports      int                      `json:"total_reports"`      // 总日报数
	SubmitRate        float64                  `json:"submit_rate"`        // 提交率
	MemberStatistics  []MemberStatistics       `json:"member_statistics"`  // 成员统计
	StatusDistribution map[string]int          `json:"status_distribution"` // 状态分布
	DailySubmitTrend  []DailySubmitCount       `json:"daily_submit_trend"`  // 每日提交趋势
}

// MemberStatistics 成员统计
type MemberStatistics struct {
	UserID      uint    `json:"user_id"`
	UserName    string  `json:"user_name"`
	RealName    string  `json:"real_name"`
	TotalReports int    `json:"total_reports"`
	SubmitRate  float64 `json:"submit_rate"`
	LateCount   int     `json:"late_count"`
}

// GetPersonalStatistics 获取个人统计
func (s *StatisticsService) GetPersonalStatistics(req *PersonalStatisticsRequest) (*PersonalStatisticsResponse, error) {
	userID := req.UserID
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	// 设置默认日期范围（最近30天）
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	if req.StartDate != nil {
		startDate = *req.StartDate
	}
	if req.EndDate != nil {
		endDate = *req.EndDate
	}

	// 获取该用户的日报列表
	conditions := map[string]interface{}{
		"user_id": userID,
		"report_date_start": startDate,
		"report_date_end": endDate,
	}
	
	reports, _, err := s.reportRepo.List(1, 10000, conditions) // 获取所有符合条件的日报
	if err != nil {
		return nil, err
	}

	// 计算统计信息
	stats := &PersonalStatisticsResponse{
		TotalReports:      len(reports),
		StatusDistribution: make(map[string]int),
		DailySubmitCount:  make([]DailySubmitCount, 0),
	}

	// 按日期统计提交数量
	dailyCountMap := make(map[string]int)
	
	for _, report := range reports {
		// 统计状态分布
		stats.StatusDistribution[report.Status]++
		
		// 统计各状态数量
		switch report.Status {
		case "submitted":
			stats.SubmittedReports++
		case "approved":
			stats.ApprovedReports++
		case "rejected":
			stats.RejectedReports++
		}
		
		if report.IsLate {
			stats.LateReports++
		}
		
		// 统计每日提交数
		dateStr := report.ReportDate.Format("2006-01-02")
		dailyCountMap[dateStr]++
	}

	// 转换为数组
	for date, count := range dailyCountMap {
		stats.DailySubmitCount = append(stats.DailySubmitCount, DailySubmitCount{
			Date:  date,
			Count: count,
		})
	}

	// 计算提交率（提交数 / 工作日数）
	workDays := calculateWorkDays(startDate, endDate)
	if workDays > 0 {
		stats.SubmitRate = float64(stats.SubmittedReports) / float64(workDays) * 100
	}

	// 计算平均提交时间
	if stats.SubmittedReports > 0 {
		var totalMinutes int64
		count := 0
		for _, report := range reports {
			if report.SubmittedAt != nil {
				// 计算提交时间（小时）
				submitHour := report.SubmittedAt.Hour()
				totalMinutes += int64(submitHour)
				count++
			}
		}
		if count > 0 {
			avgHour := totalMinutes / int64(count)
			stats.AverageSubmitTime = fmt.Sprintf("%02d:00", avgHour)
		}
	}

	return stats, nil
}

// GetTeamStatistics 获取团队统计
func (s *StatisticsService) GetTeamStatistics(req *TeamStatisticsRequest) (*TeamStatisticsResponse, error) {
	// 设置默认日期范围（最近30天）
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	if req.StartDate != nil {
		startDate = *req.StartDate
	}
	if req.EndDate != nil {
		endDate = *req.EndDate
	}

	// 构建查询条件
	conditions := map[string]interface{}{
		"report_date_start": startDate,
		"report_date_end": endDate,
	}
	if req.Department != "" {
		conditions["department"] = req.Department
	}

	// 获取日报列表
	reports, _, err := s.reportRepo.List(1, 10000, conditions)
	if err != nil {
		return nil, err
	}

	// 获取用户列表（用于统计成员）
	userConditions := make(map[string]interface{})
	if req.Department != "" {
		userConditions["department = ?"] = req.Department
	}
	users, _, err := s.userRepo.List(1, 10000, userConditions)
	if err != nil {
		return nil, err
	}

	stats := &TeamStatisticsResponse{
		TotalMembers:      len(users),
		StatusDistribution: make(map[string]int),
		MemberStatistics:  make([]MemberStatistics, 0),
		DailySubmitTrend:  make([]DailySubmitCount, 0),
	}

	// 按用户统计
	userReportMap := make(map[uint][]model.DailyReport)
	for _, report := range reports {
		userReportMap[report.UserID] = append(userReportMap[report.UserID], report)
		stats.StatusDistribution[report.Status]++
	}

	// 统计每个成员的数据
	for _, user := range users {
		userReports := userReportMap[user.ID]
		if len(userReports) > 0 {
			stats.ActiveMembers++
		}

		submitCount := 0
		lateCount := 0
		for _, report := range userReports {
			if report.Status == "submitted" || report.Status == "approved" {
				submitCount++
			}
			if report.IsLate {
				lateCount++
			}
		}

		workDays := calculateWorkDays(startDate, endDate)
		submitRate := 0.0
		if workDays > 0 {
			submitRate = float64(submitCount) / float64(workDays) * 100
		}

		stats.MemberStatistics = append(stats.MemberStatistics, MemberStatistics{
			UserID:       user.ID,
			UserName:    user.Username,
			RealName:    user.RealName,
			TotalReports: len(userReports),
			SubmitRate:  submitRate,
			LateCount:   lateCount,
		})
	}

	// 统计每日提交趋势
	dailyCountMap := make(map[string]int)
	for _, report := range reports {
		dateStr := report.ReportDate.Format("2006-01-02")
		dailyCountMap[dateStr]++
	}
	for date, count := range dailyCountMap {
		stats.DailySubmitTrend = append(stats.DailySubmitTrend, DailySubmitCount{
			Date:  date,
			Count: count,
		})
	}

	// 计算团队提交率
	stats.TotalReports = len(reports)
	workDays := calculateWorkDays(startDate, endDate)
	if workDays > 0 && stats.TotalMembers > 0 {
		expectedReports := workDays * stats.TotalMembers
		if expectedReports > 0 {
			stats.SubmitRate = float64(stats.TotalReports) / float64(expectedReports) * 100
		}
	}

	return stats, nil
}

// calculateWorkDays 计算工作日数（简单实现，排除周末）
func calculateWorkDays(startDate, endDate time.Time) int {
	workDays := 0
	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		weekday := current.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			workDays++
		}
		current = current.AddDate(0, 0, 1)
	}
	return workDays
}

