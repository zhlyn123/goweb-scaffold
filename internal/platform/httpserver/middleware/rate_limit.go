package middleware

import (
	"goweb-scaffold/internal/platform/config"
	"goweb-scaffold/internal/shared/response"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitRecord struct {
	count     int
	resetAt   time.Time
	updatedAt time.Time
}

// LoginRateLimit 按客户端 IP 对登录接口做基础限流。
// 第一阶段使用进程内存实现，后续可以替换为 Redis 分布式限流。
func LoginRateLimit(cfg config.RateLimitConfig) gin.HandlerFunc {
	if !cfg.Enabled || cfg.Requests <= 0 || cfg.Window <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	var mu sync.Mutex
	records := make(map[string]rateLimitRecord)

	return func(c *gin.Context) {
		now := time.Now()
		ip := c.ClientIP()

		mu.Lock()
		cleanupExpiredRateLimitRecords(records, now)

		record := records[ip]
		if record.resetAt.IsZero() || now.After(record.resetAt) {
			record = rateLimitRecord{
				count:     0,
				resetAt:   now.Add(cfg.Window),
				updatedAt: now,
			}
		}

		if record.count >= cfg.Requests {
			retryAfter := int(time.Until(record.resetAt).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			records[ip] = record
			mu.Unlock()

			c.Header("Retry-After", strconv.Itoa(retryAfter))
			response.Error(c, http.StatusTooManyRequests, "RATE_LIMITED", "请求过于频繁")
			c.Abort()
			return
		}

		record.count++
		record.updatedAt = now
		records[ip] = record
		mu.Unlock()

		c.Next()
	}
}

// cleanupExpiredRateLimitRecords 清理已经过期的限流记录，避免内存持续增长。
func cleanupExpiredRateLimitRecords(records map[string]rateLimitRecord, now time.Time) {
	for ip, record := range records {
		if now.After(record.resetAt) {
			delete(records, ip)
		}
	}
}
