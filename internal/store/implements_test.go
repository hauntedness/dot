package store

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/hauntedness/dot/internal/inj/liu"
)

func TestSaveImplement(t *testing.T) {
	type args struct {
		impl *ImplementStmt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				impl: &ImplementStmt{
					CmpPkgPath:   lt.PkgPath(),
					CmpTypName:   "Guan2",
					CmpKind:      3,
					IfacePkgPath: lt.PkgPath(),
					IfaceName:    lt.Name(),
				},
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				impl: &ImplementStmt{
					CmpPkgPath:   lt.PkgPath(),
					CmpTypName:   "Guan3",
					CmpKind:      3,
					IfacePkgPath: lt.PkgPath(),
					IfaceName:    lt.Name(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveImplement(tt.args.impl); (err != nil) != tt.wantErr {
				t.Errorf("SaveImplement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindImplements(t *testing.T) {
	lt := reflect.TypeFor[liu.Namer]()
	type args struct {
		InterfacePkg  string
		InterfaceName string
	}
	tests := []struct {
		name string
		args args
		want func([]ImplementStmt, error) bool
	}{
		{
			name: "",
			args: args{
				InterfacePkg:  lt.PkgPath(),
				InterfaceName: lt.Name(),
			},
			want: func(is []ImplementStmt, err error) bool {
				return len(is) == 2 && err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindImplementsByInterface(tt.args.InterfacePkg, tt.args.InterfaceName)
			if !tt.want(got, err) {
				t.Errorf("FindImplements() = %v, %v", got, err)
			}
		})
	}
}

func TestViewImplements(t *testing.T) {
	rows, err := db.Query("select * from implement_stmts")
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
