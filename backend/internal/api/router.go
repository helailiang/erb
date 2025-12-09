package api

import (
	"erb/internal/api/middleware"
	"erb/internal/api/v1"
	"erb/internal/config"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	router := gin.New()

	// 全局中间件
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
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return router
}

