package auth

import (
	authandler "github.com/dstgo/maxwell/internal/app/handler/auth"
	"github.com/dstgo/maxwell/internal/app/types/auth"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

func NewAuthAPI(token *authandler.TokenHandler, auth *authandler.AuthHandler) *AuthAPI {
	return &AuthAPI{token: token, auth: auth}
}

type AuthAPI struct {
	token *authandler.TokenHandler
	auth  *authandler.AuthHandler
}

// Ping
// @Summary      Ping
// @Description  test service if is available
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.Response
// @Router       /ping [GET]
func (a *AuthAPI) Ping(ctx *gin.Context) {
	resp.Ok(ctx).Msg("pong").JSON()
}

// Login
// @Summary      Login
// @Description  use login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginOption  query  auth.LoginOption  true "login params"
// @Success      200  {object}  types.Response{data=auth.TokenResult}
// @Router       /auth/login [POST]
func (a *AuthAPI) Login(ctx *gin.Context) {
	var loginOpt auth.LoginOption
	if err := ctx.ShouldBindJSON(&loginOpt); err != nil {
		resp.Fail(ctx).Error(err).JSON()
		return
	}
	// login by username and password
	tokenPair, err := a.auth.LoginByPassword(ctx, loginOpt)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Data(auth.TokenResult{
			AccessToken:  tokenPair.AccessToken.TokenString,
			RefreshToken: tokenPair.RefreshToken.TokenString,
		}).JSON()
	}
}

// Register
// @Summary      Register
// @Description  register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        RegisterOption  body  auth.RegisterOption  true "register params"
// @Success      200  {object}  types.Response
// @Router       /auth/register [POST]
func (a *AuthAPI) Register(ctx *gin.Context) {
	var registerOpt auth.RegisterOption
	if err := ctx.ShouldBindJSON(&registerOpt); err != nil {
		resp.Fail(ctx).Error(err).JSON()
		return
	}

	_, err := a.auth.RegisterNewUser(ctx, registerOpt)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).JSON()
	}
}

// ResetPassword
// @Summary      ResetPassword
// @Description  reset user password without login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        ResetPasswordOption   body  auth.ResetPasswordOption  true "reset params"
// @Success      200  {object}  types.Response
// @Router       /auth/reset [POST]
func (a *AuthAPI) ResetPassword(ctx *gin.Context) {
	var restOpt auth.ResetPasswordOption
	if err := ctx.ShouldBindJSON(&restOpt); err != nil {
		resp.Fail(ctx).Error(err).JSON()
		return
	}

	if err := a.auth.ResetPassword(ctx, restOpt); err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).JSON()
	}
}

// Refresh
// @Summary      Refresh
// @Description  ask for refresh token lifetime
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        RefreshTokenOption  body  auth.RefreshTokenOption  true "refresh params"
// @Success      200  {object}  types.Response{data=auth.TokenResult}
// @Router       /auth/refresh [POST]
func (a *AuthAPI) Refresh(ctx *gin.Context) {
	var refreshOpt auth.RefreshTokenOption
	if err := ctx.ShouldBindJSON(&refreshOpt); err != nil {
		resp.Fail(ctx).Error(err).JSON()
		return
	}

	// ask for refresh token
	tokenPair, err := a.token.Refresh(ctx, refreshOpt.AccessToken, refreshOpt.RefreshToken)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Data(auth.TokenResult{
			AccessToken:  tokenPair.AccessToken.TokenString,
			RefreshToken: tokenPair.RefreshToken.TokenString,
		}).JSON()
	}
}
