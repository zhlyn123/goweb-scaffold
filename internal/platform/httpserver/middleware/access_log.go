package middleware

import (
	"goweb-scaffold/internal/shared/requestid"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AccessLog 记录每个 HTTP 请求的访问日志。
func AccessLog(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		id, _ := c.Get(requestid.ConstextKey)

		log.Info(
			"HTTP 请求完成",
			zap.String("request_id", stringValue(id)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// stringValue 将上下文里的值安全转换成字符串。
func stringValue(value any) string {
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}
