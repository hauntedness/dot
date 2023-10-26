package dot

import (
	"testing"

	"github.com/hauntedness/dot/internal/doc/play"
)

func TestLoadPackage(t *testing.T) {
	pkg, err := LoadPackage(play.Book{}, nil)
	if err != nil {
		t.Error(err)
		return
	}
	ns, err := pkg.LookupStruct(play.Book{})
	if err != nil {
		t.Error(err)
		return
	}

	if num := ns.NumFields(); num != 3 {
		t.Errorf("NumFields %d", num)
		return
	}
	if name := ns.FieldName(1); name != "Words" {
		t.Errorf("FieldName %s", name)
		return
	}

	if typeString := ns.FieldTypeString(2); typeString != "*json2.Marshaler" {
		t.Errorf("FieldTypeString %s", typeString)
		return
	}
	if tag := ns.FieldTag(1); tag != "" {
		t.Errorf("FieldTag %s", tag)
		return
	}
}
