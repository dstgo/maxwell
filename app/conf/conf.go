package conf

import (
	"github.com/dstgo/contrib/gorms"
	"github.com/dstgo/contrib/logx"
	"time"
)

type AppConf struct {
	Server ServerConf   `mapstructure:"server"`
	Log    logx.Options `mapstructure:"log"`
	DB     gorms.Config `mapstructure:"db"`
	Redis  RedisConf    `mapstructure:"redis"`

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

type RedisConf struct {
	Address      string        `mapstructure:"address"`
	Password     string        `mapstructure:"password"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
}

type JwtConf struct {
	Issuer string `mapstructure:"issuer"`
	Access struct {
		Expire time.Duration `mapstructure:"expire"`
		Delay  time.Duration `mapstructure:"delay"`
		Key    string        `mapstructure:"key"`
	} `mapstructure:"access"`
	Refresh struct {
		Expire time.Duration `mapstructure:"expire"`
		Key    string        `mapstructure:"key"`
	} `mapstructure:"refresh"`
}
