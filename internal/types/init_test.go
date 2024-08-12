package types

import (
	"reflect"

	"github.com/hauntedness/dot/internal/inj/liu"
)

var (
	liu1 = reflect.TypeFor[liu.Liu]()
	pkg1 = Load(liu1.PkgPath())
)
