package types

import (
	"go/types"
	"log"
	"math"

	"golang.org/x/tools/go/packages"
)

func Load(path string) *packages.Package {
	pkgs, err := packages.Load(&packages.Config{Mode: math.MaxInt}, path)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("path is not a package")
	}
	return pkgs[0]
}

// IocGen create an init function for type T.
//
// it automatically search the dependencies tree, and fill the dependencies.
func IocGen(typ types.Type) {
	// 搜索本包和依赖的包
	// 搜索构造器
	// 关于参数注入, 分为基本类型, struct, slice, map, interface, 还有不支持的类型比如chan
	// 基本类型和不支持的类型直接放到全局map中, 初始为0值, 用户可以自定义修改
	// struct, slice, map, interface有可能有多种候选者, 根据情况选择注入, 并允许被替换
	//
	// 默认生成error
	// 递归查找
	// 需要详细的错误信息
}
