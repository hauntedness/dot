package liu2

//go:generate digen

type Guan2 struct{}

// Name implements Namer.
func (g *Guan2) Name() string {
	panic("unimplemented")
}

// digen will treat this as implement statement
//
//go:ioc implements
func (g *Guan2) Implements() Namer {
	return g
}

//go:ioc component
type Namer interface {
	Name() string
}
