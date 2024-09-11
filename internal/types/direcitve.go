package types

import (
	"flag"
	"fmt"
	"slices"
	"strings"
)

type Directive struct {
	cmd  string
	docs []string
	fs   *flag.FlagSet
}

func (d *Directive) Parse(fn func(*flag.Flag)) error {
	for _, doc := range d.docs {
		directive := strings.Replace(doc, "// go:ioc", "//go:ioc", 1)
		prefix := "//go:ioc"
		if d.cmd != "" {
			prefix = prefix + " " + d.cmd
		}
		args, ok := strings.CutPrefix(directive, prefix)
		if !ok {
			return fmt.Errorf("directive:%s invalid prefix.", doc)
		}
		args = strings.TrimSpace(args)
		if len(args) == 0 {
			continue
		}
		if !strings.HasPrefix(args, "-") {
			return fmt.Errorf("directive:%s invalid prefix.", doc)
		}
		cmd := strings.Split(args, " ")
		err := d.fs.Parse(cmd)
		if err != nil {
			return err
		}
		d.fs.VisitAll(fn)
	}
	return nil
}

func (d *Directive) Cmd() string {
	return d.cmd
}

func AppendLabels(labelstr string, labels []string) []string {
	ls := strings.Split(labelstr, ",")
	for i := range ls {
		ls[i] = strings.TrimSpace(ls[i])
	}
	return slices.Compact(slices.Sorted(slices.Values(append(labels, ls...))))
}
