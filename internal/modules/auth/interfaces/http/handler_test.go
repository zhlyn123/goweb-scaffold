package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authapp "goweb-scaffold/internal/modules/auth/application"
	userdomain "goweb-scaffold/internal/modules/user/domain"
	"goweb-scaffold/internal/shared/response"

	"github.com/gin-gonic/gin"
)

type fakeUserRepository struct {
	usersByEmail map[string]*userdomain.User
	usersByID    map[string]*userdomain.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		usersByEmail: map[string]*userdomain.User{},
		usersByID:    map[string]*userdomain.User{},
	}
}

func (r *fakeUserRepository) CreateUser(_ context.Context, user *userdomain.User) error {
	r.usersByEmail[user.Email] = user
	r.usersByID[user.ID] = user
	return nil
}

func (r *fakeUserRepository) GetUserByID(_ context.Context, id string) (*userdomain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, userdomain.ErrUserNotFound
	}
	return user, nil
}

func (r *fakeUserRepository) GetUserByEmail(_ context.Context, email string) (*userdomain.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return nil, userdomain.ErrUserNotFound
	}
	return user, nil
}

type fakeHasher struct{}

func (fakeHasher) HashPassword(password string) (string, error) {
	return "hash:" + password, nil
}

func (fakeHasher) CheckPassword(password string, passwordHash string) error {
	if passwordHash != "hash:"+password {
		return errors.New("invalid password")
	}
	return nil
}

type fakeTokenIssuer struct{}

func (fakeTokenIssuer) IssueAccessToken(context.Context, string) (string, time.Time, error) {
	return "access-token", time.Now().Add(time.Hour), nil
}

type fakeIDGenerator struct{}

func (fakeIDGenerator) NewID() string {
	return "00000000-0000-0000-0000-000000000001"
}

type passthroughTx struct{}

func (passthroughTx) Run(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestAuthRegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := newFakeUserRepository()
	usecase := authapp.NewUsecase(repo, fakeHasher{}, fakeTokenIssuer{}, fakeIDGenerator{}, passthroughTx{}, nil, "")
	router := gin.New()
	RegisterRoutes(router, NewHandler(usecase), func(c *gin.Context) { c.Next() })

	registerBody := []byte(`{"email":"user@example.com","password":"password123","nickname":"Demo"}`)
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	registerRec := httptest.NewRecorder()
	router.ServeHTTP(registerRec, registerReq)

	if registerRec.Code != http.StatusOK {
		t.Fatalf("expected register status %d, got %d: %s", http.StatusOK, registerRec.Code, registerRec.Body.String())
	}

	loginBody := []byte(`{"email":"user@example.com","password":"password123"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	router.ServeHTTP(loginRec, loginReq)

	if loginRec.Code != http.StatusOK {
		t.Fatalf("expected login status %d, got %d: %s", http.StatusOK, loginRec.Code, loginRec.Body.String())
	}

	var body response.Body
	if err := json.Unmarshal(loginRec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %#v", body.Data)
	}
	if data["access_token"] != "access-token" {
		t.Fatalf("expected access token, got %#v", data["access_token"])
	}
}

func TestAuthRegisterRejectsInvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	usecase := authapp.NewUsecase(newFakeUserRepository(), fakeHasher{}, fakeTokenIssuer{}, fakeIDGenerator{}, passthroughTx{}, nil, "")
	router := gin.New()
	RegisterRoutes(router, NewHandler(usecase), func(c *gin.Context) { c.Next() })

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(`{"email":"bad","password":"password123"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
