# 日常工作统计系统 - 后端服务

## 技术栈

- Go 1.21+
- Gin Web框架
- GORM ORM
- MySQL 8.0+
- JWT认证
- Zap日志

## 项目结构

```
backend/
├── cmd/server/          # 应用入口
├── internal/            # 内部代码
│   ├── api/            # API路由处理
│   ├── config/         # 配置管理
│   ├── model/          # 数据模型
│   ├── repository/     # 数据访问层
│   ├── service/        # 业务逻辑层
│   └── validator/      # 参数验证
├── pkg/                # 公共包
│   ├── database/       # 数据库连接
│   ├── jwt/            # JWT工具
│   ├── logger/         # 日志工具
│   ├── response/       # 响应处理
│   └── utils/          # 工具函数
├── config/             # 配置文件
└── migrations/         # 数据库迁移
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置数据库

编辑 `config/config.yaml` 文件，配置数据库连接信息：

```yaml
database:
  host: 10.0.51.35
  port: 3306
  user: root
  password: 123456
  dbname: oa
  charset: utf8mb4
```

### 3. 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

### 4. API文档

启动服务后，访问 `http://localhost:8080/swagger/index.html` 查看API文档（需要配置Swagger）

## 环境变量

可以通过环境变量覆盖配置：

- `SERVER_PORT`: 服务器端口
- `SERVER_MODE`: 运行模式 (debug/release/test)
- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名称
- `JWT_SECRET`: JWT密钥
- `LOG_LEVEL`: 日志级别

## API接口

### 认证接口

- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/auth/me` - 获取当前用户信息

## 开发规范

1. 遵循Go官方代码规范
2. 使用gofmt格式化代码
3. 所有公共函数需要注释
4. 错误处理要完善
5. 日志记录要详细

