package mids

import (
	"errors"
	"github.com/dstgo/maxwell/server/types/auth"
	"github.com/dstgo/maxwell/server/types/route"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/contribs/ratelimit"
	"github.com/ginx-contribs/ginx/contribs/ratelimit/counter"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/redis/go-redis/v9"
	"time"
)

func ByIpPath(ctx *gin.Context) string {
	return ctx.RemoteIP() + ":" + ctx.FullPath()
}

// RateLimitByCount limits the number of requests by counting for public api
func RateLimitByCount(client *redis.Client, limit int, window time.Duration, keyFn func(ctx *gin.Context) string) gin.HandlerFunc {
	limiter := counter.NewLimiter(
		counter.WithLimit(limit),
		counter.WithWindow(window),
		counter.WithKeyFn(keyFn),
		counter.WithCounter(counter.Redis(client)),
	)

	return func(ctx *gin.Context) {
		metadata := ginx.MetaFromCtx(ctx)
		if !metadata.Contains(route.Public) {
			ctx.Next()
			return
		}

		// try allowing requests
		allow, err := limiter.Allow(ctx)
		if err == nil {
			ctx.Next()
			allow()
		} else {
			if errors.Is(err, ratelimit.ErrRateLimitExceed) {
				resp.Fail(ctx).Error(auth.ErrRateLimitExceeded).JSON()
			} else {
				resp.InternalError(ctx).Error(err).JSON()
			}
			ctx.Abort()
		}
	}
}
