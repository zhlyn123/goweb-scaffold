package bootstrap

import (
	"context"
	"fmt"

	"goweb-scaffold/internal/app/lifecycle"
	"goweb-scaffold/internal/platform/config"
)

// Application 负责应用最外层的启动和关闭流程。
type Application struct {
	config    *config.Config
	lifecycle *lifecycle.Manager
}

// NewApplication 创建应用根对象。
//
// 第一阶段刻意保持简单。后续阶段会在这里注册配置、日志、数据库、
// 缓存、HTTP 服务、指标监控和业务模块路由。
func NewApplication(config *config.Config) *Application {
	return &Application{
		config:    config,
		lifecycle: lifecycle.NewManager(),
	}
}

// Run 启动所有已注册的应用组件，并在上下文取消或启动失败时执行关闭流程。
func (a *Application) Run(ctx context.Context) error {
	if err := a.lifecycle.Start(ctx); err != nil {
		return fmt.Errorf("start application: %w", err)
	}

	<-ctx.Done()

	if err := a.lifecycle.Stop(context.Background()); err != nil {
		return fmt.Errorf("stop application: %w", err)
	}

	return nil
}
