package main

import (
	"embed"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"math"

	"golang.org/x/tools/go/packages"
)

// ensure importing embed
var _ embed.FS

//go:embed p/greet.go
var greet string

//go:embed p/greet_test.go
var greet_test string

func main() {
	// Create the AST by parsing greet and greet_test.
	fset := token.NewFileSet()
	files := []*ast.File{
		mustParse(fset, "greet.go", greet),
		mustParse(fset, "greet_test.go", greet_test),
	}

	// Compute package documentation with examples.
	pkg, err := doc.NewFromFiles(fset, files, "example.com/doc/p")
	if err != nil {
		panic(err)
	}
	pkgs, err := packages.Load(&packages.Config{Mode: math.MaxInt}, ".")
	if err != nil {
		panic(err)
	}
	_ = pkgs
	fmt.Printf("package %s - %s", pkg.Name, pkg.Doc)
	printFuncs(pkg)
	printStruct(pkg)
}

func printStruct(pkg *doc.Package) {
	for _, _type := range pkg.Types {
		fmt.Printf("%#v", _type)
	}
}

func printFuncs(pkg *doc.Package) {
	fmt.Printf("func %s - %s", pkg.Funcs[0].Name, pkg.Funcs[0].Doc)
	fmt.Printf(" ⤷ example with suffix %q - %s", pkg.Funcs[0].Examples[0].Suffix, pkg.Funcs[0].Examples[0].Doc)
}

func mustParse(fset *token.FileSet, filename, src string) *ast.File {
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}
