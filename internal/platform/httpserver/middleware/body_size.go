package middleware

import (
	"goweb-scaffold/internal/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodySizeLimit 限制请求体大小，避免过大的请求占用过多内存和带宽。
func BodySizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if maxBytes <= 0 {
			c.Next()
			return
		}

		if c.Request.ContentLength > maxBytes {
			response.Error(c, http.StatusRequestEntityTooLarge, "REQUEST_BODY_TOO_LARGE", "请求体过大")
			c.Abort()
			return
		}

		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}

		c.Next()
	}
}
