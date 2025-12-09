# 日常工作统计系统

## 项目简介

日常工作统计系统是一个面向企业内部使用的管理系统，旨在帮助企业管理员工日常工作记录、提升工作效率、加强团队协作。系统通过标准化的日报管理流程，使管理者能够及时了解团队成员的工作进展，同时为员工提供便捷的工作记录工具。

## 核心功能

### 用户管理模块
- **用户账号管理**：账号创建、信息维护、账号状态管理
- **权限管理**：权限定义、权限分配、权限审查
- **角色管理**：角色设计、角色分配、角色更新
- **安全与审计**：登录安全、操作审计、数据保护
- **自助服务**：个人信息管理、密码管理、通知管理

### 日报管理模块
- **日报填写**：模板支持、富文本编辑、定时提醒、草稿保存
- **提交与修改**：提交机制、补交功能
- **审批与反馈**：自动流转、批注/评论、审批操作、已读/未读状态
- **查询与统计**：个人视图、团队/部门视图、高级筛选、导出功能、统计分析

## 技术栈

### 前端
- Vue 3.x + TypeScript
- Vite
- Element Plus / Ant Design Vue
- Pinia
- Vue Router
- Axios

### 后端
- Go (Golang)
- Gin / Echo
- GORM
- MySQL 8.0+
- Redis

## 项目结构

```
erb/
├── frontend/          # 前端项目
├── backend/           # 后端项目
├── project-docs/      # 项目文档
│   ├── overview.md          # 项目概述
│   ├── requirements.md      # 需求与功能规格说明（PRD）
│   ├── tech-specs.md        # 技术规格说明
│   ├── user-structure.md    # 用户流程与项目结构
│   └── timeline.md          # 项目时间线与进度
└── README.md          # 项目说明
```

## 快速开始

### 环境要求
- Node.js >= 16.0.0
- Go >= 1.19
- MySQL >= 8.0
- Redis >= 6.0

### 安装步骤

#### 1. 数据库配置

确保MySQL数据库已启动，数据库连接信息：
- 主机：10.0.51.35
- 端口：3306
- 用户名：root
- 密码：123456
- 数据库：oa

数据库会在首次运行时自动创建表结构。

#### 2. 后端启动

```bash
cd backend
go mod download

# 编辑配置文件（如需要）
# 配置文件位置：backend/config/config.yaml

# 运行服务
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

#### 3. 初始化数据（可选）

如果需要创建默认管理员账号，运行：

```bash
cd backend
go run scripts/init_data.go
```

默认管理员账号：
- 用户名：admin
- 密码：admin123456

#### 4. 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 `http://localhost:3000` 启动

#### 5. 访问系统

打开浏览器访问：`http://localhost:3000`

使用默认管理员账号登录即可开始使用。

详细文档请参考 [project-docs](./project-docs/) 目录下的文档。

## 项目文档

- [项目概述](./project-docs/overview.md) - 项目背景、愿景、目标
- [需求文档](./project-docs/requirements.md) - 详细的功能需求规格说明（PRD）
- [技术规格](./project-docs/tech-specs.md) - 技术栈、开发方法、数据库设计
- [用户流程](./project-docs/user-structure.md) - 用户角色、流程、项目结构
- [时间线](./project-docs/timeline.md) - 项目里程碑和进度跟踪

## 开发计划

项目预计开发周期：11周

- **第一阶段**（2周）：项目准备与设计
- **第二阶段**（4周）：基础功能开发
- **第三阶段**（3周）：高级功能开发
- **第四阶段**（2周）：测试与优化
- **第五阶段**（持续）：上线与维护

详细进度请参考 [项目时间线](./project-docs/timeline.md)。

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证。

## 联系方式

如有问题或建议，请联系项目维护者。

