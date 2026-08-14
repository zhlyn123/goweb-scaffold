package http

import (
	"context"
	"net/http"
	"time"

	"goweb-scaffold/internal/shared/requestid"
	"goweb-scaffold/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	database HealthChecker
	redis    HealthChecker
	timeout  time.Duration
}

func NewHandler(database HealthChecker, redis HealthChecker) *Handler {
	return &Handler{
		database: database,
		redis:    redis,
		timeout:  2 * time.Second,
	}
}

func (h *Handler) Live(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "UP",
	})
}

func (h *Handler) Ready(c *gin.Context) {
	status := "UP"
	dependencies := gin.H{}

	h.check(c.Request.Context(), dependencies, "database", h.database, &status)
	h.check(c.Request.Context(), dependencies, "redis", h.redis, &status)

	data := gin.H{
		"status":       status,
		"dependencies": dependencies,
	}
	if status != "UP" {
		c.JSON(http.StatusServiceUnavailable, response.Body{
			Code:      "service_unavailable",
			Message:   "service is not ready",
			Data:      data,
			RequestID: requestID(c),
		})
		return
	}

	response.Success(c, data)
}

func (h *Handler) check(parent context.Context, dependencies gin.H, name string, checker HealthChecker, status *string) {
	ctx, cancel := context.WithTimeout(parent, h.timeout)
	defer cancel()

	if checker == nil {
		*status = "DOWN"
		dependencies[name] = gin.H{
			"status": "DOWN",
			"error":  "dependency is not configured",
		}
		return
	}

	if err := checker.Ping(ctx); err != nil {
		*status = "DOWN"
		dependencies[name] = gin.H{
			"status": "DOWN",
			"error":  err.Error(),
		}
		return
	}

	dependencies[name] = gin.H{
		"status": "UP",
	}
}

func requestID(c *gin.Context) string {
	value, exists := c.Get(requestid.ConstextKey)
	if !exists {
		return ""
	}

	id, ok := value.(string)
	if !ok {
		return ""
	}

	return id
}
