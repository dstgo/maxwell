package ratelimit

import (
	"github.com/dstgo/maxwell/contribs/ginx/pkg/ratelimit/counter"
	"github.com/dstgo/maxwell/contribs/ginx/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Limiter limits rate of requests to server
type Limiter interface {
	Allow(ctx *gin.Context) (bool, error)
}

type Options struct {
	Limiter Limiter
	Handler gin.HandlerFunc
}

type Option func(*Options)

func WithLimiter(limiter Limiter) Option {
	return func(option *Options) {
		option.Limiter = limiter
	}
}

func WithHandler(handler gin.HandlerFunc) Option {
	return func(option *Options) {
		option.Handler = handler
	}
}

// RateLimit returns a new limiter handler with options
func RateLimit(opts ...Option) gin.HandlerFunc {
	opt := new(Options)
	for _, option := range opts {
		option(opt)
	}

	// default counter limiter
	if opt.Limiter == nil {
		opt.Limiter = counter.Limiter()
	}

	return func(ctx *gin.Context) {
		if opt.Limiter == nil {
			ctx.Next()
			return
		}

		allow, err := opt.Limiter.Allow(ctx)
		if err != nil {
			resp.Fail(ctx).Status(http.StatusInternalServerError).Render()
			ctx.Abort()
			return
		}

		if !allow {
			if opt.Handler != nil {
				opt.Handler(ctx)
			}
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		} else {
			ctx.Next()
		}
	}
}
