package server

import (
	"context"
	"fmt"
	"github.com/dstgo/maxwell/server/conf"
	authhandler "github.com/dstgo/maxwell/server/handler/auth"
	"github.com/dstgo/maxwell/server/mids"
	"github.com/dstgo/maxwell/server/pkg/logh"
	"github.com/dstgo/maxwell/server/types"
	"github.com/dstgo/size"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/constant/methods"
	"github.com/ginx-contribs/ginx/middleware"
	"log/slog"
	"net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/ginx-contribs/ent-sqlite"
)

// NewApp returns a new http app server
func NewApp(ctx context.Context, appConf *conf.App) (*ginx.Server, error) {

	slog.Debug("maxwell server is initializing")

	// initialize database
	slog.Debug(fmt.Sprintf("connecting to %s(%s)", appConf.DB.Driver, appConf.DB.Address))
	db, err := initializeDB(ctx, appConf.DB)
	if err != nil {
		return nil, err
	}

	// initialize redis client
	slog.Debug(fmt.Sprintf("connecting to redis(%s)", appConf.Redis.Address))
	redisClient, err := initializeRedis(ctx, appConf.Redis)
	if err != nil {
		return nil, err
	}

	// initialize email client
	slog.Debug(fmt.Sprintf("establish email client(%s:%d)", appConf.Email.Host, appConf.Email.Port))
	emailClient, err := initializeEmail(ctx, appConf.Email)
	if err != nil {
		return nil, err
	}

	// initialize ginx server
	server := ginx.New(
		ginx.WithOptions(ginx.Options{
			Mode:               gin.ReleaseMode,
			ReadTimeout:        appConf.Server.ReadTimeout,
			WriteTimeout:       appConf.Server.WriteTimeout,
			IdleTimeout:        appConf.Server.IdleTimeout,
			MaxMultipartMemory: appConf.Server.MultipartMax,
			MaxHeaderBytes:     int(size.MB * 2),
			MaxShutdownTimeout: time.Second * 5,
		}),
		ginx.WithNoMethod(middleware.NoMethod(methods.Get, methods.Post, methods.Put, methods.Delete, methods.Options)),
		ginx.WithNoRoute(middleware.NoRoute()),
		ginx.WithMiddlewares(
			// reocvery handler
			middleware.Recovery(slog.Default(), nil),
			// access logger
			middleware.Logger(slog.Default(), "accesslog"),
			// rate limit by counting
			mids.RateLimitByCount(redisClient, appConf.Limit.Public.Limit, appConf.Limit.Public.Window, mids.ByIpPath),
			// jwt authentication
			mids.TokenAuthenticator(authhandler.NewTokenHandler(appConf.Jwt, redisClient)),
		),
	)
	ginx.DefaultValidateHandler = validatePramsHandler

	// whether to enable pprof program profiling
	if appConf.Server.Pprof {
		server.Engine().GET("/pprof/profile", gin.WrapF(pprof.Profile))
		server.Engine().GET("/pprof/heap", gin.WrapH(pprof.Handler("heap")))
		server.Engine().GET("/pprof/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		slog.Info("pprof profiling enabled")
	}

	// register shutdown hook
	onShutdown := func(ctx context.Context) error {
		logh.ErrorNotNil("db closed failed", db.Close())
		logh.ErrorNotNil("redis closed failed", redisClient.Close())
		return nil
	}
	server.OnShutdown = append(server.OnShutdown, onShutdown)

	slog.Debug("setup api router")
	// initialize api router
	_, err = setup(&types.Env{
		AppConf: appConf,
		Ent:     db,
		Redis:   redisClient,
		Router:  server.RouterGroup().Group("/api"),
		Email:   emailClient,
	})
	if err != nil {
		return nil, err
	}

	return server, nil
}
