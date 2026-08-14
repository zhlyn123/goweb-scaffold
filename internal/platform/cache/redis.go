package cache

import (
	"context"
	"fmt"

	"goweb-scaffold/internal/platform/config"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg config.RedisConfig) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           cfg.DB,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}),
	}
}

func (r *Redis) Ping(ctx context.Context) error {
	if r == nil || r.client == nil {
		return fmt.Errorf("redis is not initialized")
	}

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	return nil
}

func (r *Redis) Close(context.Context) error {
	if r == nil || r.client == nil {
		return nil
	}

	if err := r.client.Close(); err != nil {
		return fmt.Errorf("close redis: %w", err)
	}

	return nil
}
