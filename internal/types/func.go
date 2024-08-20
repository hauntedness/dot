package types

import (
	"flag"
	"fmt"
	"go/types"
	"strings"
)

type Func struct {
	fn             *types.Func
	sig            *types.Signature
	paramSetttings map[string]map[string]string
	pvdName        string
	ready          bool
}

func NewFunc(obj types.Object) (*Func, error) {
	if obj == nil {
		return nil, fmt.Errorf("obj is nil")
	}
	if fn, ok := obj.(*types.Func); ok {
		return &Func{fn: fn, sig: fn.Type().(*types.Signature)}, nil
	} else {
		return nil, fmt.Errorf("obj is not a Func: %v", obj)
	}
}

func (f *Func) Kind() string {
	return "func"
}

func (f *Func) Params() *types.Tuple {
	return f.sig.Params()
}

func (f *Func) Param(i int) (*types.Var, error) {
	res := f.sig.Params()
	if res == nil || res.Len() < i+1 {
		return nil, fmt.Errorf("Func params length error for index=%d", i)
	}
	return res.At(i), nil
}

func (f *Func) ParamType(i int) (types.Type, error) {
	v, err := f.Param(i)
	if err != nil {
		return nil, err
	}
	return v.Type(), nil
}

func (f *Func) Recv() *types.Var {
	return f.sig.Recv()
}

func (f *Func) RecvTypeParams() *types.TypeParamList {
	return f.sig.RecvTypeParams()
}

func (f *Func) Results() *types.Tuple {
	return f.sig.Results()
}

func (f *Func) Result(i int) (*types.Var, error) {
	res := f.sig.Results()
	if res == nil || res.Len() < i+1 {
		return nil, fmt.Errorf("Func results length error for index=%d", i)
	}
	return res.At(i), nil
}

func (f *Func) ResultType(i int) (types.Type, error) {
	v, err := f.Result(i)
	if err != nil {
		return nil, err
	}
	return v.Type(), nil
}

func (f *Func) ReturnError() bool {
	results := f.Results()
	if results.Len() != 2 {
		return false
	}
	if typ := results.At(1).Type(); IsError(typ) {
		return false
	}
	return false
}

func (f *Func) SignatureString() string {
	return f.sig.String()
}

func (f *Func) TypeParams() *types.TypeParamList {
	return f.sig.TypeParams()
}

func (f *Func) Underlying() types.Type {
	return f.sig.Underlying()
}

func (f *Func) Variadic() bool {
	return f.sig.Variadic()
}

func (f *Func) Exported() bool {
	return f.fn.Exported()
}

func (f *Func) FullName() string {
	return f.fn.FullName()
}

func (f *Func) Id() string {
	return f.fn.Id()
}

func (f *Func) Name() string {
	return f.fn.Name()
}

func (f *Func) Parent() *types.Scope {
	return f.fn.Parent()
}

func (f *Func) Pkg() *types.Package {
	return f.fn.Pkg()
}

func (f *Func) String() string {
	return f.fn.String()
}

func (f *Func) SetDirectives(directives []string) {
	if len(f.paramSetttings) == 0 {
		f.paramSetttings = map[string]map[string]string{}
	}
	for _, d := range directives {
		directive := strings.Replace(d, "// go:ioc", "//go:ioc", 1)
		args, ok := strings.CutPrefix(directive, "//go:ioc")
		if !ok {
			panic(fmt.Errorf("directive:%s has no valid prefix.", d))
		}
		args = strings.TrimSpace(args)
		if len(args) == 0 {
			continue
		}
		cmd := strings.Split(args, " ")
		err := providerFlags.Parse(cmd)
		if err != nil {
			panic(err)
		}
		providerFlags.VisitAll(func(g *flag.Flag) {
			// process name
			if g.Name == "name" {
				f.pvdName = g.Value.String()
				return
			} else if g.Name == "param" {
				value := g.Value.String()
				if value == "" {
					return
				}
				// process command eg: --param name.provider="NewLiu2"
				key, cmd, ok := strings.Cut(value, ".")
				if !ok {
					panic(fmt.Errorf("unknown param settings: %s", value))
				}
				key = strings.TrimSpace(key)
				subkey, c, ok := strings.Cut(cmd, "=")
				if !ok {
					panic(fmt.Errorf("unknown param settings: %s", value))
				}
				subkey, c = strings.TrimSpace(subkey), strings.TrimSpace(c)
				if len(f.paramSetttings[key]) == 0 {
					f.paramSetttings[key] = map[string]string{}
				}
				f.paramSetttings[key][subkey] = c
			}
		})
	}
	f.ready = true
}

func (f *Func) PvdName() string {
	if f.pvdName != "" {
		return f.pvdName
	}
	return f.Name()
}

func (f *Func) ParamPvd(param string) string {
	if settings, ok := f.paramSetttings[param]; ok {
		return strings.Trim(settings["provider"], "\"")
	}
	return ""
}

func (f *Func) IsValid() bool {
	return f.ready
}
