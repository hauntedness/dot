package types

import (
	"log/slog"
	"testing"
)

func TestLoad(t *testing.T) {
	if pkg1.TypesInfo == nil {
		slog.Warn("pkg.TypesInfo == nil, perhaps packages.NeedTypesInfo is not specified in load config.")
	}
	for _, def := range pkg1.TypesInfo.Defs {
		fn, err := NewFunc(def)
		if err == nil {
			slog.Info("func", "v", fn)
			continue
		}
		st, err := NewStruct(def)
		if err == nil {
			slog.Info("named", "value", st)
			continue
		}
	}
}
