package types

import "testing"

func TestFunc_SetDirectives1(t *testing.T) {
	fn := &Func{}
	fn.SetDirectives([]string{`//go:ioc --param name.ident="liu" --name high_recommended`})
	if fn.pvdName != "high_recommended" {
		t.Fatalf(`fn.additionIdentity != "high_recommended"`)
	} else if fn.paramSetttings["name"]["ident"] != `"liu"` {
		t.Fatalf(`fn.paramSetttings["name"]["ident"] != "liu"`)
	}
}

func TestFunc_SetDirectives2(t *testing.T) {
	fn := &Func{}
	fn.SetDirectives([]string{`//go:ioc --param name.provider=NewLiu2 --name high_recommended`})
	if fn.pvdName != "high_recommended" {
		t.Fatalf(`fn.additionIdentity != "high_recommended"`)
	} else if fn.paramSetttings["name"]["provider"] != "NewLiu2" {
		t.Fatalf(`fn.paramSetttings["name"]["provider"] != NewLiu2`)
	}
}
