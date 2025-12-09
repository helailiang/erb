package main

import (
	"erb/internal/config"
	"erb/internal/model"
	"erb/pkg/database"
	"erb/pkg/utils"
	"fmt"
	"log"
	"time"
)

// initData 初始化系统数据
func initData() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 检查是否已有数据
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("数据库已有数据，跳过初始化")
		return
	}

	// 创建默认角色
	roles := []model.Role{
		{
			Name:        "系统管理员",
			Description: "拥有系统所有权限",
			Status:      "enabled",
			IsSystem:    true,
		},
		{
			Name:        "部门管理员",
			Description: "管理本部门用户和查看本部门日报",
			Status:      "enabled",
			IsSystem:    true,
		},
		{
			Name:        "普通员工",
			Description: "填写和查看自己的日报",
			Status:      "enabled",
			IsSystem:    true,
		},
		{
			Name:        "高级管理者",
			Description: "查看所有部门日报和统计",
			Status:      "enabled",
			IsSystem:    true,
		},
	}

	for i := range roles {
		if err := db.Create(&roles[i]).Error; err != nil {
			log.Printf("创建角色失败: %v", err)
		}
	}

	// 创建默认管理员账号
	hashedPassword, err := utils.HashPassword("admin123456")
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	now := time.Now()
	admin := model.User{
		Username:          "admin",
		Password:          hashedPassword,
		RealName:          "系统管理员",
		Email:             "admin@example.com",
		Department:        "IT部门",
		Position:          "系统管理员",
		Status:            "normal",
		PasswordChangedAt: &now,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("创建管理员账号失败: %v", err)
	}

	// 分配系统管理员角色
	if err := db.Model(&admin).Association("Roles").Append(&roles[0]); err != nil {
		log.Printf("分配角色失败: %v", err)
	}

	fmt.Println("初始化数据完成！")
	fmt.Println("默认管理员账号：")
	fmt.Println("  用户名: admin")
	fmt.Println("  密码: admin123456")
	fmt.Println("  请登录后立即修改密码！")
}

func main() {
	initData()
}

