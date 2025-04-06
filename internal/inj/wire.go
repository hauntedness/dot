//go:build wireinject

// The build tag makes sure the stub is not built in the final build.

package inj

import (
	"github.com/google/wire"
	"github.com/hauntedness/dot/internal/inj/liu"
)

func initializeBaz() (*liu.Liu, error) {
	wire.Build(LiuSet)
	return nil, nil
}
