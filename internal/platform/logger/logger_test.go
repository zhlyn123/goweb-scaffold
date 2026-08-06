package logger

import (
	"testing"

	"go.uber.org/zap"

	"goweb-scaffold/internal/platform/config"
)

func TestNew(t *testing.T) {
	log, err := New(config.LogConfig{
		Level:  "debug",
		Format: "console",
	})
	if err != nil {
		t.Fatalf("创建日志对象失败: %v", err)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Debug("测试 debug 日志")
	log.Info("测试 info 日志")
}

func TestNewWithJSONFormat(t *testing.T) {
	log, err := New(config.LogConfig{
		Level:  "info",
		Format: "json",
	})
	if err != nil {
		t.Fatalf("创建 JSON 日志对象失败: %v", err)
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("测试 json 日志", zap.String("component", "logger"))
}

func TestNewWithInvalidLevel(t *testing.T) {
	_, err := New(config.LogConfig{
		Level:  "invalid",
		Format: "console",
	})
	if err == nil {
		t.Fatal("期望返回错误，实际没有错误")
	}
}
