package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"goweb-scaffold/internal/platform/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	configPath := flag.String("config", "configs/config.local.yaml", "配置文件路径")
	migrationsDir := flag.String("dir", "migrations", "migration 文件目录")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "缺少 migration 命令，例如: up, down, status")
		os.Exit(1)
	}

	command := flag.Arg(0)

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("pgx", buildDSN(cfg.Database))
	if err != nil {
		fmt.Fprintf(os.Stderr, "打开数据库连接失败: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Fprintf(os.Stderr, "设置 goose 数据库方言失败: %v\n", err)
		os.Exit(1)
	}

	if err := goose.Run(command, db, *migrationsDir, flag.Args()[1:]...); err != nil {
		fmt.Fprintf(os.Stderr, "执行 migration 失败: %v\n", err)
		os.Exit(1)
	}
}

// buildDSN 根据数据库配置构造 PostgreSQL 连接字符串。
// 如果配置中直接提供了 DSN，则优先使用 DSN。
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