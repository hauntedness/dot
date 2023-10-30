package dot

import (
	"log"
	"log/slog"
	"reflect"
	"strings"
	"testing"

	"github.com/hauntedness/dot/internal/doc/play"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

func TestLoadPackage(t *testing.T) {
	pkg, err := LoadPackage(book, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if pkg == nil {
		t.Errorf("pkg is nil")
		return
	}
	ns, err := pkg.LookupStruct(book)
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
	if comments := ns.FieldComments(2); !reflect.DeepEqual(comments, []string{}) {
		t.Errorf("FieldComments %s", comments)
		return
	}
}

var (
	book = play.Book[string]{}
	typ  = reflect.TypeOf(book)
)

func TestGenGo(t *testing.T) {
	arguments := args.Default()
	arguments.InputDirs = []string{typ.PkgPath()}
	arguments.GoHeaderFilePath = "testdata/header.txt"
	// Override defaults.
	arguments.OutputFileBaseName = "deepcopy_generated"
	arguments.OutputBase = "./internal"
	nameSystems := namer.NameSystems{
		"public": &namer.NameStrategy{
			Join: func(pre string, in []string, post string) string {
				return strings.Join(in, "_")
			},
			PrependPackageNames: 1,
		},
		"raw": namer.NewRawNamer("", nil),
	}
	execute := func(ctx *generator.Context, _ *args.GeneratorArgs) (packages generator.Packages) {
		name := types.Name{
			Package: typ.PkgPath(),
			Name:    "Book[T any]",
			Path:    "",
		}
		ut := ctx.Universe.Type(name)
		if ut.Name.Name != name.Name {
			t.Errorf("actual is not name.Name")
			return
		}
		if len(ut.CommentLines[0]) == 0 || ut.CommentLines[0] != "Book is book" {
			t.Errorf("comment is not read")
			return
		}
		if len(ut.Members) == 0 {
			t.Errorf("members were not read")
			return
		}
		return
	}
	// Run it.
	if err := arguments.Execute(nameSystems, "public", execute); err != nil {
		log.Fatalf("Error: %v", err)
	}
	slog.Info("Completed successfully.")
}
