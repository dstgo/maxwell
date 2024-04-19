package conf

import (
	"github.com/ginx-contribs/dbx"
	"log/slog"
	"time"
)

type AppConf struct {
	Server ServerConf    `mapstructure:"server"`
	Log    LogConf       `mapstructure:"log"`
	DB     dbx.Options   `mapstructure:"db"`
	Redis  RedisConf     `mapstructure:"redis"`
	Email  EmailConf     `mapstructure:"email"`
	Jwt    JwtConf       `mapstructure:"jwt"`
	Limit  RateLimitConf `mapstructure:"limit"`

	Version   string `mapstructure:"-"`
	BuildTime string `mapstructure:"-"`
}

type ServerConf struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
	MultipartMax int64         `mapstructure:"multipartMax"`
	Pprof        bool          `mapstructure:"pprof"`
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

type RateLimitConf struct {
	IpLimit struct {
		Limit  int           `mapstructure:"limit"`
		Window time.Duration `mapstructure:"window"`
	} `mapstructure:"ipLimit"`
}

type EmailConf struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	MQ       struct {
		Topic     string   `mapstructure:"topic"`
		MaxLen    int64    `mapstructure:"maxLen"`
		BatchSize int64    `mapstructure:"batchSize"`
		Group     string   `mapstructure:"group"`
		Consumers []string `mapstructure:"consumers"`
	} `mapstructure:"mq"`
	Code struct {
		TTL      time.Duration `mapstructure:"ttl"`
		RetryTTL time.Duration `mapstructure:"retry"`
	} `mapstructure:"code"`
}
