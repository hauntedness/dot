package guan

import "github.com/hauntedness/dot/internal/inj/zhang"

//go:generate dot-ioc
var _ int

//go:ioc component
type Guan struct{}

//go:ioc --param name.provider="NewZhang"
func NewGuan(z zhang.Zhang) *Guan {
	return &Guan{}
}

func (g *Guan) Name() string {
	return "Guan, Yu"
}
