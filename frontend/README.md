# 日常工作统计系统 - 前端

## 技术栈

- Vue 3.x
- TypeScript
- Vite
- Element Plus
- Pinia
- Vue Router
- Axios

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API接口
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   ├── layouts/          # 布局组件
│   ├── router/           # 路由配置
│   ├── stores/           # 状态管理
│   ├── styles/           # 样式文件
│   ├── utils/            # 工具函数
│   ├── views/            # 页面组件
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

服务将在 `http://localhost:3000` 启动

### 3. 构建生产版本

```bash
npm run build
```

## 开发规范

1. 使用TypeScript严格模式
2. 遵循Vue 3 Composition API最佳实践
3. 组件命名使用PascalCase
4. 文件命名使用kebab-case
5. 使用ESLint进行代码检查

