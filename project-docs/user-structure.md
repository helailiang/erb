# 用户流程与项目结构

## 用户角色定义

### 1. 系统管理员
- **职责**：管理系统配置、用户账号、权限和角色
- **权限**：拥有系统所有权限
- **主要功能**：
  - 用户账号管理
  - 角色和权限管理
  - 系统配置管理
  - 数据备份和恢复
  - 审计日志查看

### 2. 部门管理员
- **职责**：管理本部门用户和查看本部门日报
- **权限**：管理本部门用户、查看和审批本部门日报
- **主要功能**：
  - 查看本部门成员日报
  - 审批本部门成员日报
  - 查看本部门统计报表
  - 管理本部门用户信息

### 3. 普通员工
- **职责**：填写和提交日报
- **权限**：填写和查看自己的日报
- **主要功能**：
  - 填写日报
  - 查看自己的历史日报
  - 查看审批反馈
  - 修改个人信息

### 4. 高级管理者
- **职责**：查看所有部门日报和统计
- **权限**：查看所有部门日报、统计报表
- **主要功能**：
  - 查看所有部门日报
  - 查看统计分析报表
  - 导出数据报表

## 用户旅程

### 1. 新用户首次使用流程

```
1. 系统管理员创建账号
   ↓
2. 用户收到账号通知（邮件/短信）
   ↓
3. 用户首次登录系统
   ↓
4. 系统提示修改密码
   ↓
5. 用户修改密码
   ↓
6. 进入系统首页
   ↓
7. 查看使用指南（可选）
   ↓
8. 开始使用系统
```

### 2. 日报填写流程

```
1. 用户登录系统
   ↓
2. 进入日报管理页面
   ↓
3. 选择填写今日日报
   ↓
4. 选择日报模板（如有多个模板）
   ↓
5. 填写日报内容
   - 今日工作
  - 明日计划
  - 问题与建议
   ↓
6. 上传附件（可选）
   ↓
7. 保存草稿或直接提交
   ↓
8. 提交日报
   ↓
9. 系统自动推送给审批人
   ↓
10. 等待审批反馈
```

### 3. 日报审批流程

```
1. 审批人收到日报提交通知
   ↓
2. 进入待审批日报列表
   ↓
3. 查看日报详情
   ↓
4. 阅读日报内容
   ↓
5. 添加批注或评论（可选）
   ↓
6. 选择审批操作
   - 通过
  - 驳回
  - 转交
   ↓
7. 填写审批意见（可选）
   ↓
8. 提交审批
   ↓
9. 系统通知日报提交人
   ↓
10. 日报提交人查看审批结果
```

### 4. 日报查询流程

```
1. 用户进入日报查询页面
   ↓
2. 选择查询范围
   - 个人日报
  - 团队日报（管理者）
  - 部门日报（部门管理员）
  - 全部日报（高级管理者）
   ↓
3. 设置筛选条件
   - 日期范围
  - 人员
  - 状态
  - 关键词
   ↓
4. 执行查询
   ↓
5. 查看查询结果列表
   ↓
6. 点击查看日报详情
   ↓
7. 导出日报（可选）
```

### 5. 用户管理流程（系统管理员）

```
1. 系统管理员进入用户管理页面
   ↓
2. 查看用户列表
   ↓
3. 执行用户管理操作
   - 创建新用户
  - 编辑用户信息
  - 禁用/启用用户
  - 删除用户
  - 分配角色
  - 分配权限
   ↓
4. 确认操作
   ↓
5. 系统执行操作并记录日志
   ↓
6. 通知相关用户（如需要）
```

## 数据流程

### 1. 日报数据流程

```
用户填写日报
   ↓
保存到草稿表（如保存草稿）
   ↓
提交日报
   ↓
保存到日报表
   ↓
创建审批记录
   ↓
发送通知给审批人
   ↓
审批人审批
   ↓
更新审批记录
   ↓
更新日报状态
   ↓
发送通知给提交人
   ↓
保存到历史记录
```

### 2. 用户认证流程

```
用户输入账号密码
   ↓
前端加密传输
   ↓
后端验证账号密码
   ↓
验证通过
   ↓
生成JWT Token
   ↓
返回Token给前端
   ↓
前端存储Token
   ↓
后续请求携带Token
   ↓
后端验证Token
   ↓
验证通过，允许访问
```

### 3. 权限验证流程

```
用户请求API
   ↓
验证Token有效性
   ↓
获取用户信息
   ↓
获取用户角色和权限
   ↓
检查请求资源所需权限
   ↓
验证用户是否有权限
   ↓
有权限：执行请求
   ↓
无权限：返回403错误
```

## 项目文件结构

### 完整项目结构

```
erb/
├── frontend/                    # 前端项目
│   ├── public/                 # 静态资源
│   ├── src/
│   │   ├── api/               # API接口
│   │   │   ├── user.ts        # 用户相关API
│   │   │   ├── report.ts      # 日报相关API
│   │   │   └── common.ts      # 通用API
│   │   ├── assets/            # 资源文件
│   │   ├── components/        # 公共组件
│   │   │   ├── common/        # 通用组件
│   │   │   ├── user/          # 用户相关组件
│   │   │   └── report/        # 日报相关组件
│   │   ├── composables/       # 组合式函数
│   │   │   ├── useAuth.ts     # 认证相关
│   │   │   ├── usePermission.ts # 权限相关
│   │   │   └── useReport.ts   # 日报相关
│   │   ├── directives/        # 自定义指令
│   │   ├── layouts/           # 布局组件
│   │   │   ├── DefaultLayout.vue
│   │   │   └── AdminLayout.vue
│   │   ├── router/            # 路由配置
│   │   │   ├── index.ts
│   │   │   ├── modules/
│   │   │   │   ├── user.ts    # 用户管理路由
│   │   │   │   └── report.ts  # 日报管理路由
│   │   ├── stores/            # 状态管理
│   │   │   ├── user.ts        # 用户状态
│   │   │   ├── report.ts      # 日报状态
│   │   │   └── app.ts         # 应用状态
│   │   ├── styles/            # 样式文件
│   │   ├── utils/             # 工具函数
│   │   │   ├── request.ts     # HTTP请求封装
│   │   │   ├── auth.ts        # 认证工具
│   │   │   ├── permission.ts  # 权限工具
│   │   │   └── format.ts      # 格式化工具
│   │   ├── views/             # 页面组件
│   │   │   ├── user/          # 用户管理页面
│   │   │   │   ├── UserList.vue
│   │   │   │   ├── UserCreate.vue
│   │   │   │   ├── UserEdit.vue
│   │   │   │   ├── RoleList.vue
│   │   │   │   └── PermissionList.vue
│   │   │   ├── report/        # 日报管理页面
│   │   │   │   ├── ReportList.vue
│   │   │   │   ├── ReportCreate.vue
│   │   │   │   ├── ReportDetail.vue
│   │   │   │   ├── ReportApproval.vue
│   │   │   │   └── ReportStatistics.vue
│   │   │   ├── dashboard/     # 仪表盘
│   │   │   │   └── Dashboard.vue
│   │   │   └── login/         # 登录页面
│   │   │       └── Login.vue
│   │   ├── App.vue
│   │   └── main.ts
│   ├── .eslintrc.js
│   ├── .prettierrc
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
│
├── backend/                    # 后端项目
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # 应用入口
│   ├── internal/
│   │   ├── api/               # API路由处理
│   │   │   ├── v1/            # API v1版本
│   │   │   │   ├── user.go    # 用户相关API
│   │   │   │   ├── report.go  # 日报相关API
│   │   │   │   └── common.go  # 通用API
│   │   │   └── middleware/    # 中间件
│   │   │       ├── auth.go    # 认证中间件
│   │   │       ├── permission.go # 权限中间件
│   │   │       └── logger.go  # 日志中间件
│   │   ├── config/            # 配置管理
│   │   │   └── config.go
│   │   ├── model/             # 数据模型
│   │   │   ├── user.go        # 用户模型
│   │   │   ├── role.go        # 角色模型
│   │   │   ├── permission.go  # 权限模型
│   │   │   ├── report.go       # 日报模型
│   │   │   └── common.go      # 通用模型
│   │   ├── repository/        # 数据访问层
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   ├── permission.go
│   │   │   └── report.go
│   │   ├── service/           # 业务逻辑层
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   ├── permission.go
│   │   │   └── report.go
│   │   ├── utils/             # 工具函数
│   │   │   ├── jwt.go         # JWT工具
│   │   │   ├── password.go    # 密码工具
│   │   │   └── validator.go   # 验证工具
│   │   └── validator/        # 参数验证
│   │       ├── user.go
│   │       └── report.go
│   ├── pkg/                   # 公共包
│   │   ├── database/          # 数据库连接
│   │   │   └── mysql.go
│   │   ├── logger/            # 日志工具
│   │   │   └── logger.go
│   │   ├── response/          # 响应处理
│   │   │   └── response.go
│   │   └── jwt/               # JWT认证
│   │       └── jwt.go
│   ├── migrations/            # 数据库迁移
│   │   ├── 001_create_users.sql
│   │   ├── 002_create_roles.sql
│   │   └── 003_create_reports.sql
│   ├── scripts/               # 脚本文件
│   ├── docs/                  # API文档
│   │   └── swagger.yaml
│   ├── go.mod
│   ├── go.sum
│   └── main.go
│
├── project-docs/              # 项目文档
│   ├── overview.md
│   ├── requirements.md
│   ├── tech-specs.md
│   ├── user-structure.md
│   └── timeline.md
│
├── .gitignore                 # Git忽略文件
├── README.md                  # 项目说明
└── docker-compose.yml         # Docker编排（可选）
```

## 核心功能模块划分

### 1. 用户管理模块
- 用户账号管理
- 角色管理
- 权限管理
- 安全与审计
- 自助服务

### 2. 日报管理模块
- 日报填写
- 日报提交与修改
- 日报审批与反馈
- 日报查询与统计

### 3. 系统管理模块
- 系统配置
- 消息通知
- 数据备份
- 日志管理

### 4. 通用功能模块
- 认证授权
- 文件上传
- 数据导出
- 统计分析

