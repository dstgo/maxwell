package counter

import (
	"github.com/dstgo/maxwell/contribs/ginx/mid/ratelimit"
	"github.com/gin-gonic/gin"
)

type Counter interface {
	Count(ctx *gin.Context, key string) (int, error)
}

type CountLimiter struct {
	Limit   int
	Window  int
	KeyFn   func(ctx *gin.Context) string
	Counter Counter
}

func (c CountLimiter) Allow(ctx *gin.Context) (bool, error) {
	key := c.KeyFn(ctx)
	count, err := c.Counter.Count(ctx, key)
	if err != nil {
		return false, err
	}
	if count > c.Limit {
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

func WithWindow(window int) Option {
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

func Limiter(options ...Option) ratelimit.Limiter {
	var opt CountLimiter
	for _, option := range options {
		option(&opt)
	}

	return
}
