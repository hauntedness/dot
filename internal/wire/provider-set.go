package wire

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/hauntedness/dot/internal/store"
	"github.com/hauntedness/dot/internal/types"
)

type bindDef struct {
	pvd *store.Provider
	dep *store.ProviderRequirement
}

type ProviderSet struct {
	Name      string //
	builder   strings.Builder
	structDef *types.Struct
	providers []*store.Provider
	binds     []bindDef
	imports   map[string]string // map[pkg_name]pkg_pkg_path
}

func (ps *ProviderSet) AddProvider(pvd *store.Provider) {
	ps.resolvePkg(pvd.CmpPkgPath, pvd.CmpPkgName)
	ps.providers = append(ps.providers, pvd)
}

func (ps *ProviderSet) SetStruct(st *types.Struct) {
	ps.resolvePkg("github.com/google/wire", "wire")
	ps.structDef = st
}

func (ps *ProviderSet) AddBind(dep *store.ProviderRequirement, pvd *store.Provider) {
	ps.resolvePkg("github.com/google/wire", "wire")
	ps.resolvePkg(dep.CmpPkgPath, dep.CmpPkgName)
	ps.resolvePkg(pvd.CmpPkgPath, pvd.CmpPkgName)
	ps.binds = append(ps.binds, bindDef{pvd: pvd, dep: dep})
}

func (ps *ProviderSet) Build() string {
	if ps.imports == nil {
		ps.imports = map[string]string{}
	}
	for _, name := range slices.Sorted(maps.Keys(ps.imports)) {
		pkg := ps.imports[name]
		fmt.Fprintf(&ps.builder, "import %s %s\n", name, strconv.Quote(pkg))
	}
	// write empty line for pretty look
	fmt.Fprintf(&ps.builder, "\n\n")
	ps.writeDeclarationStart()
	if ps.structDef != nil {
		ps.writeStruct(ps.structDef)
	}
	ps.compactProviders()
	for _, pvd := range ps.providers {
		ps.writeFuncProvider(pvd)
	}
	for _, bind := range ps.binds {
		ps.writeInterfaceBind(bind.dep, bind.pvd)
	}
	ps.writeDeclarationEnd()
	return (&ps.builder).String()
}

func (ps *ProviderSet) compactProviders() {
	type tp = *store.Provider
	list := ps.providers
	slices.SortFunc(list, func(a, b tp) int { return a.Compare(b) })
	list = slices.CompactFunc(list, func(a, b tp) bool {
		return a.PvdPkgPath == b.PvdPkgPath && a.PvdOriName == b.PvdOriName
	})
	ps.providers = list
}

func (ps *ProviderSet) writeDeclarationStart() {
	(&ps.builder).WriteString("var " + ps.Name + " = wire.NewSet(\n")
}

func (ps *ProviderSet) writeDeclarationEnd() {
	(&ps.builder).WriteString(")\n")
}

func (ps *ProviderSet) writeFuncProvider(pvd *store.Provider) {
	pkg_name := ps.resolvePkg(pvd.PvdPkgPath, pvd.PvdPkgName)
	//
	(&ps.builder).WriteString("\t" + pkg_name + "." + pvd.PvdName + ",\n")
}

func (ps *ProviderSet) writeStruct(st *types.Struct) {
	fmt.Fprintf(&ps.builder, "\twire.Struct(new(%s), \"*\"),\n", st.Name())
}

// WriteBind is used to write wire.Bind method for interface components
func (ps *ProviderSet) writeInterfaceBind(dep *store.ProviderRequirement, pvd *store.Provider) {
	//
	interface_pkg_name := ps.resolvePkg(dep.CmpPkgPath, dep.CmpPkgName)
	component_pkg_name := ps.resolvePkg(pvd.CmpPkgPath, pvd.CmpPkgName)
	// eg "wire.Bind(new(liu.Namer), new(guan.Guan)),"
	star := ""
	if types.IsStructPointerKind(types.TypeKind(pvd.CmpKind)) {
		star = "*"
	}
	fmt.Fprintf((&ps.builder), "\twire.Bind(new(%s.%s), new(%s%s.%s)),\n", interface_pkg_name, dep.CmpTypName, star, component_pkg_name, pvd.CmpTypName)
}

// resolvePkg load package path to import map and also handle conflicts on import name
func (ps *ProviderSet) resolvePkg(pkg_path, pkg_name string) string {
	if ps.imports == nil {
		ps.imports = map[string]string{}
	}
	if path := ps.imports[pkg_name]; path == pkg_path {
		return pkg_name
	} else if path == "" {
		ps.imports[pkg_name] = pkg_path
		return pkg_name
	}
	for name, pkg := range ps.imports {
		if pkg == pkg_path {
			return name
		}
	}
	for i := range 100 {
		name := pkg_name + strconv.Itoa(i)
		exists := ps.imports[name]
		if exists == "" {
			ps.imports[name] = pkg_path
			return name
		}
	}
	panic(fmt.Errorf("could not resolve package name. %s", pkg_path))
}
