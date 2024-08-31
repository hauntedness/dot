package guan

import "github.com/hauntedness/dot/internal/inj/zhang"

//go:generate digen
var _ int

//go:ioc component
type Guan struct{}

//go:ioc provider --param z.provider="NewZhang"
func NewGuan(z zhang.Zhang) *Guan {
	return &Guan{}
}

func (g *Guan) Name() string {
	return "Guan, Yu"
}
