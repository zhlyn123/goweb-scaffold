package middleware

import (
	"goweb-scaffold/internal/shared/requestid"
	"goweb-scaffold/internal/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			id, _ := c.Get(requestid.ConstextKey)

			log.Error(
				"HTTP 请求发生 panic",
				zap.Any("painc", err),
				zap.String("retquest_id", stringValue(id)),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
			)

			response.Error(
				c,
				http.StatusInternalServerError,
				"internal_error",
				"服务器内部错误",
			)

			c.Abort()
		}()

		c.Next()
	}
}
