package types

import (
	"fmt"
	"go/types"
	"testing"
)

func TestStruct_String(t *testing.T) {
	for id, def := range pkg1.TypesInfo.Defs {
		if id == nil {
			continue
		}
		st, err := NewStruct(def)
		if err != nil {
			continue
		}
		fmt.Println(
			st._TypeName.String(),
			st._TypeName.Name(),
			st._TypeName.IsAlias(),
			st._Named.String(),
			st.ParentScope(),
		)
	}
}

func TestStruct_ComponentName(t *testing.T) {
	st := Struct{}
	directives := [][]string{
		{"//go:ioc component --name book"},
		{"//go:ioc component --kind var"},
	}
	st._TypeName = types.NewTypeName(1, nil, "book", nil)
	for _, directives := range directives {
		st.directives = directives
		name := st.ComponentName()
		if name != "book" {
			t.Fatalf(`name != "book"`)
		}
	}
}
