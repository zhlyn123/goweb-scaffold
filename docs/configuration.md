# 配置说明

配置加载入口位于 `internal/platform/config`。加载优先级为：

```text
默认值 < YAML 配置文件 < GOWEB_ 环境变量
```

本地默认文件是 `configs/config.local.yaml`，Docker Compose 使用 `configs/config.docker.yaml`。

## 环境变量覆盖

环境变量前缀是 `GOWEB_`，字段路径中的 `.` 使用 `_` 替代。

示例：

```bash
GOWEB_HTTP_ADDR=":8088"
GOWEB_LOG_LEVEL="info"
GOWEB_DATABASE_HOST="localhost"
GOWEB_JWT_SECRET="change-me"
```

## 配置分组

| 分组 | 说明 |
| --- | --- |
| `app` | 应用名称和运行环境 |
| `http` | 监听地址、读写超时、优雅关闭超时 |
| `log` | 日志级别和输出格式 |
| `cors` | 跨域来源、方法、请求头、凭证和缓存时间 |
| `security` | 请求体大小限制、登录限流 |
| `database` | PostgreSQL 连接、连接池和超时 |
| `redis` | Redis 地址、数据库编号和超时 |
| `metrics` | Prometheus 指标开关和路径 |
| `telemetry` | OpenTelemetry 开关、服务名和 exporter |
| `jwt` | JWT secret、issuer 和 access token TTL |

## 生产注意事项

- 必须通过环境变量或 Secret 注入 `GOWEB_JWT_SECRET` 和数据库密码。
- `cors.allow_credentials=true` 时不要把 `cors.allow_origins` 配成 `*`。
- `log.format` 建议使用 `json`，便于日志平台采集。
- `database.max_open_conns` 需要结合数据库实例规格和服务副本数设置。
- `metrics.enabled` 可在不需要暴露指标的环境关闭。
