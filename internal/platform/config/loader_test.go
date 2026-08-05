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
}
