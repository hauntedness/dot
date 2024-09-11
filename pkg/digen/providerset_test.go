package digen

import (
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/liu"
	"github.com/hauntedness/dot/internal/store"
)

func TestProviderSetGen_GenerateFromFuncPkg(t *testing.T) {
	store.Init()
	path := reflect.TypeFor[liu.Liu]().PkgPath()
	pg := NewProviderSetGen("dev")
	err := pg.GenerateFromFuncPkg(path)
	if err != nil {
		t.Fatal(err)
	}
}
