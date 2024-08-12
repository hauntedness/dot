package inj

import (
	"github.com/google/wire"
	"github.com/hauntedness/dot/internal/inj/guan"
	"github.com/hauntedness/dot/internal/inj/liu"
)

var Liu = wire.NewSet(
	liu.NewLiu,
	guan.NewGuan,
	wire.Bind(new(liu.Namer), new(guan.Guan)),
)
