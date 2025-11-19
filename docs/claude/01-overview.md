# 项目概览

## 项目简介

Chat Nio 是新一代 AIGC 一站式 B/C 端解决方案，提供强大的 API 分发系统和丰富的用户界面。它同时作为 C 端（客户端）聊天界面和 B 端（业务端）API 代理/分发平台，支持多个 AI 模型提供商（OpenAI、Anthropic、Gemini、Midjourney 及 15+ 其他提供商）。

**技术栈：**
- **前端**: React 18 + Redux Toolkit + Radix UI + Tailwind CSS + Vite
- **后端**: Go 1.20 + Gin + MySQL + Redis
- **应用**: PWA + WebSocket 实时通信

---

## 目录结构

### 后端目录（Go）

| 目录 | 主要职责 | 关键文件 |
|------|----------|----------|
| **adapter/** | AI 提供商适配器，实现多种 AI 服务的统一接口 | `adapter.go`, `openai/`, `claude/`, `midjourney/` 等 |
| **channel/** | 渠道管理系统，负载均衡、优先级路由、模型重定向 | `manager.go`, `ticker.go`, `charge.go` |
| **manager/** | 业务逻辑层，处理聊天、图像、视频生成请求 | `chat.go`, `chat_completions.go`, `images.go`, `videos.go` |
| **auth/** | 认证授权、用户管理、配额管理 | `auth.go`, `controller.go`, `quota.go` |
| **admin/** | 管理员面板 API，用户管理、统计分析 | `controller.go`, `user.go`, `instance.go` |
| **connection/** | 数据库连接管理（MySQL/SQLite） | `database.go`, `worker.go`, `cache.go` |
| **middleware/** | HTTP 中间件（CORS、认证、限流） | `auth.go`, `cors.go`, `builtins.go` |
| **utils/** | 工具函数，WebSocket 处理、日志、配置加载 | `config.go`, `websocket.go`, `net.go` |
| **globals/** | 全局常量、类型定义、数据结构 | `channel.go`, `types.go` |
| **migration/** | 数据库迁移脚本 | 各版本迁移文件 |
| **addition/** | 附加功能（文章生成、卡片、Web 搜索） | `article/`, `card/`, `generation/`, `web/` |
| **cli/** | 命令行接口工具 | `cli.go` |

### 前端目录（React）

| 目录 | 主要职责 | 关键文件 |
|------|----------|----------|
| **app/src/components/** | 可复用 UI 组件（Radix UI + Tailwind） | 各种 UI 组件 |
| **app/src/routes/** | 页面组件 | 各页面路由组件 |
| **app/src/admin/** | 管理后台页面 | 管理面板相关页面 |
| **app/src/store/** | Redux 状态管理（auth、chat、settings） | Redux slices |
| **app/src/api/** | API 客户端函数 | API 调用封装 |
| **app/src/types/** | TypeScript 类型定义 | 类型声明文件 |
| **app/src/utils/** | 前端工具函数 | 辅助函数 |
| **app/src/dialogs/** | 模态对话框 | 各种对话框组件 |
| **app/src/assets/** | 静态资源（CSS、图片、i18n） | `globals.less`, 国际化文件 |

---

## 构建与运行方式

### 后端构建与运行

#### 方式一：直接编译运行

```bash
# 构建二进制文件
go build -o chatnio

# 运行应用（默认端口：8094）
./chatnio

# 使用自定义配置运行
MYSQL_HOST=localhost REDIS_HOST=localhost ./chatnio
```

**配置说明：**
- 配置文件位于 `config/config.yaml`
- 环境变量会覆盖配置文件（例如：`MYSQL_HOST` 覆盖 `mysql.host`）
- 需要 MySQL 和 Redis 运行中

#### 方式二：Docker 部署

```bash
# 启动所有服务（应用 + MySQL + Redis）
docker-compose up -d

# 使用稳定版本
docker-compose -f docker-compose.stable.yaml up -d

# 启用 Watchtower 自动更新
docker-compose -f docker-compose.watch.yaml up -d

# 停止所有服务
docker-compose down

# 更新镜像
docker-compose pull
```

**Docker 卷挂载：**
- `~/db` - MySQL 数据
- `~/redis` - Redis 数据
- `~/config` - 配置文件
- `~/logs` - 应用日志
- `~/storage` - 文件上传和生成内容

### 前端构建与运行

```bash
cd app

# 安装依赖（使用 pnpm，不要用 npm）
pnpm install

# 启动开发服务器（Vite）
pnpm dev

# 生产构建（包含 TypeScript 编译）
pnpm build

# 快速构建（跳过 TypeScript 检查）
pnpm fast-build

# ESLint 检查
pnpm lint

# 代码格式化
pnpm prettier
```

---

## 外部依赖

### 必需依赖

#### 数据库
- **MySQL** (推荐生产环境)
  - 版本：5.7+
  - 用途：用户数据、对话历史、配额管理
  - 配置：`config.yaml` 中的 `mysql` 部分
  - 备注：若未配置 MySQL，会自动降级使用 SQLite（`~/db/chatnio.db`）

- **SQLite** (开发环境备用)
  - 自动创建于 `./db/chatnio.db`
  - 仅在 `mysql.host` 未配置时启用

#### 缓存
- **Redis**
  - 版本：5.0+
  - 用途：会话缓存、速率限制、实时通信
  - 配置：`config.yaml` 中的 `redis` 部分

### Go 框架依赖（主要）

| 依赖包 | 版本 | 用途 |
|--------|------|------|
| **gin-gonic/gin** | v1.9.1 | HTTP Web 框架 |
| **spf13/viper** | v1.16.0 | 配置管理 |
| **go-redis/redis** | v8.11.5 | Redis 客户端 |
| **go-sql-driver/mysql** | v1.7.1 | MySQL 驱动 |
| **mattn/go-sqlite3** | v1.14.22 | SQLite 驱动 |
| **gorilla/websocket** | v1.5.0 | WebSocket 支持 |
| **dgrijalva/jwt-go** | v3.2.0 | JWT 认证 |
| **sirupsen/logrus** | v1.9.3 | 日志记录 |
| **goccy/go-json** | v0.10.2 | 高性能 JSON 解析 |

### 前端框架依赖（主要）

| 依赖包 | 版本 | 用途 |
|--------|------|------|
| **react** | 18.2.0 | 前端框架 |
| **@reduxjs/toolkit** | 1.9.5 | 状态管理 |
| **react-router-dom** | 6.17.0 | 路由管理 |
| **@radix-ui/** | 多个包 | 无障碍 UI 组件 |
| **tailwindcss** | 3.3.3 | CSS 框架 |
| **axios** | 1.5.0 | HTTP 客户端 |
| **vite** | 4.4.5 | 构建工具 |
| **typescript** | 5.0.2 | 类型系统 |
| **react-markdown** | 8.0.7 | Markdown 渲染 |
| **mermaid** | 10.9.0 | 图表渲染 |

### 第三方 AI API

项目支持多个 AI 提供商，需要相应的 API Key：
- **OpenAI API** (GPT-3.5, GPT-4)
- **Anthropic Claude API**
- **Google Gemini / PaLM2**
- **Azure OpenAI**
- **Midjourney** (图像生成)
- **百度文心一言、阿里通义千问、腾讯混元** 等国内模型

---

## 新手阅读顺序

### 第一阶段：了解项目结构

1. **CLAUDE.md** - 项目整体介绍和开发指南（当前文件的上层文档）
2. **docs/claude/01-overview.md** - 本文档，项目概览
3. **README.md** / **README_zh-CN.md** - 项目说明和功能介绍

### 第二阶段：理解启动流程

4. **docs/claude/02-entrypoint.md** - 程序入口与启动流程
5. **main.go** - 主入口文件
   - 阅读启动顺序：`ReadConf()` → `InitInstance()` → `InitManager()` → 注册路由 → 启动服务器
6. **utils/config.go** - 配置加载逻辑
7. **connection/database.go** - 数据库初始化和表创建

### 第三阶段：核心业务逻辑

8. **docs/claude/03-callchains.md** - 核心调用链分析
9. **adapter/adapter.go** - 适配器模式实现
   - 查看如何注册和创建不同 AI 提供商的适配器
10. **channel/manager.go** - 渠道管理和负载均衡
    - 理解优先级路由、权重分配、自动重试
11. **manager/chat.go** - 聊天请求处理流程
    - 从请求接收到响应返回的完整流程

### 第四阶段：模块间协作

12. **docs/claude/04-modules.md** - 模块依赖与数据流
13. **auth/auth.go** - 用户认证和授权
14. **middleware/auth.go** - 认证中间件
15. **manager/connection.go** - WebSocket 连接管理

### 第五阶段：前端架构

16. **app/src/App.tsx** - 前端入口
17. **app/src/router.tsx** - 路由配置
18. **app/src/store/** - Redux 状态管理
19. **app/src/api/** - API 调用封装

### 第六阶段：系统架构

20. **docs/claude/05-architecture.md** - 完整系统架构文档

### 实践建议

- **运行项目**: 先用 Docker 快速启动，体验完整功能
- **代码调试**: 在 `utils/config.go` 中启用 `debug: true`，查看详细日志
- **添加适配器**: 参考 `adapter/openai/` 或 `adapter/claude/` 实现新的 AI 提供商
- **API 测试**: 使用 Postman 或 curl 测试 `/api/v1/chat/completions` 接口
- **数据库查看**: 使用 MySQL Workbench 或 DBeaver 查看数据库结构

---

## 默认凭据

部署后的默认管理员账户：
- **用户名**: `root`
- **密码**: `chatnio123456`
- **建议**: 首次登录后立即在管理面板中修改密码

---

## 项目特色

1. **多提供商支持**: 统一接口适配 15+ AI 提供商
2. **智能负载均衡**: 优先级路由、权重分配、自动故障转移
3. **灵活计费系统**: 支持按 Token、按次数、按时间等多种计费模式
4. **实时通信**: WebSocket + SSE 双协议支持
5. **企业级特性**: 用户分组、配额管理、订阅系统、邀请码
6. **渐进式 Web 应用**: PWA 支持，可离线使用
7. **国际化**: 多语言支持（中文、英文、日文）
8. **双模式部署**: 一体化部署或前后端分离部署
