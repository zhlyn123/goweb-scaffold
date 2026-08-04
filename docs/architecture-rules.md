# 架构约束说明

这份文档用于约束 `goweb-scaffold` 后续开发时的代码放置和依赖方向。

## 1. 依赖方向

允许的依赖方向：

```text
cmd -> internal/app
internal/app -> internal/platform
internal/app -> internal/modules
interfaces -> application -> domain
infrastructure -> domain
internal/modules -> internal/shared
internal/platform -> third-party libraries
```

禁止的依赖方向：

```text
domain -> Gin / GORM / Redis / zap
application -> Gin
application -> GORM
handler -> database
module A infrastructure -> module B infrastructure
```

## 2. 目录职责

### cmd

只放应用入口，不写业务逻辑。

### internal/app

负责应用装配、依赖初始化、生命周期管理。

### internal/modules

负责业务模块。每个模块应包含自己的领域、应用、基础设施和接口层。

### internal/platform

负责配置、日志、数据库、缓存、HTTP 服务、可观测性等平台能力。

### internal/shared

只放稳定、通用、业务无关的共享能力。

## 3. 第一阶段规则

- 先保证应用入口和生命周期流程稳定。
- 暂不引入 Gin、GORM、Redis 等外部依赖。
- 后续技术组件通过 lifecycle hook 注册到应用中。
- 业务模块不能反向影响平台层结构。

