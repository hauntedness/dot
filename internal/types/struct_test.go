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
		fmt.Println(
			st.typeName.String(),
			st.typeName.Name(),
			st.typeName.IsAlias(),
			st.named.String(),
			st.ParentScope(),
		)
	}
}
