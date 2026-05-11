package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	infraredis "github.com/bing127/enterprise-agent/module-infra/middleware/redis"
	"github.com/bing127/enterprise-agent/module-pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// ── 周期性固定窗口限流（内存） ───────────────────────────────────────────────────

// RateLimitMiddleware 基于 IP 的固定窗口限流中间件。
//
//   - limit:  窗口内最大请求数
//   - window: 窗口时长（如 time.Minute）
func RateLimitMiddleware(limit int, window time.Duration) app.HandlerFunc {
	type entry struct {
		count   int
		resetAt time.Time
	}
	var (
		mu    sync.Mutex
		store = make(map[string]*entry)
	)
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		mu.Lock()
		e, ok := store[ip]
		if !ok || time.Now().After(e.resetAt) {
			store[ip] = &entry{count: 1, resetAt: time.Now().Add(window)}
			mu.Unlock()
			c.Next(ctx)
			return
		}
		if e.count >= limit {
			mu.Unlock()
			response.TooManyRequests(ctx, c, "rate limit exceeded, please try again later")
			c.Abort()
			return
		}
		e.count++
		mu.Unlock()
		c.Next(ctx)
	}
}

// ── Redis 令牌桶限流 ────────────────────────────────────────────────────────────

// tokenBucketScript 是 Redis Lua 脚本，实现令牌桶算法。
// KEYS[1] = bucket key (e.g. "ratelimit:<ip>")
// ARGV[1] = capacity（桶容量）
// ARGV[2] = rate（每秒补充令牌数，浮点数）
// ARGV[3] = now（当前 Unix 时间戳，秒，浮点数）
// 返回值：1 = 允许，0 = 拒绝
const tokenBucketScript = `
local key       = KEYS[1]
local capacity  = tonumber(ARGV[1])
local rate      = tonumber(ARGV[2])
local now       = tonumber(ARGV[3])

local data = redis.call("HMGET", key, "tokens", "last_time")
local tokens    = tonumber(data[1]) or capacity
local last_time = tonumber(data[2]) or now

-- 按时间差补充令牌
local elapsed = math.max(0, now - last_time)
tokens = math.min(capacity, tokens + elapsed * rate)

if tokens >= 1 then
    tokens = tokens - 1
    redis.call("HMSET", key, "tokens", tokens, "last_time", now)
    redis.call("EXPIRE", key, math.ceil(capacity / rate) + 1)
    return 1
else
    redis.call("HMSET", key, "tokens", tokens, "last_time", now)
    redis.call("EXPIRE", key, math.ceil(capacity / rate) + 1)
    return 0
end
`

// RedisTokenBucketMiddleware 基于 Redis 令牌桶的分布式限流中间件。
//
//   - rdb:      infra Redis 客户端
//   - keyPrefix: 限流 key 前缀（如 "api:rl"）
//   - capacity:  桶容量（最大突发请求数）
//   - rate:      每秒补充令牌数（平均 QPS 上限）
func RedisTokenBucketMiddleware(rdb *infraredis.Client, keyPrefix string, capacity int64, rate float64) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ip := c.ClientIP()
		key := fmt.Sprintf("%s:%s", keyPrefix, ip)
		now := float64(time.Now().UnixNano()) / 1e9 // Unix 秒（高精度）

		result, err := rdb.RedisClient.Eval(ctx, tokenBucketScript,
			[]string{key},
			capacity, rate, now,
		).Int64()
		if err != nil {
			// Redis 故障时放行（降级策略，避免 Redis 问题影响业务）
			c.Next(ctx)
			return
		}

		if result == 0 {
			response.TooManyRequests(ctx, c, "too many requests, please slow down")
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
