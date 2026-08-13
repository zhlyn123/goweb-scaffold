package middleware

import (
	"context"
	"goweb-scaffold/internal/platform/security"
	"goweb-scaffold/internal/shared/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenParser interface {
	ParseAccessToken(ctx context.Context,tokenString string) (*security.AccessTokenCliaims, error)
}

func AuthRequired(tokenParser TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "缺少访问令牌")
			c.Abort()
			return
		}

		tokenString, ok := parseBearerToken(authHeader)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "访问令牌格式错误")
			c.Abort()
			return
		} 

		claims, err := tokenParser.ParseAccessToken(c.Request.Context(), tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "访问令牌无效")
			c.Abort()
			return
		}

		security.SetCurrentUserID(c, claims.UserID)

		c.Next()
	}
}

func parseBearerToken(authHeader string) (string, bool) {
	const prefix = "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		return "", false
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if tokenString == "" {
		return "", false
	}

	return tokenString, true
}
