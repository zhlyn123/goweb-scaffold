package httpserver

import (
	"context"
	"errors"
	"fmt"
	"goweb-scaffold/internal/platform/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.HTTPConfig
	server *http.Server
}

func NewServer(cfg config.HTTPConfig, router *gin.Engine) *Server {
	return &Server{
		config: cfg,
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("启动 HTTP 服务失败: %w", err)
		}
		return nil
	default:
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	shutdownctx, cancel := context.WithTimeout(ctx, s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(shutdownctx); err != nil {
		return fmt.Errorf("关闭 HTTP 服务失败: %w", err)
	}

	return nil
}
