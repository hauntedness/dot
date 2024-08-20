package fst

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"testing"

	"golang.org/x/tools/go/ast/astutil"
)

func TestRewrite(t *testing.T) {
	fset := token.NewFileSet()
	syntax, err := parser.ParseFile(fset, "wire.go", nil, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	astutil.AddNamedImport(fset, syntax, "_", "github.com/hauntedness/dot/internal/inj/liu")
	syntax2 := astutil.Apply(syntax, func(c *astutil.Cursor) bool {
		if id, ok := c.Node().(*ast.Ident); ok {
			if id.Name == "Set" {
				id2 := ast.NewIdent("WireSet2")
				c.Replace(id2)
			}
		}
		return true
	}, nil)
	file, err := os.Create("wire.go")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	printer.Fprint(file, fset, syntax2)
}
