package types

import (
	"flag"
	"fmt"
	"strings"
)

type Directive struct {
	cmd string
	ds  []string
	fs  *flag.FlagSet
}

func (d Directive) Parse(fn func(*flag.Flag)) error {
	for _, doc := range d.ds {
		directive := strings.Replace(doc, "// go:ioc", "//go:ioc", 1)
		prefix := "//go:ioc"
		if d.cmd != "" {
			prefix = prefix + " " + d.cmd
		}
		args, ok := strings.CutPrefix(directive, prefix)
		if !ok {
			return fmt.Errorf("directive:%s has no valid prefix.", doc)
		}
		args = strings.TrimSpace(args)
		if len(args) == 0 {
			continue
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
