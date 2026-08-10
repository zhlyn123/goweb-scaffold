package bootstrap

import (
	"context"
	"fmt"

	"goweb-scaffold/internal/app/lifecycle"
	systemhttp "goweb-scaffold/internal/modules/system/interfaces/http"
	"goweb-scaffold/internal/platform/config"
	"goweb-scaffold/internal/platform/httpserver"
	"goweb-scaffold/internal/platform/httpserver/middleware"
	"goweb-scaffold/internal/platform/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Application 负责应用最外层的启动和关闭流程。
type Application struct {
	config    *config.Config
	logger    *zap.Logger
	lifecycle *lifecycle.Manager
}

// NewApplication 创建应用根对象。
//
// 第一阶段刻意保持简单。后续阶段会在这里注册配置、日志、数据库、
// 缓存、HTTP 服务、指标监控和业务模块路由。
func NewApplication(cfg *config.Config) (*Application, error) {
	log, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.Recovery(log),
		middleware.AccessLog(log),
		middleware.CORS(cfg.CORS),
	)

	systemHandler := systemhttp.NewHandler()
	systemhttp.RegisterRoutes(router, systemHandler)

	manger := lifecycle.NewManager()
	server := httpserver.NewServer(cfg.HTTP, router)

	manger.Register(
		lifecycle.Hook{
			Name:  "http_server",
			Start: server.Start,
			Stop:  server.Stop,
		},
	)

	return &Application{
		config:    cfg,
		logger:    log,
		lifecycle: manger,
	}, nil
}

// Run 启动所有已注册的应用组件，并在上下文取消或启动失败时执行关闭流程。
func (a *Application) Run(ctx context.Context) error {
	a.logger.Info(
		"应用启动中",
		zap.String("app_name", a.config.App.Name),
		zap.String("app_env", a.config.App.Env),
		zap.String("http_addr", a.config.HTTP.Addr),
	)

	if err := a.lifecycle.Start(ctx); err != nil {
		return fmt.Errorf("start application: %w", err)
	}

	<-ctx.Done()

	if err := a.lifecycle.Stop(context.Background()); err != nil {
		return fmt.Errorf("stop application: %w", err)
	}

	if err := a.logger.Sync(); err != nil {
		return fmt.Errorf("刷新日志失败: %w", err)
	}

	return nil
}
