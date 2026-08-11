package config

import (
	"os"
	"path"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := path.Join(dir, "config.yaml")

	content := `
app:
  name: test-app
  env: dev
http:
  addr: :9090
  read_timeout: 5s
  write_timeout: 10s
  shutdown_timeout: 10s
log:
  level: info
  format: json
cors:
  allow_origins:
    - "http://localhost:5173"
  allow_methods:
    - "GET"
    - "POST"
  allow_headers:
    - "Content-Type"
    - "Authorization"
  expose_headers:
    - "X-Request-ID"
  allow_credentials: true
  max_age: 1h
database:
  host: db.example.local
  port: 5433
  user: test_user
  password: test_password
  name: test_db
  ssl_mode: disable
  time_zone: Asia/Shanghai
  max_open_conns: 30
  max_idle_conns: 15
  conn_max_lifetime: 2h
  conn_max_idle_time: 20m
  connect_timeout: 3s
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("写入配置文件失败：%v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("加载配置文件失败：%v", err)
	}

	if cfg.App.Name != "test-app" {
		t.Fatalf("App.Name = %s, want test-app", cfg.App.Name)
	}

	if cfg.HTTP.Addr != ":9090" {
		t.Fatalf("HTTP.Addr = %s, want :9090", cfg.HTTP.Addr)
	}

	if cfg.Log.Level != "info" {
		t.Fatalf("Log.Level = %s, want info", cfg.Log.Level)
	}
	if len(cfg.CORS.AllowOrigins) != 1 || cfg.CORS.AllowOrigins[0] != "http://localhost:5173" {
		t.Fatalf("CORS.AllowOrigins = %v, want [http://localhost:5173]", cfg.CORS.AllowOrigins)
	}

	if cfg.CORS.MaxAge.String() != "1h0m0s" {
		t.Fatalf("CORS.MaxAge = %s, want 1h0m0s", cfg.CORS.MaxAge)
	}

	if cfg.Database.Host != "db.example.local" {
		t.Fatalf("Database.Host = %s, want db.example.local", cfg.Database.Host)
	}

	if cfg.Database.Port != 5433 {
		t.Fatalf("Database.Port = %d, want 5433", cfg.Database.Port)
	}

	if cfg.Database.User != "test_user" {
		t.Fatalf("Database.User = %s, want test_user", cfg.Database.User)
	}

	if cfg.Database.Name != "test_db" {
		t.Fatalf("Database.Name = %s, want test_db", cfg.Database.Name)
	}

	if cfg.Database.MaxOpenConns != 30 {
		t.Fatalf("Database.MaxOpenConns = %d, want 30", cfg.Database.MaxOpenConns)
	}

	if cfg.Database.ConnectTimeout.String() != "3s" {
		t.Fatalf("Database.ConnectTimeout = %s, want 3s", cfg.Database.ConnectTimeout)
	}
}
