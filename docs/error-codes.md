# 错误码说明

API 使用统一响应结构，业务错误码放在响应体 `code` 字段中，HTTP 状态码只表达协议层语义。

## 响应格式

成功：

```json
{
  "code": "ok",
  "message": "success",
  "data": {},
  "request_id": "request-id"
}
```

失败：

```json
{
  "code": "INVALID_ARGUMENT",
  "message": "request error",
  "request_id": "request-id"
}
```

## 通用错误码

| HTTP 状态码 | 错误码 | 含义 |
| --- | --- | --- |
| 400 | `INVALID_ARGUMENT` | 请求参数错误 |
| 401 | `UNAUTHORIZED` | 缺少、格式错误或无效的访问令牌 |
| 403 | `FORBIDDEN` | 已认证但没有访问权限 |
| 429 | `RATE_LIMITED` | 请求过于频繁 |
| 500 | `INTERNAL_ERROR` | 服务端内部错误 |
| 503 | `SERVICE_UNAVAILABLE` | 依赖不可用，服务未 ready |

## 业务错误码

| HTTP 状态码 | 错误码 | 场景 |
| --- | --- | --- |
| 401 | `INVALID_CREDENTIAL` | 登录邮箱或密码错误 |
| 403 | `USER_DISABLED` | 用户已禁用 |
| 409 | `USER_EMAIL_EXISTS` | 注册邮箱已存在 |

## 使用约定

- handler 负责把领域错误映射为 HTTP 状态码和业务错误码。
- application 和 domain 返回 Go error，不直接感知 HTTP 状态码。
- 新增错误码后，需要同步更新本文件和 `api/openapi/openapi.yaml`。
- 对外响应避免暴露数据库、Redis、JWT 解析等内部错误细节。
