package mid

import (
	"github.com/dstgo/maxwell/contribs/cors"
	"github.com/dstgo/maxwell/contribs/ginx/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

// CorsHandler deals with case of cors origins
func CorsHandler(corsHelper *cors.Cors) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if corsHelper == nil || !corsHelper.Enabled {
			ctx.Next()
			return
		}
		origin := ctx.GetHeader(cors.Origin)

		// write response headers
		corsHeader := ctx.Writer.Header()

		// origin == *
		if len(corsHelper.AllowOrigins) == 1 && corsHelper.AllowOrigins[0] == "*" {
			corsHeader.Set(cors.AccessControlAllowOrigin, "*")
		} else if slices.Contains(corsHelper.AllowOrigins, origin) {
			corsHeader.Set(cors.AccessControlAllowOrigin, origin)
		}

		// set attributes
		corsHeader.Set(cors.AccessControlAllowMethods, corsHelper.AccessAllowMethods())
		corsHeader.Set(cors.AccessControlAllowHeaders, corsHelper.AccessAllowHeaders())
		corsHeader.Set(cors.AccessControlExposeHeaders, corsHelper.AccessExposedHeaders())
		corsHeader.Set(cors.AccessControlMaxAge, corsHelper.AccessMaxAge())
		corsHeader.Set(cors.AccessControlAllowCredentials, corsHelper.AccessCredentials())

		// abort if is options method
		if ctx.Request.Method == http.MethodOptions {
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// ResourceNotFoundHandler deals with case of 404
func ResourceNotFoundHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp.New(ctx).Status(http.StatusNotFound).Msg("resource not found").Render()
	}
}

// MethodAllowHandler deals with case of method not allowed
func MethodAllowHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp.New(ctx).Status(http.StatusMethodNotAllowed).Msg("method not allowed").Render()
	}
}
