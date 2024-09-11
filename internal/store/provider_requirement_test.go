package store

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/guan"
)

func TestFindProviderRequirements(t *testing.T) {
	guanType := reflect.TypeFor[guan.Guan]()
	providerRequirements, err := FindProviderRequirements(
		&Provider{
			PvdPkgPath: guanType.PkgPath(),
			PvdPkgName: "guan",
			PvdName:    "NewGuan",
		})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", providerRequirements)
}

func TestFindAllProviderRequirements(t *testing.T) {
	// guanType := reflect.TypeFor[guan.Guan]()
	providerRequirements, err := Select[ProviderRequirement]("select * from provider_requirements")
	if err != nil {
		t.Fatal(err)
	}
	for _, req := range providerRequirements {
		fmt.Printf("%#v\n", req)
	}
}
