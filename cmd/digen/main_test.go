package main

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/liu"
	"github.com/hauntedness/dot/internal/inj/liu2"
	"github.com/hauntedness/dot/internal/store"
)

func TestMain(m *testing.M) {
	store.Init()
	os.Exit(m.Run())
}

func Test_main(t *testing.T) {
	path := reflect.TypeFor[liu.Liu]().PkgPath()
	err := Generate(path)
	if err != nil {
		log.Panic(err)
	}
}

func Test_main2(t *testing.T) {
	path := reflect.TypeFor[liu2.Guan2]().PkgPath()
	err := Generate(path)
	if err != nil {
		log.Panic(err)
	}
}

func TestGenerateProviderSet(t *testing.T) {
	// path := reflect.TypeFor[liu.Liu]().PkgPath()
	err := GenerateProviderSet("D:\\temp\\dot\\internal\\inj\\liu")
	if err != nil {
		t.Fatal(err)
	}
}
