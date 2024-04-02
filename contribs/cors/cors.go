package cors

import (
	"github.com/spf13/cast"
	"strings"
)

// cors client request header
const (
	Origin                     = "origin"
	AccessControlRequestMethod = "Access-Control-Request-Method"
)

// cors server response header
const (
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowHeaders     = "Access-Control-Allow-AccessAllowHeaders"
	AccessControlAllowMethods     = "Access-Control-Allow-AccessAllowsMethods"
	AccessControlExposeHeaders    = "Access-Control-Expose-AccessAllowHeaders"
	AccessControlMaxAge           = "Access-Control-Max-Age"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

type Cors struct {
	Enabled          bool     `mapstructure:"enable"`
	AllowOrigins     []string `mapstructure:"allowOrigins"`
	AllowMethods     []string `mapstructure:"allowMethods"`
	AllowHeaders     []string `mapstructure:"allowHeaders"`
	ExposeHeaders    []string `mapstructure:"exposedHeaders"`
	MaxAge           int64    `mapstructure:"maxAge"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
}

func (c Cors) AccessAllowMethods() string {
	return strings.Join(c.AllowMethods, ",")
}

func (c Cors) AccessAllowHeaders() string {
	return strings.Join(c.AllowHeaders, ",")
}

func (c Cors) AccessExposedHeaders() string {
	return strings.Join(c.ExposeHeaders, ",")
}

func (c Cors) AccessMaxAge() string {
	return cast.ToString(c.MaxAge)
}

func (c Cors) AccessCredentials() string {
	return cast.ToString(c.AllowCredentials)
}
