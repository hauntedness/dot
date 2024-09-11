package inj

import (
	"github.com/google/wire"
	"github.com/hauntedness/dot/internal/inj/guan"
	"github.com/hauntedness/dot/internal/inj/liu"
	"github.com/hauntedness/dot/internal/inj/liu2"
	"github.com/hauntedness/dot/internal/inj/zhang"
	"github.com/hauntedness/dot/internal/inj/zhang/yanyan"
)

var LiuSet = wire.NewSet(
	liu.NewLiu,
	guan.NewGuan,
	wire.Bind(new(liu.Namer), new(*guan.Guan)),
	zhang.NewZhang,
	yanyan.NewYanYan,
	liu2.NewGuan2,
)

var _ liu.Namer = (*guan.Guan)(nil)
