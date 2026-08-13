package http

import (
	"errors"
	authapp "goweb-scaffold/internal/modules/auth/application"
	userdomain "goweb-scaffold/internal/modules/user/domain"
	"goweb-scaffold/internal/shared/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authUsecase *authapp.Usecase
}

func NewHandler(authUsecase *authapp.Usecase) *Handler {
	return &Handler{
		authUsecase: authUsecase,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ARGUMENT", "请求参数错误")
		return
	}
	result, err := h.authUsecase.Register(c.Request.Context(), authapp.RegisterCommand{
		Email:    req.Email,
		Password: req.Password,
		NickName: req.NickName,
	})
	if err != nil {
		writeAuthError(c, err)
		return
	}

	response.Success(c, RegisterResponse{
		UserID: result.UserID,
		Email:  result.Email,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ARGUMENT", "请求参数错误")
		return
	}
	result, err := h.authUsecase.Login(c.Request.Context(), authapp.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeAuthError(c, err)
		return
	}

	response.Success(c, LoginResponse{
		AccessToken: result.AccessToken,
		ExpiresAt:   result.ExpiresAt,
		TokenType:   result.TokenType,
	})
}

func writeAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, userdomain.ErrUserEmailExists):
		response.Error(c, http.StatusConflict, "USER_EMAIL_EXISTS", "邮箱已被注册")
	case errors.Is(err, userdomain.ErrInvalidPassword):
		response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIAL", "邮箱或密码错误")
	case errors.Is(err, userdomain.ErrUserDisabled):
		response.Error(c, http.StatusForbidden, "USER_DISABLED", "用户已被禁用")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "系统内部错误")
	}
}
