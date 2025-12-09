package main

import (
	"erb/internal/api"
	"erb/internal/config"
	"erb/pkg/database"
	"erb/pkg/logger"
	"fmt"
	"log"
)

// @title 日常工作统计系统 API
// @version 1.0
// @description 日常工作统计系统后端API文档
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	logger.InitLogger(cfg.Log.Level, cfg.Log.Path)

	// 初始化数据库
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化路由
	router := api.InitRouter(db, cfg)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info(fmt.Sprintf("服务器启动在端口 %s", addr))
	if err := router.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

