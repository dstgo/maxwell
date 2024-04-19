package handler

import (
	"github.com/dstgo/maxwell/internal/app/handler/auth"
	emailhandler "github.com/dstgo/maxwell/internal/app/handler/email"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	// auth handlers
	auth.NewAuthHandler,
	auth.NewTokenHandler,
	auth.NewVerifyCodeHandler,

	// email handlers
	emailhandler.NewEmailHandler,
)
