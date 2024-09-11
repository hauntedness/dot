// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inj

import (
	"github.com/hauntedness/dot/internal/inj/guan"
	"github.com/hauntedness/dot/internal/inj/liu"
	"github.com/hauntedness/dot/internal/inj/zhang"
	"github.com/hauntedness/dot/internal/inj/zhang/yanyan"
)

// Injectors from wire.go:

func initializeBaz() (*liu.Liu, error) {
	yanYan := yanyan.NewYanYan()
	zhangZhang := zhang.NewZhang(yanYan)
	guanGuan := guan.NewGuan(zhangZhang)
	liuLiu := liu.NewLiu(guanGuan)
	return liuLiu, nil
}
