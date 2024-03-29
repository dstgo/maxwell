package conf

import (
	"github.com/dstgo/contrib/gorms"
	"github.com/dstgo/contrib/logx"
	"time"
)

type AppConf struct {
	DB     gorms.Config `mapstructure:"db"`
	Log    logx.Options `mapstructure:"log"`
	Server *ServerConf  `mapstructure:"server"`

	Version   string
	BuildTime string
}

type ServerConf struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
	MultipartMax int64         `mapstructure:"multipartMax"`
}
