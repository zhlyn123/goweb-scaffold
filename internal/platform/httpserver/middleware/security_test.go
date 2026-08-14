package middleware

import (
	"context"
	"goweb-scaffold/internal/platform/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestLoginRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST(
		"/login",
		LoginRateLimit(config.RateLimitConfig{
			Enabled:  true,
			Requests: 1,
			Window:   time.Minute,
		}),
		func(c *gin.Context) {
			c.Status(http.StatusOK)
		},
	)

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/login", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d", first.Code, http.StatusOK)
	}

	second := httptest.NewRecorder()
	router.ServeHTTP(second, httptest.NewRequest(http.MethodPost, "/login", nil))
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}
}

func TestBodySizeLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(BodySizeLimit(4))
	router.POST("/body", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/body", strings.NewReader("too-large"))

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestCORSRejectsWildcardOriginWithCredentialsInProd(t *testing.T) {
	_, err := CORS("prod", config.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	})
	if err == nil {
		t.Fatal("err is nil, want production CORS validation error")
	}
}

func TestRequirePermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET(
		"/protected",
		func(c *gin.Context) {
			c.Set("current_user_id", "user-1")
			c.Next()
		},
		RequirePermission(
			func(ctx context.Context, userID string, permissionCode string) error {
				if userID != "user-1" {
					t.Fatalf("userID = %s, want user-1", userID)
				}
				if permissionCode != "user:read" {
					t.Fatalf("permissionCode = %s, want user:read", permissionCode)
				}
				return nil
			},
			"user:read",
		),
		func(c *gin.Context) {
			c.Status(http.StatusOK)
		},
	)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/protected", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
}
