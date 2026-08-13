package security

import "github.com/gin-gonic/gin"

const CurrentUserIDKey = "current_user_id"

// SetCurrentUserID 将当前登录用户 ID 写入 Gin Context。
// JWT 鉴权中间件解析 token 后调用它
func SetCurrentUserID(c *gin.Context, userID string) {
	c.Set(CurrentUserIDKey, userID)
}

// GetCurrentUserID 从 Gin Context 读取当前登录用户 ID。
// 需要登录态的 handler 可以通过它拿到当前用户。
func GetCurrentUserID(c *gin.Context) (string, bool) {
	value, exists := c.Get(CurrentUserIDKey)
	if !exists {
		return "", false
	}

	userID, ok := value.(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}
