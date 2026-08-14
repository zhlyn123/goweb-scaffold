package middleware

import (
	"fmt"
	"goweb-scaffold/internal/platform/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(appEnv string, cfg config.CORSConfig) (gin.HandlerFunc, error) {
	if appEnv == "prod" && cfg.AllowCredentials && hasWildcardOrigin(cfg.AllowOrigins) {
		return nil, fmt.Errorf("生产环境不允许同时开启 allow_credentials 和通配符 CORS origin")
	}

	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		ExposeHeaders:    cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.MaxAge,
	}), nil
}

func hasWildcardOrigin(origins []string) bool {
	for _, origin := range origins {
		if origin == "*" {
			return true
		}
	}

	return false
}
