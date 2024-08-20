package liu

import (
	"net/http"
	"runtime"

	"github.com/hauntedness/dot/internal/inj/guan"
	"github.com/hauntedness/dot/internal/inj/zhang"
	"github.com/hauntedness/dot/internal/inj/zhang/yanyan"
)

//go:generate digen

// Liu
//
// this directive is only for document purpose, it takes no effects
// go:ioc component
type Liu struct{}

// NewLiu
//
//go:ioc
func NewLiu(guan *guan.Guan) *Liu {
	return &Liu{}
}

// NewLiu2
//
// turn off go:ioc --param name.ident="liu" --name high_recommended
func NewLiu2(name string, guan *guan.Guan) *http.Request {
	return nil
}

// NewLiu3
//
//go:ioc --param namer.provider=NewGuan
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
// go:ioc component
type Namer interface {
	Name() string
}

// Liu3
//
//go:ioc --name=liu3
var Liu3 *Liu = NewLiu(guan.NewGuan(zhang.NewZhang(yanyan.NewYanYan())))
