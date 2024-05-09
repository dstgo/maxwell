package auth

import (
	"github.com/dstgo/maxwell/server/types/route"
	"github.com/ginx-contribs/ginx"
)

type Router struct {
	auth *AuthAPI
}

func NewRouter(group *ginx.RouterGroup, auth *AuthAPI) Router {

	// auth
	authGroup := group.MGroup("/auth", ginx.M{route.Public})
	authGroup.POST("/login", auth.Login)
	authGroup.POST("/register", auth.Register)
	authGroup.POST("/reset", auth.ResetPassword)
	authGroup.POST("/refresh", auth.Refresh)
	authGroup.POST("/code", auth.VerifyCode)

	return Router{auth: auth}
}
