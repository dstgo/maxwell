package auth

import (
	"github.com/ginx-contribs/ginx"
)

type Router struct {
	auth *AuthAPI
}

func NewRouter(group *ginx.RouterGroup, auth *AuthAPI) Router {

	// ping
	group.GET("/ping", auth.Ping)

	authGroup := group.Group("/auth")
	authGroup.POST("/login", auth.Login)
	authGroup.POST("/register", auth.Register)
	authGroup.POST("/reset", auth.ResetPassword)
	authGroup.POST("/refresh", auth.Refresh)
	authGroup.POST("/code", auth.VerifyCode)

	return Router{auth: auth}
}
