package types

import (
	"fmt"
	"go/types"
)

type Struct struct {
	*types.Struct
	_TypeName  *types.TypeName
	_Named     *types.Named
	directives []string
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
	return &Struct{_TypeName: t, _Named: n, Struct: s}, nil
}

func (s *Struct) Pkg() *types.Package {
	return s._TypeName.Pkg()
}

func (s *Struct) String() string {
	return s._TypeName.String()
}

func (s *Struct) Name() string {
	return s._TypeName.Name()
}

func (s *Struct) ParentScope() *types.Scope {
	return s._TypeName.Parent()
}

func (s *Struct) Type() types.Type {
	return s._TypeName.Type()
}

func (s *Struct) Exported() bool {
	return s._TypeName.Exported()
}

func (s *Struct) Method(i int) *types.Func {
	return s._Named.Method(i)
}

func (s *Struct) NumMethods() int {
	return s._Named.NumMethods()
}

func (s *Struct) Origin() *types.Named {
	return s._Named.Origin()
}

func (s *Struct) TypeArgs() *types.TypeList {
	return s._Named.TypeArgs()
}

func (s *Struct) TypeParams() *types.TypeParamList {
	return s._Named.TypeParams()
}

func (s *Struct) Directives() []string {
	return s.directives
}

func (s *Struct) SetDirectives(directives []string) {
	s.directives = directives
}

// ComponentName return name of the component, specified by --name.
// default to camelCased TypeName
// eg: go:ioc component --name book
func (s *Struct) ComponentName() string {
	prefix := []string{"// go:ioc component ", "//go:ioc component "}
	return SearchDirectives(componentFlags, s.directives, prefix, "name", CamelCase(s.Name()))
}
