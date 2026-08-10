package http

import (
	"github.com/gin-gonic/gin"

	"goweb-scaffold/internal/shared/response"
)

// Handler 负责 system 模块的 HTTP 请求处理。
type Handler struct{}

// NewHandler 创建 system 模块处理器。
func NewHandler() *Handler {
	return &Handler{}
}

// Live 返回进程存活状态。
func (h *Handler) Live(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "UP",
	})
}
