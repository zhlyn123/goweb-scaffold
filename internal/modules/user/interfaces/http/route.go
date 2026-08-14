package http

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册 user 模块路由。
// 这里传入 authMiddleware，是为了让路由层决定哪些接口需要登录。
func RegisterRoutes(
	router *gin.Engine, handler *Handler,
	authMiddleware gin.HandlerFunc,
	permissionMiddleware gin.HandlerFunc,
) {
	users := router.Group("/api/v1/users")
	users.Use(authMiddleware)
	{
		users.GET("/me", permissionMiddleware, handler.Me)
	}
}
