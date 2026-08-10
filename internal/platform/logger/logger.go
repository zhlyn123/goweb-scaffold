package logger

import (
	"fmt"
	"goweb-scaffold/internal/platform/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LogConfig) (*zap.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	zapCfg := buildZapConfig(cfg, level)

	log, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("创建日志对象失败：%w", err)
	}

	return log, nil
}

// parseLevel 将配置里的日志级别转换为 zap 使用的日志级别
func parseLevel(level string) (zapcore.Level, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return zapcore.InfoLevel, fmt.Errorf("无效的日志级别：%s", level)
	}

	return zapLevel, nil
}

// buildZapConfig 根据应用配置构造zap配置
func buildZapConfig(cfg config.LogConfig, level zapcore.Level) zap.Config {
	if cfg.Format == "json" {
		zapCfg := zap.NewProductionConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
		zapCfg.Encoding = "json"
		zapCfg.EncoderConfig.TimeKey = "ts"
		zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		return zapCfg
	}

	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.Encoding = "console"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapCfg
}
