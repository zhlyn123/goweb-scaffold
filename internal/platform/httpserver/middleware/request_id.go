package middleware

import (
	"goweb-scaffold/internal/shared/requestid"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestid.HeaderName)
		if id == "" {
			id = requestid.New()
		}

		c.Set(requestid.ContextKey, id)
		c.Writer.Header().Set(requestid.HeaderName, id)

		c.Next()
	}
}
