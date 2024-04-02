package mid

import (
	"fmt"
	"github.com/dstgo/maxwell/contribs/ginx/resp"
	"github.com/dstgo/maxwell/contribs/utils/str2byte"
	"github.com/dstgo/size"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// RequestIdHandler generates a unique request id for each http request
func RequestIdHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("X-Request-ID", uuid.NewString())
		ctx.Next()
	}
}

// AccessLogHandler records each http request into log
func AccessLogHandler(logger *slog.Logger, msg string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now()

		ctx.Next()

		cost := time.Now().Sub(begin)
		method := ctx.Request.Method
		url := ctx.Request.RequestURI
		path := ctx.FullPath()
		status := ctx.Writer.Status()
		ip := ctx.ClientIP()
		requestId := ctx.Writer.Header().Get("X-Request-ID")
		requestSize := roundSize(int(ctx.Request.ContentLength))
		responseSize := roundSize(ctx.Writer.Size())

		attrs := []any{
			slog.String("method", method), slog.Int("status", status),
			slog.String("cost", fmt.Sprintf("%.2fms", float64(cost)/float64(time.Millisecond))),
			slog.String("ip", ip), slog.String("url", url), slog.String("path", path),
			slog.String("request-size", requestSize.String()), slog.String("response-size", responseSize.String()),
		}

		if requestId != "" {
			attrs = append(attrs, slog.String("request-id", requestId))
		}

		if len(ctx.Errors) != 0 {
			attrs = append(attrs, slog.String("errors", ctx.Errors.String()))
		}

		logger.InfoContext(ctx, msg, attrs...)
	}
}

func roundSize(s int) size.Size {
	if s < 0 {
		return size.NewInt(0, size.B)
	}
	bodysize := size.NewInt(s, size.B)
	if size.Unit(s) > size.KB {
		bodysize = bodysize.To(size.KB)
	} else if size.Unit(s) > size.MB {
		bodysize = bodysize.To(size.MB)
	} else if size.Unit(s) > size.GB {
		bodysize = bodysize.To(size.GB)
	}
	return bodysize
}

// RecoveryHandler handles the case of panic situation
func RecoveryHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if panicErr := recover(); panicErr != nil {
				var (
					brokenPipe bool
				)

				var err any
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				if ne, ok := panicErr.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				err = panicErr

				if !brokenPipe {
					logger.ErrorContext(ctx, "[Panic Recovered]", slog.Any("error", err), slog.String("stack", str2byte.BytesToString(debug.Stack())))
					ctx.Abort()
					resp.Fail(ctx).Status(http.StatusInternalServerError).Msg("internal server error").Render()
				} else {
					return
				}
			}
		}()

		ctx.Next()
	}
}
