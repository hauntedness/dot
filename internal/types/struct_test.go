package types

import (
	"fmt"
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
		st.SetDirectives([]string{"//go:ioc component --wire"})
		fmt.Println(
			st.typeName.String(),
			st.typeName.Name(),
			st.typeName.IsAlias(),
			st.named.String(),
			st.ParentScope(),
			st.autowire,
		)
		if st.Name() == "Liu" {
			fmt.Printf("autowire:%v\n", st.autowire)
		}
	}
}

func TestStruct_Doc(t *testing.T) {
	st := &Struct{}
	st.SetDirectives([]string{"//go:ioc component --wire"})
	if !st.autowire {
		t.Fatalf("autowire:%v\n", st.autowire)
	}
}
