# 技术规格说明

## 技术栈

### 前端技术栈
- **框架**：Vue 3.x
- **语言**：TypeScript
- **构建工具**：Vite
- **UI框架**：Element Plus
- **状态管理**：Pinia
- **路由**：Vue Router
- **HTTP客户端**：Axios
- **富文本编辑器**：TinyMCE
- **图表库**：ECharts
- **日期处理**：Day.js
- **工具库**：Lodash

### 后端技术栈
- **语言**：Go (Golang)
- **Web框架**：Gin 
- **ORM**：GORM
- **数据库**：MySQL 8.0+
- **缓存**：Redis
- **消息队列**：RabbitMQ
- **文件存储**：先本地存储，后续可以转为MinIO
- **日志**： Zap
- **配置管理**：Viper

### 开发工具
- **版本控制**：Git
- **代码规范**：ESLint (前端) / Golangci-lint (后端)
- **API文档**：Swagger / OpenAPI
- **测试框架**：Jest (前端) / Go Testing (后端)
- **CI/CD**：GitHub Actions / GitLab CI（可选）

## 开发方法

### 开发模式
- **前后端分离**：前端和后端独立开发和部署
- **RESTful API**：使用RESTful风格设计API接口
- **组件化开发**：前端采用组件化开发模式
- **模块化设计**：后端采用模块化设计，便于维护和扩展

### 代码规范

#### 前端代码规范
- 使用TypeScript严格模式
- 遵循Vue 3 Composition API最佳实践
- 组件命名使用PascalCase
- 文件命名使用kebab-case
- 使用ESLint进行代码检查
- 代码注释使用JSDoc格式

#### 后端代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 包命名使用小写字母
- 函数命名使用驼峰命名
- 使用golangci-lint进行代码检查
- 代码注释使用Go doc格式

### 项目结构

#### 前端项目结构
```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API接口定义
│   ├── assets/            # 资源文件
│   ├── components/        # 公共组件
│   ├── composables/      # 组合式函数
│   ├── directives/        # 自定义指令
│   ├── layouts/           # 布局组件
│   ├── router/            # 路由配置
│   ├── stores/            # 状态管理
│   ├── styles/            # 样式文件
│   ├── utils/             # 工具函数
│   ├── views/             # 页面组件
│   ├── App.vue            # 根组件
│   └── main.ts            # 入口文件
├── .eslintrc.js           # ESLint配置
├── .prettierrc            # Prettier配置
├── index.html             # HTML模板
├── package.json           # 依赖配置
├── tsconfig.json          # TypeScript配置
└── vite.config.ts         # Vite配置
```

#### 后端项目结构
```
backend/
├── cmd/
│   └── server/            # 应用入口
├── internal/
│   ├── api/               # API路由处理
│   ├── config/            # 配置管理
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── repository/        # 数据访问层
│   ├── service/           # 业务逻辑层
│   ├── utils/             # 工具函数
│   └── validator/         # 参数验证
├── pkg/                   # 公共包
│   ├── database/          # 数据库连接
│   ├── logger/            # 日志工具
│   ├── response/          # 响应处理
│   └── jwt/               # JWT认证
├── migrations/            # 数据库迁移
├── scripts/               # 脚本文件
├── docs/                  # API文档
├── go.mod                 # Go模块配置
├── go.sum                 # 依赖校验
└── main.go                # 主入口
```

## 数据库设计

### 数据库选型
- **主数据库**：MySQL 8.0+
- **缓存数据库**：Redis
- **数据库字符集**：utf8mb4
- **数据库排序规则**：utf8mb4_unicode_ci

### 核心数据表设计

#### 用户相关表
- `users` - 用户表
- `roles` - 角色表
- `permissions` - 权限表
- `user_roles` - 用户角色关联表
- `role_permissions` - 角色权限关联表
- `user_permissions` - 用户权限关联表（直接分配）

#### 日报相关表
- `daily_reports` - 日报表
- `report_templates` - 日报模板表
- `report_approvals` - 日报审批表
- `report_comments` - 日报评论表
- `report_attachments` - 日报附件表
- `report_drafts` - 日报草稿表

#### 系统相关表
- `audit_logs` - 审计日志表
- `notifications` - 通知表
- `system_configs` - 系统配置表
- `file_uploads` - 文件上传表

### 数据库设计原则
- 使用InnoDB存储引擎
- 合理设计索引，提高查询性能
- 使用外键约束保证数据完整性
- 使用软删除（deleted_at字段）
- 记录创建时间和更新时间
- 使用UUID或雪花算法生成主键（可选）

## API设计规范

### RESTful API设计
- 使用标准HTTP方法（GET、POST、PUT、DELETE、PATCH）
- 使用RESTful资源命名
- 统一的响应格式
- 统一的错误处理
- API版本控制（/api/v1/）

### 响应格式
```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### 错误码规范
- 2xx：成功
- 4xx：客户端错误
- 5xx：服务器错误

### API认证
- 使用JWT Token进行认证
- Token存储在HttpOnly Cookie或LocalStorage
- Token过期时间：7天
- 支持Token刷新

## 安全规范

### 前端安全
- XSS防护：对用户输入进行转义
- CSRF防护：使用CSRF Token
- 敏感信息不存储在LocalStorage
- API请求使用HTTPS

### 后端安全
- SQL注入防护：使用参数化查询
- XSS防护：对输出进行转义
- CSRF防护：验证CSRF Token
- 密码加密：使用bcrypt加密
- 接口限流：防止接口被恶意调用
- 权限验证：每个接口进行权限验证

## 部署规范

### 环境配置
- **开发环境**：本地开发
- **测试环境**：用于功能测试
- **预发布环境**：用于上线前验证
- **生产环境**：正式运行环境

### 部署方式
- **前端**：Nginx静态文件部署
- **后端**：Docker容器部署
- **数据库**：MySQL主从复制（可选）
- **缓存**：Redis集群（可选）

### 监控与日志
- 应用日志：使用ELK或类似工具收集
- 系统监控：使用Prometheus + Grafana
- 错误追踪：使用Sentry（可选）
- 性能监控：APM工具（可选）

## 测试规范

### 单元测试
- 前端：使用Jest进行单元测试，覆盖率 > 70%
- 后端：使用Go Testing进行单元测试，覆盖率 > 80%

### 集成测试
- API接口测试
- 数据库操作测试
- 第三方服务集成测试

### 端到端测试
- 使用Playwright或Cypress进行E2E测试
- 覆盖主要业务流程

## 文档规范

### 代码文档
- 所有公共函数和类需要注释
- 使用JSDoc（前端）和Go doc（后端）格式
- 复杂逻辑需要详细说明

### API文档
- 使用Swagger生成API文档
- 包含请求参数、响应格式、错误码说明
- 提供API调用示例

### 项目文档
- 项目概述文档
- 需求文档
- 技术规格文档
- 部署文档
- 用户手册

