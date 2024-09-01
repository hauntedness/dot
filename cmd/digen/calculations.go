package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/hauntedness/dot/internal/store"
	"github.com/hauntedness/dot/internal/types"
	"github.com/hauntedness/dot/internal/wire"
)

type ProviderGen struct {
	label string
}

func (pg *ProviderGen) GenerateProviderSet(pkg string) error {
	loaded := types.Load(pkg)
	providers, err := pg.FindProviderByPkg(loaded.PkgPath)
	if err != nil {
		return err
	}
	for _, provider := range providers {
		ps := &wire.ProviderSet{Name: "Wire" + provider.PvdFuncName + "Set"}
		ps.AddProvider(&provider)
		err := pg.walk(ps, &provider)
		if err != nil {
			return err
		}
		fmt.Println(ps.Build())
	}
	return nil
}

func (pg *ProviderGen) walk(ps *wire.ProviderSet, provider *store.Provider) error {
	//
	requirements, err := pg.FindProviderRequirements(provider)
	if err != nil {
		return err
	}
	for _, requirement := range requirements {
		// 如果是interface需要按名字查找, 因为类型和package并不匹配
		if types.IsInterfaceKind(types.TypeKind(requirement.CmpKind)) {
			provider, err := pg.findInterfaceProvider(&requirement)
			if err != nil {
				return err
			}
			ps.AddProvider(provider)
			ps.AddBind(&requirement, provider)
			err = pg.walk(ps, provider)
			if err != nil {
				return err
			}
		} else {
			providers, err := pg.FindProviderByCmp(requirement.CmpPkgPath, requirement.CmpTypName, requirement.CmpKind)
			if err != nil {
				return err
			}
			if l := len(providers); l == 0 {
				return fmt.Errorf("no providers found for %#v", requirement)
			} else if l > 1 {
				if requirement.CmpPvdName == "" {
					return fmt.Errorf("could not determine providers. req: %v, providers: %v", requirement, providers)
				}
				var p *store.Provider
				var count int
				for _, provider := range providers {
					name := cmp.Or(provider.PvdName, provider.PvdFuncName)
					if requirement.CmpPvdName == name {
						count++
						p = &provider
					}
				}
				if count != 1 {
					return fmt.Errorf("could not determine providers. req: %v, providers: %v", requirement, providers)
				}
				ps.AddProvider(p)
				err = pg.walk(ps, p)
				if err != nil {
					return err
				}
			} else if l == 1 {
				provider := &providers[0]
				if requirement.CmpPvdName != "" {
					name := cmp.Or(provider.PvdName, provider.PvdFuncName)
					if name != requirement.CmpPvdName {
						return fmt.Errorf("could not found provider %s, please check the provide name.", requirement.CmpPvdName)
					}
				}
				ps.AddProvider(provider)
				err := pg.walk(ps, provider)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (pg *ProviderGen) checkLabel(requiredLabels string) bool {
	if requiredLabels == "" {
		return true
	}
	if pg.label == "" {
		return false
	}
	lb := types.Labels{}
	return lb.Append(requiredLabels).Labeled(pg.label)
}

func (pg *ProviderGen) findInterfaceProvider(req *store.ProviderRequirement) (*store.Provider, error) {
	checked := func(providers []store.Provider, err error) (*store.Provider, error) {
		if err != nil {
			return nil, err
		}
		if len(providers) != 1 {
			return nil, fmt.Errorf("could not determine provider for interface kind. req: %v, providers: %v", req, providers)
		}
		return &providers[0], nil
	}
	if req.CmpPvdName == "" {
		//
		implements, err := pg.FindImplementsByInterface(req.CmpPkgPath, req.CmpTypName)
		if err != nil {
			return nil, err
		}
		if len(implements) != 1 {
			return nil, fmt.Errorf("could not determine implementation for interface kind. req: %v, implementations: %v", req, implements)
		}
		// find by implementation
		providers1, err1 := pg.FindProviderByCmp(implements[0].CmpPkgPath, implements[0].CmpTypName, implements[0].CmpKind)
		// find by interface directly
		providers2, err2 := pg.FindProviderByCmp(req.CmpPkgPath, req.CmpTypName, req.CmpKind)
		return checked(append(providers1, providers2...), errors.Join(err1, err2))
	}
	providers, err := pg.FindProviderByName(req.CmpPvdName)
	if err != nil {
		return nil, err
	}
	cp := make([]store.Provider, 0, 1)
	for _, provider := range providers {
		impl := &store.ImplementStmt{
			CmpPkgPath:   provider.CmpPkgPath,
			CmpTypName:   provider.CmpTypName,
			CmpKind:      provider.CmpKind,
			IfacePkgPath: req.CmpPkgPath,
			IfaceName:    req.CmpTypName,
		}
		implements, err := pg.FindImplementsByImpl(impl)
		if err != nil {
			return nil, err
		}
		if len(implements) == 1 {
			cp = append(cp, provider)
		}
	}
	return checked(cp, err)
}

func (pg *ProviderGen) FindProviderByCmp(pkg string, typ string, kind int) ([]store.Provider, error) {
	list, err := store.FindProviderByCmp(pkg, typ, kind)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(list, func(i store.Provider) bool {
		return !pg.checkLabel(i.Labels)
	}), nil
}

func (pg *ProviderGen) FindProviderByPkg(pkg string) ([]store.Provider, error) {
	list, err := store.FindProviderByPkg(pkg)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(list, func(i store.Provider) bool {
		return !pg.checkLabel(i.Labels)
	}), nil
}

func (pg *ProviderGen) FindProviderByName(component string) ([]store.Provider, error) {
	list, err := store.FindProviderByName(component)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(list, func(i store.Provider) bool {
		return !pg.checkLabel(i.Labels)
	}), nil
}

func (pg *ProviderGen) FindProviderRequirements(c *store.Provider) ([]store.ProviderRequirement, error) {
	requirements, err := store.FindProviderRequirements(c)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(requirements, func(i store.ProviderRequirement) bool {
		return !pg.checkLabel(i.Labels)
	}), nil
}

func (pg *ProviderGen) FindImplementsByInterface(InterfacePkg string, InterfaceName string) ([]store.ImplementStmt, error) {
	list, err := store.FindImplementsByInterface(InterfacePkg, InterfaceName)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(list, func(i store.ImplementStmt) bool {
		return !pg.checkLabel(i.Labels)
	}), nil
}

func (pg *ProviderGen) FindImplementsByImpl(impl *store.ImplementStmt) ([]store.ImplementStmt, error) {
	list, err := store.FindImplementsByImpl(impl)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(list, func(i store.ImplementStmt) bool {
		return !pg.checkLabel(i.Labels)
	}), nil
}
