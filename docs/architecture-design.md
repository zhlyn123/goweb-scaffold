# Go Web Scaffold 架构设计文档

## 1. 项目定位

`goweb-scaffold` 定位为一个面向中大型业务系统的 Go Web 应用脚手架。

它不是一个简单的 CRUD 模板，而是一个可长期复用的应用底座，重点解决以下问题：

- 高可用：具备健康检查、优雅启停、超时控制、限流、可观测性、配置治理等能力。
- 低耦合：业务逻辑不直接依赖 Gin、GORM、Redis 等具体基础设施。
- 可复用：沉淀通用工程能力、模块规范、鉴权体系、错误模型、部署模板。
- 可演进：初期采用模块化单体架构，后续可平滑拆分为微服务。

## 2. 设计目标

### 2.1 核心目标

- 提供清晰稳定的项目结构。
- 内置生产级基础设施能力。
- 建立统一的模块开发规范。
- 支持 REST API 服务快速开发。
- 支持数据库、缓存、鉴权、日志、监控、链路追踪。
- 支持 Docker Compose 本地开发环境。
- 预留 Kubernetes、消息队列、任务队列、微服务拆分能力。

### 2.2 非目标

第一阶段不追求以下能力：

- 完整微服务治理体系。
- Service Mesh。
- 多租户复杂隔离。
- 复杂工作流引擎。
- 大而全的后台管理系统。
- 过度抽象的 DDD 框架。

脚手架应优先保证可理解、可运行、可维护，再逐步增强。

## 3. 架构风格

推荐采用：

```text
Clean Architecture
+
Hexagonal Architecture
+
Modular Monolith
```

核心思想：

- 业务规则放在内层，框架和数据库放在外层。
- 业务模块之间边界清晰。
- 应用初期部署为一个进程，但代码结构为未来拆服务预留边界。
- 依赖方向必须从外层指向内层，内层不能依赖外层实现。

## 4. 推荐技术栈

| 能力 | 推荐方案 | 说明 |
|---|---|---|
| Web 框架 | Gin | 成熟、生态丰富、学习资料多 |
| 数据库 | PostgreSQL | 长期扩展性强，事务和索引能力优秀 |
| ORM | GORM | 上手快，配合 Repository 接口隔离 |
| 缓存 | Redis | 缓存、限流、分布式锁、会话扩展 |
| 配置 | Viper | 支持配置文件和环境变量 |
| 日志 | zap | 高性能结构化日志 |
| 鉴权 | JWT + RBAC | 先满足常见后台和 API 场景 |
| 参数校验 | validator | Go 生态常用校验库 |
| 数据迁移 | goose | 简洁，适合工程化迁移 |
| API 文档 | OpenAPI / Swagger | 支持接口协作和调试 |
| 指标监控 | Prometheus | 暴露服务指标 |
| 链路追踪 | OpenTelemetry | 预留分布式追踪能力 |
| 容器化 | Docker / Docker Compose | 本地和测试环境标准化 |
| 部署预留 | Kubernetes | 后续生产环境部署模板 |
| 测试 | testing + testify | 单元测试和集成测试 |

## 5. 总体分层

```text
cmd
  api
    main.go

internal
  app
    bootstrap
    container
    lifecycle

  modules
    user
      domain
      application
      infrastructure
      interfaces
    auth
      domain
      application
      infrastructure
      interfaces
    system
      interfaces

  platform
    config
    logger
    database
    cache
    httpserver
    telemetry
    metrics
    security
    eventbus
    transaction
    migration
    validator

  shared
    response
    errors
    pagination
    idgen
    clock

api
  openapi

configs
deployments
  docker
  k8s
migrations
scripts
tests
```

## 6. 分层职责

### 6.1 cmd

应用入口层。

职责：

- 创建应用上下文。
- 加载配置。
- 初始化依赖容器。
- 启动 HTTP 服务。
- 监听系统信号，执行优雅关闭。

### 6.2 internal/app

应用装配层。

职责：

- 统一初始化配置、日志、数据库、缓存、路由、监控等组件。
- 管理依赖注入。
- 管理应用生命周期。
- 避免初始化逻辑散落在业务模块中。

### 6.3 internal/modules

业务模块层。每个模块应尽量自包含。

推荐模块内部结构：

```text
modules/user
  domain
    entity.go
    repository.go
    service.go
    errors.go
  application
    command
    query
    usecase.go
  infrastructure
    persistence
    cache
    event
  interfaces
    http
      handler.go
      route.go
      request.go
      response.go
```

职责说明：

- `domain`：领域实体、领域服务、仓储接口、领域错误。
- `application`：用例编排、事务边界、命令查询模型。
- `infrastructure`：数据库、缓存、消息等具体实现。
- `interfaces`：HTTP handler、路由注册、请求响应转换。

### 6.4 internal/platform

平台能力层。

职责：

- 封装与业务无关的基础设施能力。
- 对外提供稳定接口或初始化方法。
- 集中管理第三方库接入方式。

典型能力：

- 配置加载。
- 日志初始化。
- 数据库连接池。
- Redis 客户端。
- HTTP Server。
- JWT 安全组件。
- OpenTelemetry。
- Prometheus。
- 事务管理器。
- 事件总线。

### 6.5 internal/shared

共享工具层。

只放真正跨模块、业务无关、稳定复用的内容。

适合放入：

- 统一响应结构。
- 统一错误码。
- 分页结构。
- ID 生成器。
- 时间抽象。

不适合放入：

- 用户业务逻辑。
- 权限业务规则。
- 具体模块的请求结构。
- 被多个模块暂时引用但语义并不通用的代码。

## 7. 依赖规则

必须遵守以下依赖方向：

```text
interfaces -> application -> domain
infrastructure -> domain
application -> shared
platform -> 第三方库
```

禁止：

- `domain` 依赖 Gin、GORM、Redis、zap。
- `application` 直接使用 Gin Context。
- `handler` 直接操作数据库。
- 模块之间直接访问对方数据库表。
- 基础设施实现反向污染业务接口。

推荐：

- 业务层依赖接口。
- 基础设施层实现接口。
- HTTP 层只做协议转换。
- 跨模块协作通过应用服务或领域事件完成。

## 8. 请求处理链路

典型请求链路：

```text
HTTP Request
  -> Middleware
  -> Handler
  -> Request DTO Validation
  -> Application Usecase
  -> Domain Service / Entity
  -> Repository Interface
  -> Infrastructure Implementation
  -> Database / Cache
  -> Response DTO
  -> Unified Response
```

Handler 只负责：

- 读取参数。
- 调用校验。
- 获取登录上下文。
- 调用 Usecase。
- 返回响应。

Usecase 负责：

- 编排业务流程。
- 控制事务边界。
- 调用领域服务。
- 调用仓储接口。
- 发布领域事件。

Repository 负责：

- 持久化实体。
- 隔离 GORM 或 SQL 细节。
- 不承载复杂业务规则。

## 9. 高可用设计

### 9.1 优雅启停

应用需要支持：

- 启动阶段依赖检查。
- HTTP Server 优雅关闭。
- 数据库连接关闭。
- Redis 连接关闭。
- 后台任务停止。
- 上下文超时控制。

### 9.2 健康检查

提供两个接口：

```text
GET /health/live
GET /health/ready
```

`live` 表示进程是否存活。

`ready` 表示服务是否可接收流量，应检查：

- 数据库连接。
- Redis 连接。
- 必要外部依赖。

### 9.3 超时控制

必须建立统一超时策略：

- HTTP Server read timeout。
- HTTP Server write timeout。
- Handler context timeout。
- 数据库查询 timeout。
- Redis 操作 timeout。
- 外部 HTTP 调用 timeout。

### 9.4 限流

第一阶段支持：

- IP 级基础限流。
- 登录接口单独限流。

后续扩展：

- 用户级限流。
- 接口级限流。
- Redis 分布式限流。

### 9.5 可观测性

默认内置：

- 结构化日志。
- Request ID。
- Trace ID。
- Prometheus metrics。
- OpenTelemetry tracing。

关键指标：

- 请求总数。
- 请求耗时。
- HTTP 状态码分布。
- Panic 次数。
- 数据库连接池状态。
- Redis 操作耗时。
- 业务错误码分布。

## 10. 低耦合设计

### 10.1 Repository 接口隔离

业务层定义接口：

```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
}
```

基础设施层实现接口：

```go
type GormUserRepository struct {
    db *gorm.DB
}
```

这样业务层不感知 GORM。

### 10.2 Cache 接口隔离

业务层只依赖缓存接口，不直接依赖 Redis 客户端。

### 10.3 EventBus 接口隔离

初期使用进程内事件总线。

后续可替换为：

- Kafka。
- NATS。
- RabbitMQ。
- Redis Stream。

### 10.4 Transaction Manager

事务不应散落在每个 Repository 中。

推荐由 Application Usecase 控制事务：

```text
transaction.Run(ctx, func(ctx context.Context) error {
  // 执行业务操作
})
```

Repository 从 context 或 transaction manager 获取事务连接。

## 11. 模块设计

### 11.1 system 模块

职责：

- 健康检查。
- 版本信息。
- 服务信息。
- 指标暴露。

### 11.2 auth 模块

职责：

- 登录。
- 刷新 Token。
- 退出登录。
- 密码校验。
- JWT 签发。
- 当前登录用户上下文。

### 11.3 user 模块

职责：

- 用户创建。
- 用户查询。
- 用户状态管理。
- 用户基础资料。

### 11.4 rbac 模块

职责：

- 角色管理。
- 权限管理。
- 用户角色绑定。
- 接口权限校验。

第一阶段可以先实现基础模型和核心校验，不做复杂后台界面。

## 12. 统一错误模型

错误应分为：

- 系统错误。
- 参数错误。
- 鉴权错误。
- 权限错误。
- 业务错误。
- 外部依赖错误。

推荐统一结构：

```json
{
  "code": "USER_NOT_FOUND",
  "message": "user not found",
  "request_id": "req_xxx",
  "details": {}
}
```

HTTP 状态码和业务错误码分开管理。

示例：

| HTTP 状态码 | 业务错误码 | 含义 |
|---|---|---|
| 400 | INVALID_ARGUMENT | 参数错误 |
| 401 | UNAUTHORIZED | 未登录 |
| 403 | FORBIDDEN | 无权限 |
| 404 | NOT_FOUND | 资源不存在 |
| 409 | CONFLICT | 资源冲突 |
| 429 | RATE_LIMITED | 请求过于频繁 |
| 500 | INTERNAL_ERROR | 系统内部错误 |

## 13. 统一响应模型

成功响应：

```json
{
  "code": "OK",
  "message": "success",
  "data": {},
  "request_id": "req_xxx"
}
```

分页响应：

```json
{
  "code": "OK",
  "message": "success",
  "data": {
    "items": [],
    "page": 1,
    "page_size": 20,
    "total": 100
  },
  "request_id": "req_xxx"
}
```

## 14. 配置设计

配置来源优先级：

```text
默认值 < 配置文件 < 环境变量 < 启动参数
```

推荐配置文件：

```text
configs
  config.local.yaml
  config.dev.yaml
  config.test.yaml
  config.prod.yaml
```

配置项分类：

- app。
- http。
- log。
- database。
- redis。
- jwt。
- telemetry。
- metrics。
- security。

敏感信息不应提交到 Git。

## 15. 数据库设计规范

建议：

- 使用 PostgreSQL。
- 主键默认使用 UUID 或雪花 ID。
- 所有业务表包含 `created_at`、`updated_at`。
- 需要软删除时使用 `deleted_at`。
- 重要业务表包含 `version` 字段用于乐观锁。
- 所有 schema 变化通过 migration 管理。
- 禁止应用启动时自动修改生产表结构。

## 16. 安全设计

第一阶段内置：

- 密码 bcrypt 加密。
- JWT access token。
- Refresh token 预留。
- CORS 配置。
- 请求体大小限制。
- 登录限流。
- 统一鉴权中间件。

后续增强：

- OAuth2 / OIDC。
- API Key。
- 操作审计。
- 敏感字段脱敏。
- CSRF 防护。

## 17. 部署设计

第一阶段提供：

- Dockerfile。
- docker-compose.yml。
- PostgreSQL 容器。
- Redis 容器。
- 应用容器。
- 本地开发 Makefile。

后续提供：

- Kubernetes Deployment。
- Service。
- ConfigMap。
- Secret 示例。
- Ingress。
- HPA。
- Prometheus scrape 配置。

## 18. 测试策略

测试分层：

- Domain 单元测试。
- Application 用例测试。
- Repository 集成测试。
- Handler API 测试。
- Middleware 测试。

最低要求：

- 核心业务用例必须有单元测试。
- Repository 需要基于测试数据库做集成测试。
- 登录、鉴权、健康检查需要 API 测试。

## 19. 代码规范

基本规则：

- 所有函数必须接收 `context.Context`，除非明显不需要。
- 不在业务层使用全局变量。
- 不在 Handler 中写业务逻辑。
- 不直接 panic，启动失败除外。
- 日志必须结构化。
- 错误必须包装上下文。
- 接口定义靠近使用方。
- 不为未来幻想过度抽象。

## 20. 第一阶段交付范围

第一阶段建议完成：

- 项目基础目录结构。
- 配置加载。
- zap 日志。
- Gin HTTP Server。
- 统一响应。
- 统一错误。
- recovery、request id、access log 中间件。
- PostgreSQL + GORM。
- Redis 客户端。
- goose migration。
- 健康检查。
- Prometheus metrics。
- JWT 鉴权。
- 用户注册、登录、查询当前用户。
- 基础 RBAC 数据模型。
- Docker Compose 本地环境。
- README 和开发文档。

## 21. 演进路线

后续可以按以下方向演进：

- 增加后台管理模块。
- 增加审计日志。
- 增加异步任务队列。
- 增加领域事件。
- 增加 Redis 分布式限流。
- 增加 OpenTelemetry Collector 示例。
- 增加 Kubernetes 部署模板。
- 增加 CI/CD。
- 增加代码生成器。
- 增加模块脚手架命令。
