package database

import (
	"context"
	"fmt"
	"goweb-scaffold/internal/platform/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	gormDB *gorm.DB
	sqlDB  interface {
		Close() error
		PingContext(context.Context) error
	}
}

func New(ctx context.Context, cfg config.DatabaseConfig) (*DB, error) {
	gormDB, err := gorm.Open(postgres.Open(buildDSN(cfg)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接池失败: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return &DB{
		gormDB: gormDB,
		sqlDB:  sqlDB,
	}, nil
}

func (d *DB) Gorm() *gorm.DB {
	return d.gormDB
}

func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.sqlDB == nil {
		return fmt.Errorf("database is not initialized")
	}

	if err := d.sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

func (d *DB) Close(context.Context) error {
	if d == nil || d.sqlDB == nil {
		return nil
	}

	if err := d.sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接池失败: %w", err)
	}

	return nil
}

func buildDSN(cfg config.DatabaseConfig) string {
	if cfg.DSN != "" {
		return cfg.DSN
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
		cfg.TimeZone,
	)

}
