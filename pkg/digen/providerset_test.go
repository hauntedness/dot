package digen

import (
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/liu"
)

func TestProviderSetGen_GenerateFromFuncPkg(t *testing.T) {
	path := reflect.TypeFor[liu.Liu]().PkgPath()
	pg := NewProviderSetGen("dev")
	err := pg.GenerateFromFuncPkg(path)
	if err != nil {
		t.Fatal(err)
	}
}
