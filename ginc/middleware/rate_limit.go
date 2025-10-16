package middleware

import (
	"fmt"
	"time"

	"github.com/PirateDreamer/going/ginc"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

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

type DTBLimiterParam struct {
	Client   *redis.Client             // redis 客户端
	Rate     float64                   // 每秒生成多少令牌
	Capacity float64                   // 桶容量
	BizErr   error                     // 令牌不够业务错误
	KeyFn    func(*gin.Context) string // 获取限流key的函数
}

// DTBLimiter 分布式令牌桶限流器，使用redis
func DTBLimiter(param DTBLimiterParam) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := param.KeyFn(c)
		now := float64(time.Now().UnixMilli())
		keys := []string{
			fmt.Sprintf("rate_limiter:%s:last_refill", key),
			fmt.Sprintf("rate_limiter:%s:tokens", key),
		}
		args := []interface{}{param.Rate, param.Capacity, now, 1}

		res, err := param.Client.Eval(c.Request.Context(), luaScript, keys, args...).Int()
		if err != nil {
			ginc.ResFail(c.Request.Context(), c, err)
			c.Abort()
			return
		}
		if res != 1 {
			ginc.ResFail(c.Request.Context(), c, param.BizErr)
			c.Abort()
			return
		}
		c.Next()
	}
}
