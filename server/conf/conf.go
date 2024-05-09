package conf

import (
	"github.com/ginx-contribs/dbx"
	"log/slog"
	"time"
)

type App struct {
	Server Server      `mapstructure:"server"`
	Log    Log         `mapstructure:"log"`
	DB     dbx.Options `mapstructure:"db"`
	Redis  Redis       `mapstructure:"redis"`
	Email  Email       `mapstructure:"email"`
	Jwt    Jwt         `mapstructure:"jwt"`
	Limit  RateLimit   `mapstructure:"limit"`

	Version   string `mapstructure:"-"`
	BuildTime string `mapstructure:"-"`
}

type Server struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
	MultipartMax int64         `mapstructure:"multipartMax"`
	Pprof        bool          `mapstructure:"pprof"`
	TLS          struct {
		Cert string `mapstructure:"cert"`
		Key  string `mapstructure:"key"`
	} `mapstructure:"tls"`
}

type Log struct {
	Filename string     `mapstructure:"filename"`
	Prompt   string     `mapstructure:"-"`
	Level    slog.Level `mapstructure:"level"`
	Format   string     `mapstructure:"format"`
	Source   bool       `mapstructure:"source"`
	Color    bool       `mapstructure:"color"`
}

type Redis struct {
	Address      string        `mapstructure:"address"`
	Password     string        `mapstructure:"password"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
}

type Jwt struct {
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
	Public struct {
		Limit  int           `mapstructure:"limit"`
		Window time.Duration `mapstructure:"window"`
	} `mapstructure:"public"`
}

type Email struct {
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
