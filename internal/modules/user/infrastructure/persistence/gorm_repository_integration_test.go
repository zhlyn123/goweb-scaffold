//go:build integration

package persistence

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	userdomain "goweb-scaffold/internal/modules/user/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormRepositoryCreateAndFindUser(t *testing.T) {
	dsn := os.Getenv("GOWEB_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("GOWEB_TEST_DATABASE_DSN is not set")
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("ping postgres: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set goose dialect: %v", err)
	}
	migrationsDir := filepath.Join("..", "..", "..", "..", "..", "migrations")
	if err := goose.Up(sqlDB, migrationsDir); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}

	repo := NewGormRepository(db)
	ctx := context.Background()
	email := "repo-integration@example.com"
	_ = db.Exec("DELETE FROM user_roles WHERE user_id IN (SELECT id FROM users WHERE email = ?)", email).Error
	_ = db.Exec("DELETE FROM users WHERE email = ?", email).Error

	user := &userdomain.User{
		ID:           "11111111-1111-1111-1111-111111111111",
		Email:        email,
		PasswordHash: "hash",
		NickName:     "Repo Test",
		Status:       userdomain.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := repo.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	got, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		t.Fatalf("get user by email: %v", err)
	}
	if got.ID != user.ID {
		t.Fatalf("expected user ID %q, got %q", user.ID, got.ID)
	}
}
