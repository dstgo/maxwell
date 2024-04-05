package conf

import (
	"github.com/ginx-contribs/dbx"
	"log/slog"
	"time"
)

type AppConf struct {
	Server ServerConf  `mapstructure:"server"`
	Log    LogConf     `mapstructure:"log"`
	DB     dbx.Options `mapstructure:"db"`
	Redis  RedisConf   `mapstructure:"redis"`

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

type LogConf struct {
	Filename string     `mapstructure:"filename"`
	Prompt   string     `mapstructure:"-"`
	Level    slog.Level `mapstructure:"level"`
	Format   string     `mapstructure:"format"`
	Source   bool       `mapstructure:"source"`
	Color    bool       `mapstructure:"color"`
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

type RateLimit struct {
	IpLimit struct {
		Limit  int           `mapstructure:"limit"`
		Window time.Duration `mapstructure:"window"`
	} `mapstructure:"ipLimit"`

	Token struct {
	}
}
