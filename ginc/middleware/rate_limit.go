package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/PirateDreamer/going/ginc"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// TokenBucketLimiter 实现令牌桶限流
type TokenBucketLimiter struct {
	client   *redis.Client // redis 客户端
	rate     float64       // 每秒生成多少令牌
	capacity float64       // 桶容量
	bizErr   error         // 令牌不够业务错误
}

// Lua脚本（保证原子性）
const luaScript = `
local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local last_refill = tonumber(redis.call("get", KEYS[1]) or 0)
local tokens = tonumber(redis.call("get", KEYS[2]) or capacity)

local delta = math.max(0, now - last_refill)
local add_tokens = delta * rate / 1000
tokens = math.min(capacity, tokens + add_tokens)

redis.call("set", KEYS[1], now)

if tokens >= requested then
    tokens = tokens - requested
    redis.call("set", KEYS[2], tokens)
    return 1
else
    redis.call("set", KEYS[2], tokens)
    return 0
end
`

// NewTokenBucketLimiter 创建限流器实例
func NewTokenBucketLimiter(client *redis.Client, rate, capacity float64, bizErr error) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		client:   client,
		rate:     rate,
		capacity: capacity,
		bizErr:   bizErr,
	}
}

// Allow 检查是否允许通过
func (l *TokenBucketLimiter) Allow(ctx context.Context, key string, requested float64) (bool, error) {
	now := float64(time.Now().UnixMilli())
	keys := []string{
		fmt.Sprintf("rate_limiter:%s:last_refill", key),
		fmt.Sprintf("rate_limiter:%s:tokens", key),
	}
	args := []interface{}{l.rate, l.capacity, now, requested}

	res, err := l.client.Eval(ctx, luaScript, keys, args...).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// GinMiddleware 返回一个 Gin 中间件
func (l *TokenBucketLimiter) GinMiddleware(keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFunc(c)
		allowed, err := l.Allow(c.Request.Context(), key, 1)
		if err != nil {
			ginc.ResFail(c.Request.Context(), c, err)
			return
		}
		if !allowed {
			ginc.ResFail(c.Request.Context(), c, l.bizErr)
			c.Abort()
			return
		}
		c.Next()
	}
}
