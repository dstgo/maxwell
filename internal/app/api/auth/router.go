package auth

import (
	"github.com/ginx-contribs/ginx"
)

type Router struct {
	auth *AuthAPI
}

func NewRouter(group *ginx.RouterGroup, auth *AuthAPI) Router {

	authGroup := group.Group("/auth", nil)
	authGroup.POST("/login", nil, auth.Login)
	authGroup.POST("/register", nil, auth.Register)
	authGroup.POST("/reset", nil, auth.ResetPassword)
	authGroup.POST("/refresh", nil, auth.Refresh)

	return Router{auth: auth}
}
