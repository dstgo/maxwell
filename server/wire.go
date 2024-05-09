//go:build wireinject

// The build tag makes sure the stub is not built in the final build
package server

import (
	"github.com/dstgo/maxwell/server/api"
	"github.com/dstgo/maxwell/server/data"
	"github.com/dstgo/maxwell/server/handler"
	"github.com/dstgo/maxwell/server/types"
	"github.com/google/wire"
)

// initialize and setup app environment
func setup(env *types.Env) (api.Router, error) {
	panic(wire.Build(EnvProvider, data.Provider, handler.Provider, api.Provider))
}
