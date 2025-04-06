package liu2

//go:generate digen -cmd=scan

type Guan2 struct{}

// Name implements Namer.
func (g *Guan2) Name() string {
	panic("unimplemented")
}

func NewGuan2() *Guan2 {
	return nil
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

// PkgAAAProviderSet (A0,A1,A2)
// PkgBBBProviderSet (B0,B1,B2)

type (
	Open  struct{}
	Close struct{}
	Put   struct{}
)

//go:ioc component --name box
type Box struct {
	NoName string `wire:"-"`
	open   Open
	close  Close
	put    Put
}
