# goweb-scaffold

一个面向生产级 Go Web 服务的脚手架规划项目，目标是构建高可用、低耦合、可复用、可演进的后端应用底座。

## 文档

- [架构设计文档](docs/architecture-design.md)
- [架构约束说明](docs/architecture-rules.md)
- [计划周期表](docs/roadmap.md)

## 当前阶段

当前正在执行第 1 周：架构骨架。

已完成：

- Go module 初始化。
- `cmd/api` 应用入口。
- `internal/app` 应用装配层。
- 生命周期管理器。
- `internal/modules`、`internal/platform`、`internal/shared` 基础目录。

## 本地验证

```bash
go test ./...
```
