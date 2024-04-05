package types

import (
	"github.com/dstgo/ent-sqlite/testdata/ent"
	"github.com/dstgo/maxwell/conf"
	"github.com/ginx-contribs/ginx"
	"github.com/redis/go-redis/v9"
)

// Response just use for documentation
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type Env struct {
	AppConf *conf.AppConf
	Ent     *ent.Client
	Redis   *redis.Client
	Router  *ginx.RouterGroup
}
