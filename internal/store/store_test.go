package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/liu"
)

var lt = reflect.TypeFor[liu.Liu]()

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}

func TestSaveProvider(t *testing.T) {
	err := SaveProvider(&Provider{
		PvdPkgPath: lt.PkgPath(),
		PvdPkgName: "liu",
		PvdName:    "NewLiu",
		CmpPkgPath: lt.PkgPath(),
		CmpPkgName: "liu",
		CmpTypName: "Liu",
		CmpKind:    3,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveProviderRequirement(t *testing.T) {
	err := SaveProviderRequirement(&ProviderRequirement{
		PvdPkgPath: lt.PkgPath(),
		PvdPkgName: "liu",
		PvdName:    "NewLiu",
		CmpPkgPath: lt.PkgPath(),
		CmpPkgName: "guan",
		CmpTypName: "Guan2",
		CmpKind:    3,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestView(t *testing.T) {
	rows, err := db.Query("select * from provider_requirements")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	results, err := MapRows(rows)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bytes))
}

func MapRows(rows *sql.Rows) ([]map[string]any, error) {
	results := []map[string]any{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		mp := map[string]any{}
		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}
		dest := make([]any, len(columns))
		for i, ct := range columnTypes {
			dest[i] = reflect.New(ct.ScanType()).Interface()
			mp[columns[i]] = dest[i]
		}
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		results = append(results, mp)
	}
	return results, nil
}
