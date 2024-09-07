package store

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/guan"
)

func TestFindProvider1(t *testing.T) {
	providers, err := FindProviderByPkg(lt.PkgPath())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range providers {
		fmt.Printf("%#v\n", p)
	}
}

func TestFindProvider2(t *testing.T) {
	guan := reflect.TypeFor[guan.Guan]()
	providers, err := FindProviderByCmp(guan.PkgPath(), guan.Name(), 3)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range providers {
		fmt.Printf("%#v\n", p)
	}
}
