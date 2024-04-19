package app

import (
	"github.com/dstgo/maxwell/internal/app/types"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/go-playground/validator/v10"
)

// handler to process when params validating failed
func validatePramsHandler(ctx *gin.Context, val any, err error) {
	if _, ok := err.(validator.ValidationErrors); ok {
		ctx.Error(err)
		resp.Fail(ctx).Error(types.ErrBadParams).JSON()
	} else {
		resp.InternalError(ctx).Error(err).JSON()
	}
}
