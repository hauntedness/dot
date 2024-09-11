package main

import (
	"flag"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/guan"
	"github.com/hauntedness/dot/internal/inj/liu"
	"github.com/hauntedness/dot/internal/inj/zhang/yanyan"
	"github.com/hauntedness/dot/internal/store"
)

func TestMain(m *testing.M) {
	flag.Parse()
	store.Init()
	os.Exit(m.Run())
}

func Test_main(t *testing.T) {
	path := reflect.TypeFor[liu.Liu]().PkgPath()
	err := Scan(path)
	if err != nil {
		log.Panic(err)
	}
}

func Test_main2(t *testing.T) {
	path := reflect.TypeFor[guan.Guan]().PkgPath()
	err := Scan(path)
	if err != nil {
		log.Panic(err)
	}
}

func Test_main3(t *testing.T) {
	path := reflect.TypeFor[yanyan.YanYan]().PkgPath()
	err := Scan(path)
	if err != nil {
		log.Panic(err)
	}
}
