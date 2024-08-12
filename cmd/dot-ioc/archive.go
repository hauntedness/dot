package main

import (
	"go/ast"

	"github.com/hauntedness/dot/internal/store"
	"github.com/hauntedness/dot/internal/types"
	"golang.org/x/tools/go/packages"
)

// todo(思考: 标注component真的有用吗? 看起来只标注provider也是可行的. 而且第三方的包也没法标注component. )
// 将components存到表里
//  1. 处理struct类型的component
//  2. 处理interface类型的component
func (c *Container) Component() error {
	for _, st := range c.structs {
		if len(st.Directives()) == 0 {
			continue
		}
		pkg := st.Pkg()
		comp := store.Component{
			CmpPkgPath: pkg.Path(),
			CmpPkgName: pkg.Name(),
			CmpTypName: st.Name(),
			CmpName:    st.ComponentName(),
			CmpKind:    int(types.TypeKindStruct),
		}
		err := store.SaveComponent(&comp)
		if err != nil {
			return err
		}
	}
	for _, it := range c.interfaces {
		if len(it.Directives()) == 0 {
			continue
		}
		pkg := it.Pkg()
		comp := store.Component{
			CmpPkgPath: pkg.Path(),
			CmpPkgName: pkg.Name(),
			CmpTypName: it.Name(),
			CmpName:    it.ComponentName(),
			CmpKind:    int(types.KindInterface),
		}
		err := store.SaveComponent(&comp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (*Container) LoadComponentDirective(_decl *ast.GenDecl, commentMap map[*ast.Ident][]string, comments ast.CommentMap, decl ast.Decl, pkg *packages.Package) {
	if len(_decl.Specs) != 1 {
		return
	}
	switch _decl := _decl.Specs[0].(type) {
	case *ast.TypeSpec:
		commentMap[_decl.Name] = types.Directives(comments.Filter(decl), pkg)
	case *ast.ValueSpec:
		if len(_decl.Names) != 1 {
			return
		}
		commentMap[_decl.Names[0]] = types.Directives(comments.Filter(decl), pkg)
	}
}
