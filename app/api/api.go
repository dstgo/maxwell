package api

import "github.com/dstgo/maxwell/app/types"

// RegisterRouter
// @title	                        MaxWell HTTP API
// @version		                    v0.0.0Beta
// @description                     This is maxwell swagger generated api documentation, know more information about maxwell on GitHub.
// @contact.name                    dstgo
// @contact.url                     https://github.com/dstgo/maxwell
// @BasePath	                    /api/
// @license.name                    MIT LICENSE
// @license.url                     https://mit-license.org/
// @securityDefinitions.apikey      BearerAuth
// @in                              header
// @name                            Authorization
//
//go:generate swag init --ot yaml --generatedTime --instanceName "maxwell" -g api.go -d ./,../types --output ./docs && swag fmt -g api.go -d ./
func RegisterRouter(env *types.Env) {

}
