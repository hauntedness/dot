package types

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"iter"
	"reflect"
)

type Struct struct {
	structType    *types.Struct
	typeName      *types.TypeName
	named         *types.Named
	labels        Labels
	componentName string
	isComponent   bool
	autowire      bool
}

func LoadStructs(pkg string) []*Struct {
	loaded := Load(pkg)
	commentMap := map[*ast.Ident][]string{}
	for _, syntax := range loaded.Syntax {
		comments := ast.NewCommentMap(loaded.Fset, syntax, syntax.Comments)
		for _, decl := range syntax.Decls {
			switch _decl := decl.(type) {
			case *ast.GenDecl:
				if len(_decl.Specs) != 1 {
					continue
				}
				switch _decl := _decl.Specs[0].(type) {
				case *ast.TypeSpec:
					commentMap[_decl.Name] = Directives(comments.Filter(decl), loaded)
					continue
				}
			}
		}
	}
	res := make([]*Struct, 0, 1)
	for id, def := range loaded.TypesInfo.Defs {
		if id == nil {
			continue
		}
		if c := commentMap[id]; len(c) != 0 {
			st, err := NewStruct(def)
			if err != nil {
				continue
			}
			st.SetDirectives(commentMap[id])
			res = append(res, st)
		}
	}
	return res
}

func NewStruct(obj types.Object) (*Struct, error) {
	t, ok := obj.(*types.TypeName)
	if !ok {
		return nil, fmt.Errorf("obj is not a struct: %v", obj)
	}
	n, ok := t.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("obj is not a named type: %v", obj)
	}
	s, ok := n.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("obj is not a struct: %v", obj)
	}
	return &Struct{typeName: t, named: n, structType: s}, nil
}

func (s *Struct) Field(i int) *types.Var {
	return s.structType.Field(i)
}

func (s *Struct) NumFields() int {
	return s.structType.NumFields()
}

func (s *Struct) Fields() iter.Seq2[int, *Var] {
	return func(yield func(int, *Var) bool) {
		for i := range s.NumFields() {
			v, err := NewVar(s.Field(i))
			if err != nil {
				panic(err)
			}
			if !yield(i, v) {
				return
			}
		}
	}
}

func (s *Struct) Tag(i int) reflect.StructTag {
	return reflect.StructTag(s.structType.Tag(i))
}

func (s *Struct) Pkg() *types.Package {
	return s.typeName.Pkg()
}

func (s *Struct) String() string {
	return s.typeName.String()
}

func (s *Struct) Name() string {
	return s.typeName.Name()
}

func (s *Struct) ParentScope() *types.Scope {
	return s.typeName.Parent()
}

func (s *Struct) Type() types.Type {
	return s.typeName.Type()
}

func (s *Struct) Exported() bool {
	return s.typeName.Exported()
}

func (s *Struct) Method(i int) *types.Func {
	return s.named.Method(i)
}

func (s *Struct) NumMethods() int {
	return s.named.NumMethods()
}

func (s *Struct) Origin() *types.Named {
	return s.named.Origin()
}

func (s *Struct) TypeArgs() *types.TypeList {
	return s.named.TypeArgs()
}

func (s *Struct) TypeParams() *types.TypeParamList {
	return s.named.TypeParams()
}

func (s *Struct) SetDirectives(directives []string) {
	dir := Directive{cmd: "component", docs: directives, fs: flag.NewFlagSet("component", flag.PanicOnError)}
	dir.fs.String("name", "", "name of the struct")
	dir.fs.Bool("wire", false, "whether generate wire definition for this function.")
	dir.fs.String("labels", "", "label that take this func into account.")
	err := dir.Parse(func(g *flag.Flag) {
		if g.Name == "name" {
			s.componentName = g.Value.String()
		} else if g.Name == "wire" && g.Value.String() == "true" {
			s.autowire = true
		} else if g.Name == "labels" {
			s.labels.Append(g.Value.String())
		}
	})
	if err != nil {
		panic(err)
	} else {
		s.isComponent = true
	}
}

// ComponentName return name of the component, specified by --name.
// default to camelCased TypeName
// eg: go:ioc component --name book
func (s *Struct) ComponentName() string {
	return s.componentName
}

func (s *Struct) IsComponent() bool {
	return s.isComponent
}

func (s *Struct) Autowire() bool {
	return s.autowire
}
