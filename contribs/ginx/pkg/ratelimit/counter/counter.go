package counter

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"time"
)

type Counter interface {
	Count(ctx context.Context, key string, limit int, window time.Duration) (int, error)
}

type CountLimiter struct {
	Limit   int
	Window  time.Duration
	KeyFn   func(ctx *gin.Context) string
	Counter Counter
}

func (c CountLimiter) Allow(ctx *gin.Context) (bool, error) {
	key := c.KeyFn(ctx)
	count, err := c.Counter.Count(ctx, key, c.Limit, c.Window)
	if err != nil {
		return false, err
	}
	if count >= c.Limit {
		return false, nil
	}
	return true, nil
}

type Option func(options *CountLimiter)

func WithLimit(limit int) Option {
	return func(cl *CountLimiter) {
		cl.Limit = limit
	}
}

func WithWindow(window time.Duration) Option {
	return func(cl *CountLimiter) {
		cl.Window = window
	}
}

func WithKeyFn(keyFn func(ctx *gin.Context) string) Option {
	return func(options *CountLimiter) {
		options.KeyFn = keyFn
	}
}

func WithCounter(counter Counter) Option {
	return func(cl *CountLimiter) {
		cl.Counter = counter
	}
}

func ClientIpKey() func(ctx *gin.Context) string {
	return func(ctx *gin.Context) string {
		return ctx.ClientIP()
	}
}

func UrlKey() func(ctx *gin.Context) string {
	return func(ctx *gin.Context) string {
		return ctx.Request.RequestURI
	}
}

func Limiter(options ...Option) CountLimiter {
	var limiter CountLimiter
	for _, option := range options {
		option(&limiter)
	}

	if limiter.Counter == nil {
		limiter.Counter = Cache()
	}

	if limiter.KeyFn == nil {
		limiter.KeyFn = UrlKey()
	}

	if limiter.Limit == 0 {
		limiter.Limit = 100
	}

	if limiter.Window == 0 {
		limiter.Window = time.Minute
	}

	return limiter
}
