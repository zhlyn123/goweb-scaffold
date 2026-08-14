package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	health := router.Group("/health")
	{
		health.GET("/live", handler.Live)
		health.GET("/ready", handler.Ready)
	}
}
