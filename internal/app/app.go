package app

import (
	"context"
	entsql "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/dstgo/maxwell/ent"
	"github.com/dstgo/maxwell/internal/app/conf"
	"github.com/dstgo/maxwell/internal/app/pkg/logh"
	"github.com/dstgo/maxwell/internal/app/types"
	"github.com/dstgo/size"
	"github.com/ginx-contribs/dbx"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/constant/methods"
	"github.com/ginx-contribs/ginx/middleware"
	"github.com/ginx-contribs/logx"
	"github.com/google/wire"
	"github.com/wneessen/go-mail"
	"log/slog"
	"net/http/pprof"
	"time"

	"github.com/dstgo/maxwell/assets"
	"github.com/gin-gonic/gin"
	_ "github.com/ginx-contribs/ent-sqlite"
	"github.com/redis/go-redis/v9"
	"io"
)

// EnvProvider only use for wire injection
var EnvProvider = wire.NewSet(
	wire.FieldsOf(new(*types.Env), "AppConf"),
	wire.FieldsOf(new(*types.Env), "Ent"),
	wire.FieldsOf(new(*types.Env), "Redis"),
	wire.FieldsOf(new(*types.Env), "Router"),
	wire.FieldsOf(new(*types.Env), "Email"),
	wire.FieldsOf(new(*conf.AppConf), "Jwt"),
	wire.FieldsOf(new(*conf.AppConf), "Email"),
)

// NewApp returns a new app server, cleanup func
func NewApp(ctx context.Context, appConf *conf.AppConf) (*ginx.Server, error) {

	slog.Debug("maxwell server is initializing")

	// initialize database
	slog.Debug(fmt.Sprintf("connecting to database(%s)", appConf.DB.Address))
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
			middleware.Recovery(slog.Default(), nil),
			middleware.Logger(slog.Default(), "access record"),
		),
	)

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
		Router:  server.RouterGroup(),
		Email:   emailClient,
	})
	if err != nil {
		return nil, err
	}

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

	return entClient, err
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

func initializeEmail(ctx context.Context, emailConf conf.EmailConf) (*mail.Client, error) {
	client, err := mail.NewClient(emailConf.Host,
		mail.WithPort(emailConf.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(emailConf.Username),
		mail.WithPassword(emailConf.Password),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
