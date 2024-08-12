package store

import (
	"fmt"
	"testing"
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
