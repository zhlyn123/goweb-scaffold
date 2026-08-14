# goweb-scaffold

`goweb-scaffold` 是一个面向生产级 Go Web 服务的脚手架项目，目标是沉淀高可用、低耦合、可复用、可演进的后端应用底座。

## 文档

- [架构设计文档](docs/architecture-design.md)
- [架构约束说明](docs/architecture-rules.md)
- [计划周期表](docs/roadmap.md)

## 当前能力

- Gin HTTP 服务
- Viper 配置加载
- zap 结构化日志
- PostgreSQL + GORM
- goose migration
- 统一响应结构
- request id、access log、recovery 中间件
- 用户注册、登录、JWT 鉴权
- RBAC 角色和权限模型
- 权限校验中间件
- 登录接口基础限流
- 请求体大小限制
- CORS 安全配置校验

## RBAC 说明

RBAC 基础表包括：

- `roles`：角色表，例如 `admin`、`user`
- `permissions`：权限表，例如 `user:read`、`user:manage`
- `user_roles`：用户和角色绑定关系
- `role_permissions`：角色和权限绑定关系

权限校验链路：

```text
JWT 鉴权中间件
-> 写入当前用户 ID
-> 权限中间件读取用户 ID
-> RBAC usecase 检查权限
-> 有权限继续执行 handler，无权限返回 403
```

当前 `/api/v1/users/me` 已接入 `user:read` 权限校验。

默认 RBAC 数据由 migration 初始化：

- `user` 角色默认拥有 `user:read`
- `admin` 角色拥有 `user:read` 和 `user:manage`
- 新注册用户会自动绑定 `user` 角色

## 安全配置

本地配置文件示例位于 `configs/config.local.yaml`。

```yaml
security:
  max_body_bytes: 1048576
  login_rate_limit:
    enabled: true
    requests: 5
    window: 1m
```

配置含义：

- `max_body_bytes`：请求体最大字节数，超过后返回 `413 REQUEST_BODY_TOO_LARGE`
- `login_rate_limit.enabled`：是否启用登录限流
- `login_rate_limit.requests`：单个 IP 在窗口期内允许的登录请求数
- `login_rate_limit.window`：限流窗口期

生产环境中，如果 `cors.allow_credentials` 为 `true`，则不允许 `cors.allow_origins` 使用 `*`。

## 本地验证

```bash
go test ./...
```

执行数据库迁移：

```bash
go run ./cmd/migrate up
```

启动服务：

```bash
go run ./cmd/api
```
