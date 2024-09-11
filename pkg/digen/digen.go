package digen

import (
	"fmt"
	"reflect"

	"github.com/hauntedness/dot/internal/store"
	"github.com/hauntedness/dot/internal/types"
	"github.com/hauntedness/dot/internal/wire"
)

// Struct
func (pg *ProviderGen) GenerateFromStruct(st *types.Struct) error {
	ps := &wire.ProviderSet{Name: "Wire" + st.Name() + "Set"}
	ps.SetStruct(st)
	for i, field := range st.Fields() {
		tag, ok := st.Tag(i).Lookup("wire")
		if ok && tag == "-" {
			continue
		}
		if field.Name() == "_" {
			continue
		}
		pkg, typ, kind := field.TypePkg().Path(), field.TypeName(), int(field.TypeKind())
		providers, err := pg.FindProviderByCmp(pkg, typ, kind)
		if err != nil {
			return err
		}
		if len(providers) != 1 {
			return fmt.Errorf("could not determine providers. (pkg,typ,kind): (%v,%v,%v), providers: %v", pkg, typ, kind, providers)
		}
		p := &providers[0]
		ps.AddProvider(p)
		err = pg.walk(ps, p)
		if err != nil {
			return err
		}
	}
	fmt.Println(ps.Build())
	return nil
}

// Struct
func (pg *ProviderGen) GenerateFromStructPkg(pkg string) error {
	structs := types.LoadStructs(pkg)
	for _, st := range structs {
		err := pg.GenerateFromStruct(st)
		if err != nil {
			return err
		}
	}
	return nil
}

// GenerateAnyStruct print the provider set needed for this struct.
func (pg *ProviderGen) GenerateAnyStruct(value any) error {
	typ := reflect.TypeOf(value)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	structs := types.LoadStructs(typ.PkgPath())
	for _, st := range structs {
		if st.Name() != typ.Name() {
			continue
		}
		err := pg.GenerateFromStruct(st)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateAnyStruct(value any) error {
	store.Init()
	pg := NewProviderSetGen("")
	err := pg.GenerateAnyStruct(value)
	return err
}
