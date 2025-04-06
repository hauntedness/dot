package liu

import (
	"runtime"

	"github.com/hauntedness/dot/internal/inj/guan"
)

//go:generate digen -cmd=scan
//go:generate digen -cmd=wire --struct
var _ int

// Liu
//
// this directive is only for document purpose, it takes no effects
//
//go:ioc component --wire
type Liu struct {
	Guan *guan.Guan
}

// NewLiu
//
//go:ioc provider
func NewLiu(guan *guan.Guan) *Liu {
	return &Liu{}
}

// NewLiu3
//
//go:ioc provider --labels dev
func NewLiu3(namer Namer) *Liu {
	return nil
}

func FileName() string {
	pc, file, line, ok := runtime.Caller(0)
	_, _ = pc, line
	if !ok {
		panic("can not get file info")
	}
	return file
}

func (l *Liu) Name() string {
	return "Liu, Bei"
}

// Namer is who has a name
//
// this directive is only for document purpose, it takes no effects
//
//go:ioc component
type Namer interface {
	Name() string
}

//go:ioc implements
var _ Namer = (*guan.Guan)(nil)
