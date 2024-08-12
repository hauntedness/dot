package main

import (
	"cmp"
	"flag"
	"fmt"
	"go/ast"
	"log"
	"log/slog"
	"os"

	"github.com/hauntedness/dot/internal/store"
	"github.com/hauntedness/dot/internal/types"
	"github.com/hauntedness/dot/internal/wire"
	"golang.org/x/tools/go/packages"
)

func main() {
	store.Init()
	path := flag.String("package", ".", "the package or dir to be scanned")
	wireset := flag.Bool("wireset", false, "generate wire set")
	flag.Parse()
	if *path == "" {
		log.Panic("Must specify the package flags or GOPACKAGE env when running without go generate.")
	}
	if ok1, ok2 := *wireset, os.Getenv("wireset"); ok1 || ok2 == "true" {
		err := GenerateProviderSet(*path)
		if err != nil {
			log.Panic(err)
		}
	} else {
		err := Generate(*path)
		if err != nil {
			log.Panic(err)
		}
		dir, err := os.Getwd()
		if err != nil {
			log.Panic(err)
		}
		slog.Info(dir, "msg", "dot-ioc definition loading is finished.")
	}
}

func Generate(path string) error {
	pkg := types.Load(path)
	c := &Container{}
	err := c.LoadProvider(pkg)
	if err != nil {
		return err
	}
	for _, fn := range c.funcs {
		err = CheckFuncProvider(fn)
		if err != nil {
			continue
		}
		err = SaveFuncProvider(fn)
		if err != nil {
			return err
		}
		err = SaveFuncProviderRequirement(fn)
		if err != nil {
			return err
		}
	}
	return nil
}

type Container struct {
	structs    []*types.Struct
	funcs      []*types.Func
	interfaces []*types.Interface
	vars       []*types.Var
}

func (c *Container) LoadProvider(pkg *packages.Package) error {
	commentMap := map[*ast.Ident][]string{}
	for _, syntax := range pkg.Syntax {
		comments := ast.NewCommentMap(pkg.Fset, syntax, syntax.Comments)
		for _, decl := range syntax.Decls {
			switch _decl := decl.(type) {
			case *ast.FuncDecl:
				commentMap[_decl.Name] = types.Directives(comments.Filter(decl), pkg)
				continue
			default:
				// omit other case
			}
		}
	}
	// skip processing when no directives
	if len(commentMap) == 0 {
		return nil
	}
	for id, def := range pkg.TypesInfo.Defs {
		if id == nil {
			continue
		}
		switch types.Kind(def) {
		case types.KindFunc:
			dd := commentMap[id]
			if len(dd) > 0 {
				fn, err := types.NewFunc(def)
				if err != nil {
					return err
				}
				fn.SetDirectives(dd)
				c.funcs = append(c.funcs, fn)
			}
		case types.KindVar:
			var1, err := types.NewVar(def)
			if err != nil {
				return err
			}
			var1.SetDirectives(commentMap[id])
			c.vars = append(c.vars, var1)
		}
	}
	return nil
}

func CheckFuncProvider(fn *types.Func) error {
	if !fn.IsValid() {
		return fmt.Errorf("fn doesn't have directives")
	}
	if l := fn.Results().Len(); l > 2 || l < 1 {
		return fmt.Errorf("Function Provider should be in form of fn(...) T or fn(...) (T, error)")
	} else if l == 2 {
		typ, err := fn.ResultType(1)
		if err != nil {
			return err
		}
		if !types.IsError(typ) {
			return fmt.Errorf("Function Provider should be in form of fn(...) T or fn(...) (T, error)")
		}
	}
	result, err := fn.Result(0)
	if err != nil {
		return err
	}
	res, err := types.NewVar(result)
	if err != nil {
		return err
	}
	kind := res.TypeKind()
	pass := false
	if kind&types.TypeKindStruct == types.TypeKindStruct {
		pass = true
	} else if kind&types.TypeKindStructPointer == types.TypeKindStructPointer {
		pass = true
	} else if kind&types.TypeKindInterface == types.TypeKindInterface {
		pass = true
	}
	if !pass {
		return fmt.Errorf("fn doesn't have correct return type. %v", fn)
	}
	return nil
}

func SaveFuncProvider(fn *types.Func) error {
	result, err := fn.Result(0)
	if err != nil {
		return fmt.Errorf("provider should at least have one result. %w", err)
	}
	res, err := types.NewVar(result)
	if err != nil {
		return err
	}
	pkg := fn.Pkg()
	pvd := store.Provider{
		PvdPkgPath:  pkg.Path(),
		PvdPkgName:  pkg.Name(),
		PvdFuncName: fn.Name(),
		PvdName:     fn.PvdName(),
		PvdKind:     fn.Kind(),
		CmpPkgPath:  res.TypePkg().Path(),
		CmpPkgName:  res.TypePkg().Name(),
		CmpTypName:  res.TypeName(),
		//
		CmpKind: int(res.TypeKind()),
	}
	if fn.ReturnError() {
		pvd.PvdError = 1
	}
	return store.SaveProvider(&pvd)
}

func SaveFuncProviderRequirement(fn *types.Func) error {
	if !fn.IsValid() {
		return fmt.Errorf("fn doesn't have directives")
	}
	for i := range fn.Params().Len() {
		param, err := fn.Param(i)
		if err != nil {
			return err
		}
		v, err := types.NewVar(param)
		if err != nil {
			return err
		}
		// 如果是基本类型, 需要在go:ioc指令下指令字面量值
		// 如果是接口类型, 需要指定具体的实现类, (可选指定provider)
		kind, pass := v.TypeKind(), false
		if kind&types.TypeKindStruct == types.TypeKindStruct {
			pass = true
		} else if kind&types.TypeKindStructPointer == types.TypeKindStructPointer {
			pass = true
		} else if kind&types.TypeKindInterface == types.TypeKindInterface {
			// 对于接口类型,
			pass = true
		}
		if !pass {
			return fmt.Errorf("fn doesn't have correct param type. %v", fn)
		}
		pkg, pkg1 := fn.Pkg(), v.TypePkg()
		pvd := store.ProviderRequirement{
			PvdPkgPath:  pkg.Path(),
			PvdPkgName:  pkg.Name(),
			PvdFuncName: fn.Name(),
			PvdName:     fn.Name(),
			PvdKind:     fn.Kind(),
			CmpPkgPath:  pkg1.Path(),
			CmpPkgName:  pkg1.Name(),
			CmpTypName:  v.TypeName(),
			CmpName:     v.Name(),
			CmpPvdName:  fn.ParamPvd(v.Name()),
			//
			CmpKind: int(v.TypeKind()),
		}
		if err := store.SaveProviderRequirement(&pvd); err != nil {
			return err
		}
	}
	return nil
}

func GenerateProviderSet(pkg string) error {
	pkg0 := types.Load(pkg)
	providers, err := store.FindProviderByPkg(pkg0.PkgPath)
	if err != nil {
		return err
	}
	for _, provider := range providers {
		ps := &wire.ProviderSet{Name: "Wire" + provider.PvdFuncName + "Set"}
		ps.AddProvider(&provider)
		err := visit(ps, &provider)
		if err != nil {
			return err
		}
		fmt.Println(ps.Build())
	}
	return nil
}

func visit(ps *wire.ProviderSet, provider *store.Provider) error {
	//
	requirements, err := store.FindProviderRequirements(provider)
	if err != nil {
		return err
	}
	for _, requirement := range requirements {
		// 如果是interface需要按名字查找, 因为类型和package并不匹配
		if types.IsInterfaceKind(types.TypeKind(requirement.CmpKind)) {
			if requirement.CmpPvdName == "" {
				return fmt.Errorf("provide name must be specified for interface kind. %#v", requirement)
			}
			// TODO(j) 是否要继续验证类型?
			providers, err := store.FindProviderByName(requirement.CmpPvdName)
			if err != nil {
				return err
			} else if len(providers) != 1 {
				return fmt.Errorf("provider name should be unique when requiring a interface value. %#v", requirement)
			}
			provider := &providers[0]
			ps.AddProvider(provider)
			ps.AddBind(&requirement, provider)
			err = visit(ps, provider)
			if err != nil {
				return err
			}
		} else {
			providers, err := store.FindProviderByCmp(requirement.CmpPkgPath, requirement.CmpTypName, requirement.CmpKind)
			if err != nil {
				return err
			}
			if l := len(providers); l == 0 {
				return fmt.Errorf("no providers found for %#v", requirement)
			} else if l > 1 {
				if requirement.CmpPvdName == "" {
					return fmt.Errorf("found multiple providers, you must specify the provider name. %#v", requirement)
				}
				var p *store.Provider
				var count int
				for _, provider := range providers {
					name := cmp.Or(provider.PvdName, provider.PvdFuncName)
					if requirement.CmpPvdName == name {
						count++
						p = &provider
					}
				}
				if count != 1 {
					return fmt.Errorf("could not find provider for %#v", requirement)
				}
				ps.AddProvider(p)
				err = visit(ps, p)
				if err != nil {
					return err
				}
			} else if l == 1 {
				provider := &providers[0]
				if requirement.CmpPvdName != "" {
					name := cmp.Or(provider.PvdName, provider.PvdFuncName)
					if name != requirement.CmpPvdName {
						return fmt.Errorf("could not found provider %s, please check the provide name.", requirement.CmpPvdName)
					}
				}
				ps.AddProvider(provider)
				err := visit(ps, provider)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
