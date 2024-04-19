package auth

import (
	authandler "github.com/dstgo/maxwell/internal/app/handler/auth"
	"github.com/dstgo/maxwell/internal/app/types/auth"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

func NewAuthAPI(token *authandler.TokenHandler, auth *authandler.AuthHandler, verifycode *authandler.VerifyCodeHandler) *AuthAPI {
	return &AuthAPI{token: token, auth: auth, verifycode: verifycode}
}

type AuthAPI struct {
	token      *authandler.TokenHandler
	auth       *authandler.AuthHandler
	verifycode *authandler.VerifyCodeHandler
}

// Ping
// @Summary      Ping
// @Description  test server if is available
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.Response{data=string}
// @Router       /ping [GET]
func (a *AuthAPI) Ping(ctx *gin.Context) {
	resp.Ok(ctx).Msg("pong").JSON()
}

// Login
// @Summary      Login
// @Description  login with password, and returns jwt token pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginOption  body  auth.LoginOption  true "login params"
// @Success      200  {object}  types.Response{data=auth.TokenResult}
// @Router       /auth/login [POST]
func (a *AuthAPI) Login(ctx *gin.Context) {
	var loginOpt auth.LoginOption
	if err := ginx.ShouldValidateJSON(ctx, &loginOpt); err != nil {
		return
	}

	// login by username and password
	tokenPair, err := a.auth.LoginWithPassword(ctx, loginOpt)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("login ok").Data(auth.TokenResult{
			AccessToken:  tokenPair.AccessToken.TokenString,
			RefreshToken: tokenPair.RefreshToken.TokenString,
		}).JSON()
	}
}

// Register
// @Summary      Register
// @Description  register a new user with verification code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        RegisterOption  body  auth.RegisterOption  true "register params"
// @Success      200  {object}  types.Response
// @Router       /auth/register [POST]
func (a *AuthAPI) Register(ctx *gin.Context) {
	var registerOpt auth.RegisterOption
	if err := ginx.ShouldValidateJSON(ctx, &registerOpt); err != nil {
		return
	}

	_, err := a.auth.RegisterNewUser(ctx, registerOpt)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("register ok").JSON()
	}
}

// ResetPassword
// @Summary      ResetPassword
// @Description  reset user password with verification code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        ResetPasswordOption   body  auth.ResetPasswordOption  true "reset params"
// @Success      200  {object}  types.Response
// @Router       /auth/reset [POST]
func (a *AuthAPI) ResetPassword(ctx *gin.Context) {
	var restOpt auth.ResetPasswordOption
	if err := ginx.ShouldValidateJSON(ctx, &restOpt); err != nil {
		return
	}

	if err := a.auth.ResetPassword(ctx, restOpt); err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("reset ok").JSON()
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
	if err := ginx.ShouldValidateJSON(ctx, &refreshOpt); err != nil {
		return
	}

	// ask for refresh token
	tokenPair, err := a.token.Refresh(ctx, refreshOpt.AccessToken, refreshOpt.RefreshToken)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("refresh ok").Data(auth.TokenResult{
			AccessToken:  tokenPair.AccessToken.TokenString,
			RefreshToken: tokenPair.RefreshToken.TokenString,
		}).JSON()
	}
}

// VerifyCode
// @Summary      VerifyCode
// @Description  send verification code mail to specified email address
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        VerifyCodeOption   body   auth.VerifyCodeOption  true  "VerifyCodeOption"
// @Success      200  {object}  types.Response
// @Router       /auth/code [POST]
func (a *AuthAPI) VerifyCode(ctx *gin.Context) {
	var verifyOpt auth.VerifyCodeOption
	if err := ginx.ShouldValidateJSON(ctx, &verifyOpt); err != nil {
		return
	}

	// check usage
	if err := auth.CheckValidUsage(verifyOpt.Usage); err != nil {
		resp.Fail(ctx).Error(auth.ErrVerifyCodeUsageUnsupported).JSON()
		return
	}

	err := a.verifycode.SendVerifyCodeEmail(ctx, verifyOpt.To, verifyOpt.Usage)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Msg("mail has been sent").JSON()
	}
}
