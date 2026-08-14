package bootstrap

import (
	"context"
	"fmt"
	"os"

	"goweb-scaffold/internal/app/lifecycle"
	authapp "goweb-scaffold/internal/modules/auth/application"
	authhttp "goweb-scaffold/internal/modules/auth/interfaces/http"
	rbacapp "goweb-scaffold/internal/modules/rbac/application"
	rbacpersistence "goweb-scaffold/internal/modules/rbac/infrastructure/persistence"
	systemhttp "goweb-scaffold/internal/modules/system/interfaces/http"
	userpersistence "goweb-scaffold/internal/modules/user/infrastructure/persistence"
	userhttp "goweb-scaffold/internal/modules/user/interfaces/http"
	"goweb-scaffold/internal/platform/cache"
	"goweb-scaffold/internal/platform/config"
	"goweb-scaffold/internal/platform/database"
	"goweb-scaffold/internal/platform/httpserver"
	"goweb-scaffold/internal/platform/httpserver/middleware"
	"goweb-scaffold/internal/platform/logger"
	"goweb-scaffold/internal/platform/metrics"
	"goweb-scaffold/internal/platform/security"
	"goweb-scaffold/internal/platform/telemetry"
	"goweb-scaffold/internal/platform/transaction"
	"goweb-scaffold/internal/shared/idgen"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

type Application struct {
	config    *config.Config
	logger    *zap.Logger
	db        *database.DB
	redis     *cache.Redis
	lifecycle *lifecycle.Manager
}

func NewApplication(cfg *config.Config) (*Application, error) {
	log, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	manager := lifecycle.NewManager()

	db, err := database.New(context.Background(), cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("create database: %w", err)
	}
	manager.Register(lifecycle.Hook{Name: "database", Stop: db.Close})

	redisClient := cache.NewRedis(cfg.Redis)
	manager.Register(lifecycle.Hook{Name: "redis", Stop: redisClient.Close})

	if cfg.Telemetry.Enabled {
		telemetryProvider, err := telemetry.New(context.Background(), cfg.Telemetry)
		if err != nil {
			return nil, fmt.Errorf("create telemetry: %w", err)
		}
		manager.Register(lifecycle.Hook{Name: "telemetry", Stop: telemetryProvider.Shutdown})
	}

	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	corsMiddleware, err := middleware.CORS(cfg.App.Env, cfg.CORS)
	if err != nil {
		return nil, fmt.Errorf("create cors middleware: %w", err)
	}

	router := gin.New()
	router.Use(middleware.RequestID())
	if cfg.Telemetry.Enabled {
		router.Use(otelgin.Middleware(cfg.Telemetry.ServiceName))
	}
	router.Use(
		middleware.Recovery(log),
		middleware.AccessLog(log),
		middleware.BodySizeLimit(cfg.Security.MaxBodyBytes),
		corsMiddleware,
	)

	if cfg.Metrics.Enabled {
		metricSet := metrics.New()
		router.Use(metricSet.Middleware())
		router.GET(cfg.Metrics.Path, gin.WrapH(metricSet.Handler()))
	}
	registerOpenAPI(router)

	systemHandler := systemhttp.NewHandler(db, redisClient)
	systemhttp.RegisterRoutes(router, systemHandler)

	userRepo := userpersistence.NewGormRepository(db.Gorm())
	passwordHasher := security.NewPasswordHasher()
	tokenService := security.NewTokenService(&cfg.JWT)
	idGenerator := idgen.NewUUIDGenerator()
	txManager := transaction.NewGormManager(db.Gorm())

	rbacRepo := rbacpersistence.NewGormRepository(db.Gorm())
	rbacUsecase := rbacapp.NewUsecase(rbacRepo, idGenerator)

	authUsecase := authapp.NewUsecase(
		userRepo,
		passwordHasher,
		tokenService,
		idGenerator,
		txManager,
		authapp.RoleBinderFunc(func(ctx context.Context, userID string, roleCode string) error {
			return rbacUsecase.BindRoleToUser(ctx, rbacapp.BindRoleToUserCommand{
				UserID:   userID,
				RoleCode: roleCode,
			})
		}),
		"user",
	)

	authHandler := authhttp.NewHandler(authUsecase)
	authhttp.RegisterRoutes(router, authHandler, middleware.LoginRateLimit(cfg.Security.LoginRateLimit))

	authMiddleware := middleware.AuthRequired(tokenService)
	requireUserReadPermission := middleware.RequirePermission(
		func(ctx context.Context, userID string, permissionCode string) error {
			return rbacUsecase.CheckPermission(ctx, rbacapp.CheckPermissionCommand{
				UserID:         userID,
				PermissionCode: permissionCode,
			})
		},
		"user:read",
	)

	userHandler := userhttp.NewHandler(userRepo)
	userhttp.RegisterRoutes(router, userHandler, authMiddleware, requireUserReadPermission)

	server := httpserver.NewServer(cfg.HTTP, router)
	manager.Register(lifecycle.Hook{Name: "http_server", Start: server.Start, Stop: server.Stop})

	return &Application{
		config:    cfg,
		logger:    log,
		db:        db,
		redis:     redisClient,
		lifecycle: manager,
	}, nil
}

func registerOpenAPI(router *gin.Engine) {
	for _, path := range []string{
		"api/openapi/openapi.yaml",
		"api-docs/openapi/openapi.yaml",
	} {
		if _, err := os.Stat(path); err == nil {
			router.StaticFile("/openapi.yaml", path)
			return
		}
	}
}

func (a *Application) Run(ctx context.Context) error {
	a.logger.Info(
		"application starting",
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
		return fmt.Errorf("sync logger: %w", err)
	}

	return nil
}
