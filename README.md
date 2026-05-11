# Enterprise Agent

基于 **CloudWeGo Hertz + Eino** 构建的企业级智能 Agent 后端服务。采用 Go Workspace 多模块架构，支持用户认证、知识库管理、多轮对话及 AI Agent 调度，内置完整的可靠性中间件栈。

---

## 技术栈

| 层次 | 技术 |
|------|------|
| HTTP 框架 | [CloudWeGo Hertz](https://github.com/cloudwego/hertz) v0.10.4 |
| AI 框架 | [CloudWeGo Eino](https://github.com/cloudwego/eino) |
| 依赖注入 | [Google Wire](https://github.com/google/wire) |
| 结构化日志 | [Uber Zap](https://github.com/uber-go/zap) v1.27 |
| JWT 鉴权 | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) v5 |
| 熔断器 | [sony/gobreaker](https://github.com/sony/gobreaker) v0.5 |
| 密码散列 | bcrypt（golang.org/x/crypto） |
| 向量数据库 | Weaviate |
| 搜索引擎 | Elasticsearch 8 |
| 缓存 / 限流 | Redis 7 |
| API 文档 | Swagger（swaggo/swag） |

---

## 项目结构

```
enterprise-agent/
├── go.work                    # Go Workspace，统一管理所有子模块
├── Makefile                   # 常用命令入口
├── manifest/
│   ├── config/                # 配置文件（config.yaml / dev / prod）
│   └── docker/                # Dockerfile + docker-compose.yaml
│
├── module-pkg/                # 共享基础包（无业务依赖）
│   ├── appctx/                # 全局用户上下文存取（WithUser / GetUser）
│   ├── constants/             # 全局常量（Header / CtxKey / 用户状态等）
│   ├── jwt/                   # JWT 签发与解析（Sign / Parse）
│   ├── logger/                # Zap 日志初始化与全局实例
│   └── response/              # 统一响应格式与错误码
│
├── module-config/             # 配置加载（Viper，支持多环境）
│
├── module-infra/              # 基础设施客户端
│   └── middleware/
│       ├── es/                # Elasticsearch 客户端
│       ├── redis/             # Redis 客户端
│       └── weaviate/          # Weaviate 向量数据库客户端
│
├── module-core/               # 用户认证领域模块
│   ├── domain/
│   │   ├── model/user.go      # User / EmailCode 实体
│   │   └── repo/              # UserRepo / EmailCodeRepo / EmailSender 接口
│   ├── service/auth_service.go # 发验证码 / 注册 / 登录（bcrypt + JWT）
│   ├── handler/auth_handler.go # HTTP Handler（send-code / register / login / me）
│   └── wire/                  # Wire Provider Set
│
├── module-agent/              # AI Agent 业务模块
│   ├── domain/                # Agent 实体 & 仓储接口
│   ├── service/               # Agent 业务逻辑
│   ├── handler/               # HTTP Handler
│   └── wire/
│
├── module-knowledge/          # 知识库模块
│   ├── domain/
│   ├── service/
│   ├── handler/
│   └── wire/
│
├── module-conversation/       # 多轮对话模块
│   ├── domain/
│   ├── service/
│   ├── handler/
│   └── wire/
│
└── module-gateway/            # API 网关（唯一可执行入口）
    ├── main.go
    ├── router/router.go       # 路由注册 + 中间件挂载
    ├── middleware/
    │   ├── auth.go            # JWT 验证 + appctx 注入
    │   ├── trace.go           # 链路追踪（X-Request-ID）
    │   ├── rate_limit.go      # IP 固定窗口限流 + Redis 令牌桶限流
    │   ├── circuit_breaker.go # 熔断器（gobreaker）
    │   ├── timeout.go         # 请求超时控制
    │   ├── shedding.go        # 自适应降载（CPU + 在途请求）
    │   ├── cpu_linux.go       # Linux /proc/stat CPU 采样
    │   └── cpu_other.go       # 非 Linux 平台（返回 0，禁用 CPU 降载）
    ├── docs/                  # swag 生成的 Swagger 文档
    └── wire/                  # Wire Provider Set
```

---

## API 路由

### 公开接口（无需鉴权）

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/v1/auth/send-code` | 发送邮箱验证码（注册 / 重置密码） |
| POST | `/api/v1/auth/register` | 邮箱注册 |
| POST | `/api/v1/auth/login` | 邮箱密码登录，返回 JWT |
| GET  | `/health` | 健康检查 |
| GET  | `/swagger/*` | Swagger UI |

### 需要 JWT 鉴权（`Authorization: Bearer <token>`）

| 方法 | 路径 | 描述 |
|------|------|------|
| GET  | `/api/v1/auth/me` | 获取当前用户信息 |
| `*`  | `/api/v1/agent/*` | Agent 相关接口 |
| `*`  | `/api/v1/knowledge/*` | 知识库相关接口 |
| `*`  | `/api/v1/conversation/*` | 对话相关接口 |

---

## 中间件栈

请求经过全局中间件后进入对应路由组的中间件：

```
请求
 │
 ├─ [全局] TimeoutMiddleware        超时控制（默认 30s）
 ├─ [全局] TraceMiddleware          生成/透传 X-Request-ID，记录访问日志
 ├─ [全局] CORSMiddleware           跨域处理
 ├─ [全局] AdaptiveSheddingMiddleware  自适应降载（CPU > 75% 且在途 > 500）
 ├─ [全局] RateLimitMiddleware      IP 级固定窗口限流（200 req/min）
 │
 ├─ [/auth] CircuitBreakerMiddleware("auth")  认证路由熔断（5次失败 → 开路 30s）
 │
 └─ [/api/v1] AuthMiddleware        JWT 解析 + 用户信息注入 context
              CircuitBreakerMiddleware("api")   业务路由熔断（10次失败 → 开路 60s）
```

---

## 快速开始

### 前置依赖

- Go 1.21+
- Docker & Docker Compose（用于本地基础设施）

### 1. 启动基础设施

```bash
make docker-up
```

启动 Elasticsearch、Redis、Weaviate 三个服务。

### 2. 本地开发运行

```bash
# 安装热重载工具（首次）
go install github.com/air-verse/air@latest

# 热重载模式
make air

# 或者普通开发模式
make dev
```

默认监听 `:8080`，Swagger 文档：http://localhost:8080/swagger/index.html

### 3. 编译运行

```bash
make build      # 编译到 tmp/main
make run        # 编译并以生产模式运行
```

---

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `APP_ENV` | `dev` | 运行环境（`dev` / `prod`） |
| `ADDR` | `:8080` | 监听地址 |
| `LOG_LEVEL` | `info` | 日志级别（`debug` / `info` / `warn` / `error`） |
| `JWT_SECRET_KEY` | *(见配置文件)* | JWT 签名密钥，**生产环境必须通过此变量注入** |
| `RESOURCE_DIR` | `./manifest/config` | 配置文件目录 |

配置文件优先级：环境变量 > `config.{APP_ENV}.yaml` > `config.yaml`

---

## 常用命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译 gateway 到 `tmp/main` |
| `make run` | 编译并以生产模式运行 |
| `make dev` | 编译并以开发模式运行 |
| `make air` | 使用 Air 启动热重载开发服务 |
| `make swagger` | 生成 Swagger 文档（自动安装 swag） |
| `make wire` | 重新生成所有模块 Wire 代码（自动安装 wire） |
| `make tidy` | 对所有子模块执行 `go mod tidy` |
| `make test` | 运行所有模块单元测试 |
| `make test-cover` | 生成测试覆盖率报告（输出到 `tmp/*.coverage.out`） |
| `make docker` | 构建 Docker 镜像 |
| `make docker-up` | 启动所有基础设施容器（ES / Redis / Weaviate） |
| `make docker-down` | 停止所有容器 |
| `make docker-logs` | 实时查看 gateway 容器日志 |
| `make new-module NAME=xxx` | 从模板脚手架创建新业务模块 |
| `make clean` | 清除编译产物（`tmp/` 目录） |
| `make help` | 打印所有可用命令及说明 |

---

## 开发规范

### 新增业务模块

使用脚手架命令一键生成模块骨架：

```bash
make new-module NAME=foo
```

命令会将 `resource/template/module-xxx/` 复制为 `module-foo/`，并自动替换所有占位符。

生成后需手动完成以下步骤：

1. 在 `go.work` 的 `use` 块中添加 `./module-foo`
2. 在 `module-gateway/go.mod` 的 `require` 和 `replace` 块中引入新模块
3. 在 `module-gateway/router/router.go` 中注册路由
4. 在 `module-gateway/main.go` 中初始化 Handler
5. 执行 `make wire && make swagger` 更新依赖注入代码和 API 文档

模板位于 `resource/template/module-xxx/`，包含完整的 `domain / service / handler / wire` 骨架，可按需修改后作为团队标准模板。

### 依赖注入（Wire）

每个业务模块在 `wire/wire.go` 中导出 `SuperSet`，包含该模块的 service 层 Provider。Handler 在 `module-gateway/main.go` 中手动实例化，无需纳入 Wire 图。

### 用户上下文

在任意已认证路由的 Handler 或 Service 中获取当前用户：

```go
import "github.com/bing127/enterprise-agent/module-pkg/appctx"

// 读取（已认证路由中使用）
user, ok := appctx.GetUser(ctx)

// 强制读取（确保已认证，否则 panic）
user := appctx.MustGetUser(ctx)
```

---

## License

MIT
