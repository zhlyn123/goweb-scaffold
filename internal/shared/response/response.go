package response

import (
	"goweb-scaffold/internal/shared/requestid"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:      "ok",
		Message:   "success",
		Data:      data,
		RequestID: getRequestID(c),
	})
}

func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, Body{
		Code:      code,
		Message:   message,
		RequestID: getRequestID(c),
	})
}

func getRequestID(c *gin.Context) string {
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
