package types

import (
	"go/ast"
	"path/filepath"
	"testing"

	"github.com/hauntedness/dot/internal/inj/liu"
)

func TestComment(t *testing.T) {
	commentMap := map[*ast.Ident][]string{}
	for _, syntax := range pkg1.Syntax {
		comments := ast.NewCommentMap(pkg1.Fset, syntax, syntax.Comments)
		for _, decl := range syntax.Decls {
			switch _decl := decl.(type) {
			case *ast.GenDecl:
				if len(_decl.Specs) != 1 {
					continue
				}
				switch _decl := _decl.Specs[0].(type) {
				case *ast.TypeSpec:
					commentMap[_decl.Name] = Directives(comments.Filter(decl), pkg1)
					continue
				case *ast.ValueSpec:
					if len(_decl.Names) != 1 {
						continue
					}
					commentMap[_decl.Names[0]] = Directives(comments.Filter(decl), pkg1)
					continue
				}
			case *ast.FuncDecl:
				commentMap[_decl.Name] = Directives(comments.Filter(decl), pkg1)
				continue
			default:
				// omit other case
			}
		}
	}
	for id, def := range pkg1.TypesInfo.Defs {
		if id == nil {
			continue
		}
		switch Kind(def) {
		case KindStruct:
			st, err := NewStruct(def)
			if err != nil {
				t.Fatal(err)
			}
			st.directives = commentMap[id]
			pos := pkg1.Fset.Position(id.Pos())
			if fileName := liu.FileName(); filepath.ToSlash(pos.Filename) != filepath.ToSlash(fileName) {
				t.Fatalf("pos.Filename(=%s) != fileName(=%s) , pos.Line = %d", pos.Filename, fileName, pos.Line)
			}
			if st.Name() == "Liu" && len(st.Directives()) == 0 {
				t.Fatalf("Error: len(st.Directives()) == 0")
			}
		case KindInterface:
			it, err := NewInterface(def)
			if err != nil {
				t.Fatal(err)
			}
			it.directives = commentMap[id]
			pos := pkg1.Fset.Position(id.Pos())
			if fileName := liu.FileName(); filepath.ToSlash(pos.Filename) != filepath.ToSlash(fileName) {
				t.Fatalf("pos.Filename(=%s) != fileName(=%s) , pos.Line = %d", pos.Filename, fileName, pos.Line)
			}
			if it.Name() == "Namer" && len(it.Directives()) == 0 {
				t.Fatalf("Error: len(st.Directives()) == 0")
			}
		case KindFunc:
			fn, err := NewFunc(def)
			if err != nil {
				t.Fatal(err)
			}
			fn.SetDirectives(commentMap[id])
			pos := pkg1.Fset.Position(id.Pos())
			if fileName := liu.FileName(); filepath.ToSlash(pos.Filename) != filepath.ToSlash(fileName) {
				t.Fatalf("pos.Filename(=%s) != fileName(=%s) , pos.Line = %d", pos.Filename, fileName, pos.Line)
			}
			if fn.Name() == "NewLiu2" && len(fn.paramSetttings) == 0 {
				t.Fatalf("Error: len(fn.Directives()) == 0")
			}
		case KindVar:
			var1, err := NewVar(def)
			if err != nil {
				t.Fatal(err)
			}
			var1.directives = commentMap[id]
			pos := pkg1.Fset.Position(id.Pos())
			if fileName := liu.FileName(); filepath.ToSlash(pos.Filename) != filepath.ToSlash(fileName) {
				t.Fatalf("pos.Filename(=%s) != fileName(=%s) , pos.Line = %d", pos.Filename, fileName, pos.Line)
			}
			if var1.Name() == "Liu3" {
				if len(var1.Directives()) == 0 {
					t.Fatalf("Error: len(fn.Directives()) == 0")
				}
			}
		}
	}
}
