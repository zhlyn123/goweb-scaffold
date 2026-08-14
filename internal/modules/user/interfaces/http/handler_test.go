package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	userdomain "goweb-scaffold/internal/modules/user/domain"
	"goweb-scaffold/internal/platform/security"
	"goweb-scaffold/internal/shared/response"

	"github.com/gin-gonic/gin"
)

func TestMeReturnsCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newUserRepo(map[string]*userdomain.User{
		"user-1": {
			ID:        "user-1",
			Email:     "user@example.com",
			NickName:  "Demo",
			Status:    userdomain.UserStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})
	router := gin.New()
	router.GET("/api/v1/users/me", func(c *gin.Context) {
		security.SetCurrentUserID(c, "user-1")
		c.Next()
	}, NewHandler(repo).Me)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var body response.Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %#v", body.Data)
	}
	if data["email"] != "user@example.com" {
		t.Fatalf("expected email user@example.com, got %#v", data["email"])
	}
}

func TestMeReturnsUnauthorizedWithoutCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/users/me", NewHandler(newUserRepo(nil)).Me)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

type userRepo struct {
	users map[string]*userdomain.User
}

func newUserRepo(users map[string]*userdomain.User) *userRepo {
	if users == nil {
		users = map[string]*userdomain.User{}
	}
	return &userRepo{users: users}
}

func (r *userRepo) CreateUser(_ context.Context, _ *userdomain.User) error {
	return nil
}

func (r *userRepo) GetUserByID(_ context.Context, id string) (*userdomain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, userdomain.ErrUserNotFound
	}
	return user, nil
}

func (r *userRepo) GetUserByEmail(_ context.Context, email string) (*userdomain.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, userdomain.ErrUserNotFound
}
