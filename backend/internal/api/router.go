package api

import (
	"erb/internal/api/middleware"
	v1 "erb/internal/api/v1"
	"erb/internal/config"

	"github.com/gin-gonic/gin"
	third_middleware "github.com/zhufuyi/sponge/pkg/gin/middleware"
	"gorm.io/gorm"
)

// InitRouter 初始化路由
func InitRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	router := gin.New()

	// 全局中间件
	router.Use(third_middleware.RequestID())

	router.Use(middleware.LoggerMiddleware())
	router.Use(gin.Recovery())

	// CORS中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API路由组
	apiV1 := router.Group("/api/v1")
	{
		// 认证相关路由
		authHandler := v1.NewAuthHandler(db, cfg)
		authHandler.RegisterRoutes(apiV1, cfg)

		// 用户管理路由
		userHandler := v1.NewUserHandler(db, cfg)
		userHandler.RegisterRoutes(apiV1, cfg)

		// 角色管理路由
		roleHandler := v1.NewRoleHandler(db, cfg)
		roleHandler.RegisterRoutes(apiV1, cfg)

		// 权限管理路由
		permissionHandler := v1.NewPermissionHandler(db, cfg)
		permissionHandler.RegisterRoutes(apiV1, cfg)

		// 日报管理路由
		reportHandler := v1.NewReportHandler(db, cfg)
		reportHandler.RegisterRoutes(apiV1, cfg)

		// 统计分析路由
		statisticsHandler := v1.NewStatisticsHandler(db, cfg)
		statisticsHandler.RegisterRoutes(apiV1, cfg)
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 通用接口
	common := router.Group("/api/v1/common")
	{
		common.GET("/departments", v1.GetDepartments)
	}

	return router
}
