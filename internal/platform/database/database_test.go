package database

import (
	"context"
	"testing"
	"time"

	"goweb-scaffold/internal/platform/config"
)

func TestBuildDSN(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secret",
		Name:     "goweb_scaffold",
		SSLMode:  "disable",
		TimeZone: "Asia/Shanghai",
	}

	got := buildDSN(cfg)
	want := "host=localhost port=5432 user=postgres password=secret dbname=goweb_scaffold sslmode=disable TimeZone=Asia/Shanghai"

	if got != want {
		t.Fatalf("buildDSN() = %q, want %q", got, want)
	}
}

func TestBuildDSNPreferRawDSN(t *testing.T) {
	cfg := config.DatabaseConfig{
		DSN: "postgres://user:pass@localhost:5432/app?sslmode=disable",
	}

	got := buildDSN(cfg)
	want := "postgres://user:pass@localhost:5432/app?sslmode=disable"

	if got != want {
		t.Fatalf("buildDSN() = %q, want %q", got, want)
	}
}

func TestCloseNilDB(t *testing.T) {
	var db *DB

	if err := db.Close(context.Background()); err != nil {
		t.Fatalf("Close() error = %v, want nil", err)
	}
}

func TestDatabaseConfigDurationsCompile(t *testing.T) {
	cfg := config.DatabaseConfig{
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
		ConnectTimeout:  5 * time.Second,
	}

	if cfg.ConnectTimeout != 5*time.Second {
		t.Fatalf("ConnectTimeout = %s, want 5s", cfg.ConnectTimeout)
	}
}
