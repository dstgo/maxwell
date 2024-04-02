package app

import (
	"context"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/dstgo/ent-sqlite/testdata/ent"
	"github.com/dstgo/maxwell/app/api"
	"github.com/dstgo/maxwell/app/types"
	"github.com/dstgo/maxwell/conf"
	"github.com/dstgo/maxwell/contribs/db"
	"github.com/dstgo/maxwell/contribs/ginx/mid"
	"github.com/dstgo/maxwell/contribs/logx"

	// sqlite ent adapter
	_ "github.com/dstgo/ent-sqlite"
	"github.com/dstgo/maxwell/assets"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"io"
	"log/slog"
	"net/http"
)

// NewApp returns a new app server, cleanup func
func NewApp(ctx context.Context, appConf *conf.AppConf) (*http.Server, error) {

	// initialize logger first
	logger, err := logx.New(
		logx.WithOutput(appConf.Log.Output),
		logx.WithLevel(appConf.Log.Level),
		logx.WithFormat(appConf.Log.Format),
		logx.WithSource(appConf.Log.Source),
	)
	if err != nil {
		return nil, err
	}

	// set global logger
	slog.SetDefault(logger.Slog())
	if err := printBanner(logger.Writer()); err != nil {
		return nil, err
	}
	slog.Info("server started successfully!", slog.String("buildtime", appConf.BuildTime), slog.String("version", appConf.Version))

	// initialize sql database
	sqldb, err := db.Open(appConf.DB)
	if err != nil {
		return nil, err
	}
	entClient := ent.NewClient(
		ent.Driver(entsql.OpenDB(appConf.DB.Driver, sqldb)),
	)
	// migrate database
	if err := entClient.Schema.Create(ctx); err != nil {
		return nil, err
	}

	// initialize redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:         appConf.Redis.Address,
		Password:     appConf.Redis.Password,
		ReadTimeout:  appConf.Redis.ReadTimeout,
		WriteTimeout: appConf.Redis.WriteTimeout,
	})
	pingresult := redisClient.Ping(ctx)
	if pingresult.Err() != nil {
		return nil, pingresult.Err()
	}

	// initialize http server
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(
		mid.RecoveryHandler(logger.Slog()),
		mid.AccessLogHandler(logger.Slog(), "accesslog"),
	)
	engine.NoRoute(mid.ResourceNotFoundHandler())
	engine.NoMethod(mid.MethodAllowHandler())
	engine.MaxMultipartMemory = appConf.Server.MultipartMax
	server := &http.Server{
		Handler:      engine,
		Addr:         appConf.Server.Address,
		ReadTimeout:  appConf.Server.ReadTimeout,
		WriteTimeout: appConf.Server.WriteTimeout,
		IdleTimeout:  appConf.Server.IdleTimeout,
	}

	// register router
	api.RegisterRouter(&types.Env{
		AppConf: appConf,
		Ent:     entClient,
		Redis:   redisClient,
	})

	// cleanup will be called when server shutdown
	server.RegisterOnShutdown(func() {
		entClient.Close()
		logger.Close()
	})

	return server, nil
}

func printBanner(writer io.Writer) error {
	bytes, err := assets.FS.ReadFile("banner.txt")
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}
