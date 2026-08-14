# goweb-scaffold

`goweb-scaffold` 是一个面向生产级 Go Web 服务的脚手架项目，目标是沉淀高可用、低耦合、可复用、可演进的后端应用底座。

## 文档

- [架构设计文档](docs/architecture-design.md)
- [架构约束说明](docs/architecture-rules.md)
- [计划周期表](docs/roadmap.md)
- [模块开发指南](docs/module-development-guide.md)
- [配置说明](docs/configuration.md)
- [错误码说明](docs/error-codes.md)
- [v0.1.0 checklist](docs/v0.1.0-checklist.md)
- [OpenAPI 文档](api/openapi/openapi.yaml)

## 当前能力

- Gin HTTP 服务
- Viper 配置加载，支持 `GOWEB_` 环境变量覆盖
- zap 结构化日志
- PostgreSQL + GORM
- Redis 客户端
- goose migration
- 统一响应结构
- request id、access log、recovery 中间件
- Prometheus metrics 和 OpenTelemetry tracing
- 用户注册、登录、JWT 鉴权
- RBAC 角色、权限、用户角色绑定和权限校验
- 登录限流、请求体大小限制、CORS 安全配置
- Docker、Docker Compose、Makefile、本地 OpenAPI/Swagger UI

## 快速启动

启动完整本地环境：

```bash
docker compose up --build
```

启动后可访问：

- API: http://localhost:8080
- 健康检查: http://localhost:8080/health/live
- Ready 检查: http://localhost:8080/health/ready
- Metrics: http://localhost:8080/metrics
- OpenAPI YAML: http://localhost:8080/openapi.yaml
- Swagger UI: http://localhost:8081
- Prometheus: http://localhost:9090

Compose 会启动 PostgreSQL、Redis、migration 任务、API 服务、Prometheus 和 Swagger UI。

## 本地开发

本地依赖 PostgreSQL 和 Redis 已启动后，执行迁移：

```bash
go run ./cmd/migrate -config configs/config.local.yaml up
```

启动服务：

```bash
go run ./cmd/api -config configs/config.local.yaml
```

也可以使用 Makefile：

```bash
make migrate-up
make run
```

常用命令：

```bash
make fmt
make tidy
make test
make docker-up
make docker-down
make openapi
```

## 测试

运行单元测试和默认 API 测试：

```bash
make test
```

运行 repository 集成测试需要显式提供测试数据库 DSN：

```bash
make test-integration TEST_DATABASE_DSN="host=localhost port=5432 user=postgres password=postgres dbname=goweb_scaffold sslmode=disable TimeZone=Asia/Shanghai"
```

集成测试使用 `integration` build tag，默认不会在 `go test ./...` 中执行。

## API 示例

注册：

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","nickname":"Demo"}'
```

登录：

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

获取当前用户：

```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <access_token>"
```

## RBAC 说明

RBAC 基础表包括：

- `roles`: 角色表，例如 `admin`、`user`
- `permissions`: 权限表，例如 `user:read`、`user:manage`
- `user_roles`: 用户和角色绑定关系
- `role_permissions`: 角色和权限绑定关系

默认 RBAC 数据由 migration 初始化：

- `user` 角色默认拥有 `user:read`
- `admin` 角色拥有 `user:read` 和 `user:manage`
- 新注册用户会自动绑定 `user` 角色

当前 `/api/v1/users/me` 已接入 JWT 鉴权和 `user:read` 权限校验。

## 配置

本地配置文件位于 `configs/config.local.yaml`，Docker Compose 使用 `configs/config.docker.yaml`。

关键安全配置：

```yaml
security:
  max_body_bytes: 1048576
  login_rate_limit:
    enabled: true
    requests: 5
    window: 1m
```

生产环境中，如果 `cors.allow_credentials` 为 `true`，则不允许 `cors.allow_origins` 使用 `*`。更多说明见 [配置说明](docs/configuration.md)。
