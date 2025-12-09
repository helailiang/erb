# 项目启动指南

## 快速启动

### 前置条件

1. **Go环境**：Go 1.21 或更高版本
2. **Node.js环境**：Node.js 16.0 或更高版本
3. **MySQL数据库**：MySQL 8.0 或更高版本
4. **数据库已创建**：确保数据库 `oa` 已存在

### 数据库配置

数据库连接信息已在 `backend/config/config.yaml` 中配置：

```yaml
database:
  host: 10.0.51.35
  port: 3306
  user: root
  password: 123456
  dbname: oa
  charset: utf8mb4
```

如需修改，请编辑该配置文件。

### 启动步骤

#### 第一步：启动后端服务

```bash
# 进入后端目录
cd backend

# 下载依赖
go mod download

# 启动服务
go run cmd/server/main.go
```

看到以下输出表示启动成功：
```
服务器启动在端口 :8080
```

#### 第二步：初始化数据（首次运行）

首次运行需要初始化默认数据：

```bash
# 在backend目录下运行
go run scripts/init_data.go
```

这将创建：
- 4个默认角色（系统管理员、部门管理员、普通员工、高级管理者）
- 1个默认管理员账号（用户名：admin，密码：admin123456）

#### 第三步：启动前端服务

打开新的终端窗口：

```bash
# 进入前端目录
cd frontend

# 安装依赖（首次运行）
npm install

# 启动开发服务器
npm run dev
```

看到以下输出表示启动成功：
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: use --host to expose
```

#### 第四步：访问系统

1. 打开浏览器访问：`http://localhost:3000`
2. 使用默认管理员账号登录：
   - 用户名：`admin`
   - 密码：`admin123456`
3. 登录成功后进入系统首页

## 常见问题

### 1. 数据库连接失败

**错误信息**：`数据库连接失败`

**解决方法**：
- 检查MySQL服务是否启动
- 检查数据库连接信息是否正确
- 确认数据库 `oa` 已创建
- 检查防火墙设置

### 2. 端口被占用

**错误信息**：`bind: address already in use`

**解决方法**：
- 修改 `backend/config/config.yaml` 中的端口号
- 或关闭占用端口的程序

### 3. 前端依赖安装失败

**错误信息**：`npm ERR!`

**解决方法**：
- 检查Node.js版本是否符合要求
- 尝试清除缓存：`npm cache clean --force`
- 使用国内镜像：`npm config set registry https://registry.npmmirror.com`

### 4. 跨域问题

如果遇到跨域问题，检查：
- 后端CORS配置是否正确
- 前端代理配置是否正确（`frontend/vite.config.ts`）

## 开发模式

### 后端开发

后端使用热重载，修改代码后需要手动重启服务。

### 前端开发

前端使用Vite，支持热模块替换（HMR），修改代码后自动刷新。

## 生产部署

### 后端部署

```bash
# 编译
go build -o erb-server cmd/server/main.go

# 运行
./erb-server
```

### 前端部署

```bash
# 构建
npm run build

# 构建产物在 dist 目录
# 使用Nginx等Web服务器部署
```

## 下一步

- 查看 [项目文档](./project-docs/) 了解详细功能
- 查看 [API文档](./backend/README.md) 了解后端接口
- 查看 [前端文档](./frontend/README.md) 了解前端开发

