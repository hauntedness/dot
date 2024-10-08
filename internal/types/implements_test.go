package types

import (
	"go/ast"
	"log/slog"
	"slices"
	"testing"
)

func TestImplements(t *testing.T) {
	pkg1 := Load("github.com/hauntedness/dot/internal/inj/liu2")
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
				case *ast.ValueSpec:
					stmt, ok := NewImplementStmt(_decl, pkg1)
					slog.Info("value spec", "stmt", stmt, "ok", ok)
					commentMap[_decl.Names[0]] = Directives(comments.Filter(decl), pkg1)
				}
			case *ast.FuncDecl:
				commentMap[_decl.Name] = Directives(comments.Filter(decl), pkg1)
				continue
			}
		}
	}
	for id, def := range pkg1.TypesInfo.Defs {
		if id == nil {
			continue
		}
		switch Kind(def) {
		case KindFunc:
			fn, err := NewFunc(def)
			if err != nil {
				continue
			}
			stmtList, ok := NewImplementStmtSlice(fn, pkg1)
			slog.Info("value spec", "stmt", stmtList, "ok", ok)
		}
	}
}

func TestImplementStmt_SetDirectives(t *testing.T) {
	t.Run("test implements argument 'labels' will be processed properly.", func(t *testing.T) {
		impl := &ImplementStmt{}
		impl.SetDirectives([]string{"//go:ioc implements --labels local,local2,local,local,local"})
		equal := slices.Compare(impl.labels.labels, []string{"local", "local2"})
		if equal != 0 {
			t.Fatalf("unexpected results: %v", impl.labels)
		}
		if !impl.labels.Labeled("local") || !impl.labels.Labeled("local2") {
			t.Fatalf("unexpected results: %v", impl.labels)
		}
	})
}
