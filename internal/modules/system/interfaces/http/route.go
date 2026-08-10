package http

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册 system 模块路由。
func RegisterRoutes(router *gin.Engine, handler *Handler) {
	health := router.Group("/health")
	{
		health.GET("/live", handler.Live)
	}
}
