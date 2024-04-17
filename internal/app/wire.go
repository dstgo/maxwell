//go:build wireinject

// The build tag makes sure the stub is not built in the final build
package app

import (
	"github.com/dstgo/maxwell/internal/app/api"
	"github.com/dstgo/maxwell/internal/app/data"
	"github.com/dstgo/maxwell/internal/app/handler"
	"github.com/dstgo/maxwell/internal/app/types"
	"github.com/google/wire"
)

// initialize and setup app environment
func setup(env *types.Env) (api.Router, error) {
	panic(wire.Build(EnvProvider, data.Provider, handler.Provider, api.Provider))
}
