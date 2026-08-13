package http

import (
	"errors"
	"goweb-scaffold/internal/modules/user/domain"
	"goweb-scaffold/internal/platform/security"
	"goweb-scaffold/internal/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userRepo domain.Repository
}

func NewHandler(userRepo domain.Repository) *Handler {
	return &Handler{
		userRepo: userRepo,
	}
}

func (h *Handler) Me(c *gin.Context) {
	userID, ok := security.GetCurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "未登录")
		return
	}
	user, err := h.userRepo.GetUserByID(c, userID)
	if errors.Is(err, domain.ErrUserNotFound) {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "用户不存在")
		return
	}

	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "系统内部错误")
		return
	}

	response.Success(c, MeResponse{
		UserID:   user.ID,
		Email:    user.Email,
		NickName: user.NickName,
		Status:   string(user.Status),
	})
}
