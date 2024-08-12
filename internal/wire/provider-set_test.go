package wire

import (
	"os"
	"testing"

	"github.com/hauntedness/dot/internal/store"
)

func TestProviderSet_Build(t *testing.T) {
	modbase := "github.com/some/mod/"
	var ps ProviderSet
	ps.AddProvider(&store.Provider{
		PvdPkgPath:  modbase + "liu",
		PvdPkgName:  "liu",
		PvdFuncName: "NewLiu",
		PvdName:     "NewLiu",
		PvdKind:     "func",
		PvdError:    0,
		CmpPkgPath:  modbase + "liu",
		CmpPkgName:  "liu",
		CmpTypName:  "Liu",
		CmpKind:     0,
	})
	ps.AddProvider(&store.Provider{
		PvdPkgPath:  modbase + "guan",
		PvdPkgName:  "guan",
		PvdFuncName: "NewGuan",
		PvdName:     "NewGuan",
		PvdKind:     "func",
		PvdError:    0,
		CmpPkgPath:  modbase + "guan",
		CmpPkgName:  "guan",
		CmpTypName:  "Guan",
		CmpKind:     0,
	})
	ps.AddBind(
		&store.ProviderRequirement{
			PvdPkgPath: modbase + "liu",
			PvdPkgName: "liu",
			PvdName:    "NewLiu",
			PvdKind:    "func",
			CmpPkgPath: modbase + "liu",
			CmpPkgName: "liu",
			CmpTypName: "Namer",
			CmpKind:    0,
		},
		&store.Provider{
			PvdPkgPath:  modbase + "guan",
			PvdPkgName:  "guan",
			PvdFuncName: "NewGuan",
			PvdName:     "NewGuan",
			PvdKind:     "func",
			PvdError:    0,
			CmpPkgPath:  modbase + "guan",
			CmpPkgName:  "guan",
			CmpTypName:  "Guan",
			CmpKind:     0,
		},
	)
	bytes, err := os.ReadFile("testdata/provider-set-data.txt")
	if err != nil {
		t.Fatal(err)
	}
	if text := ps.Build(); text != string(bytes) {
		t.Fatalf("ps.String() != string(bytes)")
	}
}
