package middleware

import (
	"context"
	"errors"
	rbacdomain "goweb-scaffold/internal/modules/rbac/domain"
	"goweb-scaffold/internal/platform/security"
	"goweb-scaffold/internal/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PermissionCheckFunc 定义权限检查函数。
// 中间件只依赖这个函数类型，避免 platform 层直接依赖 RBAC usecase。
type PermissionCheckFunc func(ctx context.Context, userID string, permissionCode string) error

// RequirePermission 创建权限校验中间件。
// 使用方式：在已登录路由后面追加 RequirePermission("权限编码")。
func RequirePermission(checkPermission PermissionCheckFunc, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := security.GetCurrentUserID(c)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未登录")
			c.Abort()
			return
		}

		err := checkPermission(c.Request.Context(), userID, permissionCode)

		if errors.Is(err, rbacdomain.ErrPermissionDenied) {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "权限不足")
			c.Abort()
			return
		}

		if err != nil {
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "系统内部错误")
			c.Abort()
			return
		}

		c.Next()
	}
}
