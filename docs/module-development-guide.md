# 模块开发指南

本文档用于指导在 `goweb-scaffold` 中新增业务模块。默认架构是模块化单体，模块内部遵循 `interfaces -> application -> domain` 的依赖方向，基础设施实现只依赖领域层接口。

## 推荐目录

```text
internal/modules/<module>
  domain
    entity.go
    repository.go
    errors.go
  application
    usecase.go
    usecase_test.go
  infrastructure
    persistence
      gorm_repository.go
  interfaces
    http
      handler.go
      request.go
      response.go
      route.go
      handler_test.go
```

## 开发步骤

1. 在 `domain` 中定义实体、仓储接口和领域错误。
2. 在 `application` 中编排用例，接收 `context.Context`，只依赖领域接口和必要的跨模块应用接口。
3. 在 `infrastructure` 中实现仓储接口，隔离 GORM、Redis 等第三方细节。
4. 在 `interfaces/http` 中完成请求绑定、响应转换和路由注册。
5. 在 `internal/app/bootstrap` 中装配仓储、用例、handler 和路由。
6. 在 `api/openapi/openapi.yaml` 中补充接口契约。
7. 按风险补测试：用例单测、handler API 测试、必要的 repository 集成测试。

## 约束

- `domain` 不依赖 Gin、GORM、Redis、zap。
- `application` 不依赖 Gin 和 GORM。
- handler 不直接操作数据库。
- Repository 不承载复杂业务规则。
- 跨模块协作优先通过应用服务接口或领域事件，不直接访问对方基础设施实现。

## 当前参考模块

- `auth`：注册、登录、JWT 签发和默认角色绑定。
- `user`：当前用户查询和用户仓储。
- `rbac`：角色、权限、用户角色绑定和权限校验。
- `system`：健康检查、ready 检查、OpenAPI 和 metrics 暴露。
