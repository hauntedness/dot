package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"log"
	"log/slog"
	"os"

	"github.com/hauntedness/dot/internal/store"
	"github.com/hauntedness/dot/internal/types"
	"golang.org/x/tools/go/packages"
)

func main() {
	store.Init()
	path := flag.String("pkg", ".", "the package or dir to be scanned")
	pset := flag.Bool("gen", false, "generate wire provider set")
	label := flag.String("label", "", "generate wire provider set")
	flag.Parse()
	if *path == "" {
		log.Panic("Must specify the package flags or GOPACKAGE env when running without go generate.")
	}
	if ok1, ok2 := *pset, os.Getenv("gen_provider_set"); ok1 || ok2 == "true" {
		pg := &ProviderGen{label: *label}
		err := pg.GenerateProviderSet(*path)
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
		slog.Info(dir, "msg", "digen definition loading is finished.")
	}
}

func Generate(path string) error {
	pkg := types.Load(path)
	c := &Container{}
	err := c.LoadDefinitions(pkg)
	if err != nil {
		return err
	}
	cleanup := func() error {
		return errors.Join(
			store.DeleteComponentByPkg(pkg.PkgPath),
			store.DeleteImplementByPkg(pkg.PkgPath),
			store.DeleteProviderByPkg(pkg.PkgPath),
			store.DeleteProviderRequirementsByPkg(pkg.PkgPath),
		)
	}
	if err := cleanup(); err != nil {
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
	for _, impl := range c.implements {
		err := SaveImplements(impl)
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
	implements []*types.ImplementStmt
}

func (c *Container) LoadDefinitions(pkg *packages.Package) error {
	directives := map[*ast.Ident][]string{}
	for _, syntax := range pkg.Syntax {
		comments := ast.NewCommentMap(pkg.Fset, syntax, syntax.Comments)
		for _, decl := range syntax.Decls {
			switch _decl := decl.(type) {
			case *ast.GenDecl:
				if len(_decl.Specs) != 1 {
					continue
				}
				switch _decl := _decl.Specs[0].(type) {
				case *ast.ValueSpec:
					stmt, ok := types.NewImplementStmt(_decl, pkg)
					if ok {
						c.implements = append(c.implements, stmt)
						ds := types.Directives(comments.Filter(decl), pkg)
						stmt.SetDirectives(ds)
						directives[_decl.Names[0]] = ds
					}
				}
			case *ast.FuncDecl:
				directives[_decl.Name] = types.Directives(comments.Filter(decl), pkg)
				continue
			default:
				// omit other case
			}
		}
	}
	// skip processing when no directives
	if len(directives) == 0 {
		return nil
	}
	for id, def := range pkg.TypesInfo.Defs {
		if id == nil {
			continue
		}
		switch types.Kind(def) {
		case types.KindFunc:
			dd := directives[id]
			if len(dd) > 0 {
				fn, err := types.NewFunc(def)
				if err != nil {
					return err
				}
				if list, ok := types.NewImplementStmtSlice(fn, pkg); ok {
					for _, impl := range list {
						impl.SetDirectives(dd)
						c.implements = append(c.implements, impl)
						slog.Debug("implements", "impl", impl)
					}
				} else {
					fn.SetDirectives(dd)
					c.funcs = append(c.funcs, fn)
				}
			}
		case types.KindVar:
			var1, err := types.NewVar(def)
			if err != nil {
				return err
			}
			var1.SetDirectives(directives[id])
			c.vars = append(c.vars, var1)
		}
	}
	return nil
}

func CheckImplementStatementFunc(fn *types.Func) error {
	return nil
}

func CheckFuncProvider(fn *types.Func) error {
	if !fn.IsValid() {
		return fmt.Errorf("fn: %v missing directives.", fn)
	}
	if l := fn.Results().Len(); l > 2 || l < 1 {
		return fmt.Errorf("fn: %v incorrect number of results.", fn)
	} else if l == 2 {
		typ, ok := fn.ResultType(1)

		if !ok || !types.IsError(typ) {
			return fmt.Errorf("fn: %v 2nd result type is not error", fn)
		}
	}
	result, ok := fn.Result(0)
	if !ok {
		return fmt.Errorf("fn: %v should at least have one result.", fn)
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
		return fmt.Errorf("1st result is not supported type. %v", fn)
	}
	return nil
}

func SaveImplements(impl *types.ImplementStmt) error {
	stmt := store.ImplementStmt{
		IfacePkgPath: impl.IfacePkg().Path(),
		IfaceName:    impl.IfaceName(),
		CmpPkgPath:   impl.ImplPkg().Path(),
		CmpTypName:   impl.ImplName(),
		//
		CmpKind: 0,
		Labels:  impl.Labels(),
	}
	typeKind := types.TypeKindOf(impl.ImplType())
	if impl.IsPointerImpl() {
		typeKind = typeKind | types.TypeKindPointer
	}
	stmt.CmpKind = int(typeKind)
	return store.SaveImplement(&stmt)
}

func SaveFuncProvider(fn *types.Func) error {
	result, ok := fn.Result(0)
	if !ok {
		panic("result should be checked already")
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
		Labels:  fn.Labels(),
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
			Labels:  fn.Labels(),
		}
		if err := store.SaveProviderRequirement(&pvd); err != nil {
			return err
		}
	}
	return nil
}
