package middleware

import (
	"time"

	"goweb-scaffold/internal/shared/requestid"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func AccessLog(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		id, _ := c.Get(requestid.ConstextKey)
		traceID := ""
		spanContext := trace.SpanContextFromContext(c.Request.Context())
		if spanContext.HasTraceID() {
			traceID = spanContext.TraceID().String()
		}

		log.Info(
			"HTTP request completed",
			zap.String("request_id", stringValue(id)),
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

func stringValue(value any) string {
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}
