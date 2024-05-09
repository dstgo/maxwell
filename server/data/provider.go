package data

import (
	"github.com/dstgo/maxwell/server/data/cache"
	"github.com/dstgo/maxwell/server/data/mq"
	"github.com/dstgo/maxwell/server/data/repo"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	// cache
	cache.NewRedisTokenCache,
	wire.Bind(new(cache.TokenCache), new(*cache.RedisTokenCache)),
	cache.NewRedisCodeCache,
	wire.Bind(new(cache.VerifyCodeCache), new(*cache.RedisCodeCache)),

	// repo
	repo.NewUserRepo,

	// mq
	mq.NewStreamQueue,
	wire.Bind(new(mq.Queue), new(*mq.StreamQueue)),
)
