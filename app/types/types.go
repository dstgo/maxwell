package types

import (
	"github.com/dstgo/maxwell/app/conf"
	"github.com/dstgo/maxwell/app/data/ent"
	"github.com/redis/go-redis/v9"
)

type Env struct {
	AppConf *conf.AppConf
	Ent     *ent.Client
	Redis   *redis.Client
}
