package app

import (
	"context"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/dstgo/ent-sqlite/testdata/ent"
	"github.com/dstgo/maxwell/app/api"
	"github.com/dstgo/maxwell/app/types"
	"github.com/dstgo/maxwell/conf"
	"github.com/dstgo/size"
	"github.com/ginx-contribs/dbx"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/constant/methods"
	"github.com/ginx-contribs/ginx/middleware"
	"github.com/ginx-contribs/logx"
	"log/slog"
	"time"

	// sqlite ent adapter
	"github.com/dstgo/maxwell/assets"
	"github.com/gin-gonic/gin"
	_ "github.com/ginx-contribs/ent-sqlite"
	"github.com/redis/go-redis/v9"
	"io"
)

// NewApp returns a new app server, cleanup func
func NewApp(ctx context.Context, appConf *conf.AppConf) (*ginx.Server, error) {

	prefix := "[maxwell]"

	// initialize database
	slog.Debug("connecting to database", slog.String("address", appConf.DB.Address))
	db, err := initializeDB(ctx, appConf.DB)
	if err != nil {
		return nil, err
	}

	// initialize redis client
	slog.Debug("connecting to redis", slog.String("address", appConf.Redis.Address))
	redisClient, err := initializeRedis(ctx, appConf.Redis)
	if err != nil {
		return nil, err
	}

	// initialize ginx server
	server := ginx.New(
		ginx.WithOptions(ginx.Options{
			Mode:               gin.ReleaseMode,
			LogPrefix:          prefix,
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
			middleware.Recovery(slog.Default(), nil),
			middleware.Logger(slog.Default(), prefix),
		),
	)

	// register shutdown hook
	onShutdown := func(ctx context.Context) error {
		errorLogIf("db closed failed", db.Close())
		errorLogIf("redis closed failed", redisClient.Close())
		return nil
	}
	server.OnShutdown = append(server.OnShutdown, onShutdown)

	slog.Debug("initialize api router")
	// initialize api router
	api.RegisterRouter(&types.Env{
		AppConf: appConf,
		Ent:     db,
		Redis:   redisClient,
		Router:  server.RouterGroup(),
	})

	return server, nil
}

func PrintBanner(writer io.Writer) error {
	bytes, err := assets.FS.ReadFile("banner.txt")
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}

// NewLogger returns a new app logger with the given options
func NewLogger(option conf.LogConf) (*logx.Logger, error) {

	writer, err := logx.NewWriter(&logx.WriterOptions{
		Filename: option.Filename,
	})
	if err != nil {
		return nil, err
	}
	handler, err := logx.NewHandler(writer, &logx.HandlerOptions{
		Level:       option.Level,
		Format:      option.Format,
		Prompt:      option.Prompt,
		Source:      option.Source,
		ReplaceAttr: nil,
		Color:       option.Color,
	})
	if err != nil {
		return nil, err
	}
	logger, err := logx.New(
		logx.WithHandlers(handler),
	)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func initializeDB(ctx context.Context, options dbx.Options) (*ent.Client, error) {
	sqldb, err := dbx.Open(options)
	if err != nil {
		return nil, err
	}
	entClient := ent.NewClient(
		ent.Driver(entsql.OpenDB(options.Driver, sqldb)),
	)
	// migrate database
	if err := entClient.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return nil, err
}

func initializeRedis(ctx context.Context, redisConf conf.RedisConf) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisConf.Address,
		Password:     redisConf.Password,
		ReadTimeout:  redisConf.ReadTimeout,
		WriteTimeout: redisConf.WriteTimeout,
	})
	pingResult := redisClient.Ping(ctx)
	if pingResult.Err() != nil {
		return nil, pingResult.Err()
	}

	return redisClient, nil
}

func errorLogIf(msg string, err error) {
	if err != nil {
		return
	}
	slog.Error(msg, slog.Any("error", err))
}
