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
}
