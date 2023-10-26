package dot

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/types"
	"math"
	"reflect"
	"strconv"

	"golang.org/x/tools/go/packages"
)

// LoadPackage loads and returns the Go package named by the given patterns.
//
// Config specifies loading options; nil behaves the same as &packages.Config{Mode: math.MaxInt}.
func LoadPackage(patterns string, config *packages.Config) (*Package, error) {
	if config == nil {
		config = &packages.Config{Mode: math.MaxInt}
	}
	pkgs, err := packages.Load(config, patterns)
	if err != nil {
		return nil, fmt.Errorf("LoadPackage: %w", err)
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("LoadPackage: %d packages found", len(pkgs))
	}
	pkg := pkgs[0]
	docs, err := doc.NewFromFiles(pkg.Fset, pkg.Syntax, pkg.PkgPath)
	if err != nil {
		return nil, fmt.Errorf("LoadPackage: %w", err)
	}
	return &Package{Package: pkg, docs: docs}, nil
}

// Package describes a loaded Go package and the comments.
type Package struct {
	*packages.Package
	docs *doc.Package
}

// LookupStruct retrieves the struct in package by the type of v.
// If the struct is present in the package the value is returned.
// Otherwise the returned value will be nil.
func (p *Package) LookupStruct(v any) (*NamedStruct, error) {
	typ := reflect.TypeOf(v)
	var pkg, typeName string = typ.PkgPath(), typ.Name()
	return p.LookupStructById(pkg, typeName)
}

// LookupStructById works similarly to LookupStruct(v any), but use the literal package path and struct name
func (p *Package) LookupStructById(pkg, structName string) (*NamedStruct, error) {
	result := &NamedStruct{}
	for ident, def := range p.Package.TypesInfo.Defs {
		if def == nil {
			continue
		} else if _, ok := def.(*types.TypeName); !ok {
			continue
		}
		if named, ok := def.Type().(*types.Named); ok {
			if named.Obj().Pkg().Path() == pkg && named.Obj().Name() == structName {
				if _, ok := named.Underlying().(*types.Struct); !ok {
					return nil, fmt.Errorf("LookupStructById: %s is not a struct", structName)
				}
				result.Type = named
				for _, f := range p.Package.Syntax {
					if f.Pos() < ident.NamePos && ident.NamePos < f.End() {
						result.Filename = p.Package.Fset.File(ident.NamePos).Name()
						result.File = f
					}
				}
				for _, t := range p.docs.Types {
					if t.Name == structName {
						result.Doc = *t
					}
				}
				break
			}
		}
	}
	if result.Type == nil {
		return nil, fmt.Errorf("LookupStructById: struct %s not found", structName)
	}
	return result, nil
}

// NamedStruct
type NamedStruct struct {
	// typ represent loaded go types
	Type     *types.Named
	Doc      doc.Type
	Filename string
	File     *ast.File
}

// Underlying return underlying *types.Struct
func (ns *NamedStruct) Underlying() *types.Struct {
	return ns.Type.Origin().Underlying().(*types.Struct)
}

func (ns *NamedStruct) NumFields() int {
	return ns.Underlying().NumFields()
}

func (ns *NamedStruct) Field(i int) *types.Var {
	return ns.Underlying().Field(i)
}

// FieldName return the name of field i
func (ns *NamedStruct) FieldName(i int) string {
	return ns.Underlying().Field(i).Name()
}

// FieldType return types.Type of field i
func (ns *NamedStruct) FieldType(i int) types.Type {
	return ns.Underlying().Field(i).Type()
}

func (ns *NamedStruct) FieldTag(i int) string {
	return ns.Underlying().Tag(i)
}

// FieldTypeString return the string repretension of field i
// the return value would consider if the package is aliased
func (ns *NamedStruct) FieldTypeString(i int) string {
	typ := ns.FieldType(i)
	return types.TypeString(typ, func(p *types.Package) string {
		if alias := ns.Imports()[p.Path()]; alias != "" {
			return alias
		}
		return p.Name()
	})
}

// Imports return the all of the import package path and corresponding alias
func (ns *NamedStruct) Imports() map[string]string {
	if ns.File == nil {
		return nil
	}
	imports := make(map[string]string)
	for _, is := range ns.File.Imports {
		alias := ""
		if is.Name != nil {
			alias = is.Name.Name
		}
		path, err := strconv.Unquote(is.Path.Value)
		if err != nil {
			path = is.Path.Value
		}
		imports[path] = alias
	}
	return imports
}
