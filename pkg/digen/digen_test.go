package digen

import (
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/liu"
	"github.com/hauntedness/dot/internal/store"
	"github.com/hauntedness/dot/internal/types"
)

func TestProviderGen_GenerateFromStruct(t *testing.T) {
	store.Init()
	pg := NewProviderSetGen("")
	structs := types.LoadStructs(reflect.TypeFor[liu.Liu]().PkgPath())
	err := pg.GenerateFromStruct(structs[0])
	if err != nil {
		t.Fatal(err)
	}
}

func TestProviderGen_GenerateAnyStruct(t *testing.T) {
	store.Init()
	pg := NewProviderSetGen("")
	err := pg.GenerateAnyStruct(liu.Liu{})
	if err != nil {
		t.Fatal(err)
	}
}
