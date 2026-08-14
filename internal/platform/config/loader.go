package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load(path string) (*Config, error) {
	v := viper.New()

	//指定配置文件路径
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	//设置默认值，避免配置文件漏字段时应用拿到空值
	setDefaults(v)

	// 允许环境变量覆盖配置项。
	// 例如 GOWEB_HTTP_ADDR 可以覆盖 http.addr。
	v.SetEnvPrefix("GOWEB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败：%w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败：%w", err)
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "goweb-scaffold")
	v.SetDefault("app.env", "dev")

	v.SetDefault("http.addr", ":8080")
	v.SetDefault("http.read_timeout", "5s")
	v.SetDefault("http.write_timeout", "10s")
	v.SetDefault("http.shutdown_timeout", "10s")

	v.SetDefault("log.level", "debug")
	v.SetDefault("log.format", "console")

	v.SetDefault("cors.allow_origins", []string{"http://localhost:5173", "http://localhost:3000"})
	v.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"})
	v.SetDefault("cors.expose_headers", []string{"X-Request-ID"})
	v.SetDefault("cors.allow_credentials", true)
	v.SetDefault("cors.max_age", "12h")

	v.SetDefault("security.max_body_bytes", 1048576)
	v.SetDefault("security.login_rate_limit.enabled", true)
	v.SetDefault("security.login_rate_limit.requests", 5)
	v.SetDefault("security.login_rate_limit.window", "1m")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "goweb_scaffold")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.time_zone", "Asia/Shanghai")
	v.SetDefault("database.max_open_conns", 20)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "1h")
	v.SetDefault("database.conn_max_idle_time", "30m")
	v.SetDefault("database.connect_timeout", "5s")

	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")

	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.path", "/metrics")

	v.SetDefault("telemetry.enabled", true)
	v.SetDefault("telemetry.service_name", "goweb-scaffold")
	v.SetDefault("telemetry.exporter", "stdout")
}
