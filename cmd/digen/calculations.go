package main

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/hauntedness/dot/internal/store"
	"github.com/hauntedness/dot/internal/types"
	"github.com/hauntedness/dot/internal/wire"
)

type ProviderGen struct {
	label string
}

func (pg *ProviderGen) GenerateProviderSet(pkg string) error {
	pkg0 := types.Load(pkg)
	providers, err := store.FindProviderByPkg(pkg0.PkgPath)
	if err != nil {
		return err
	}
	for _, provider := range providers {
		if ok := pg.checkLabel(provider.Labels); !ok {
			continue
		}
		ps := &wire.ProviderSet{Name: "Wire" + provider.PvdFuncName + "Set"}
		ps.AddProvider(&provider)
		err := pg.visit(ps, &provider)
		if err != nil {
			return err
		}
		fmt.Println(ps.Build())
	}
	return nil
}

func (pg *ProviderGen) visit(ps *wire.ProviderSet, provider *store.Provider) error {
	//
	requirements, err := store.FindProviderRequirements(provider)
	if err != nil {
		return err
	}
	for _, requirement := range requirements {
		if ok := pg.checkLabel(requirement.Labels); !ok {
			continue
		}
		// 如果是interface需要按名字查找, 因为类型和package并不匹配
		if types.IsInterfaceKind(types.TypeKind(requirement.CmpKind)) {
			provider, err := findInterfaceProvider(&requirement)
			if err != nil {
				return err
			}
			ps.AddProvider(provider)
			ps.AddBind(&requirement, provider)
			err = pg.visit(ps, provider)
			if err != nil {
				return err
			}
		} else {
			providers, err := store.FindProviderByCmp(requirement.CmpPkgPath, requirement.CmpTypName, requirement.CmpKind)
			if err != nil {
				return err
			}
			if l := len(providers); l == 0 {
				return fmt.Errorf("no providers found for %#v", requirement)
			} else if l > 1 {
				if requirement.CmpPvdName == "" {
					return fmt.Errorf("found multiple providers, you must specify the provider name. %#v", requirement)
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
				if count == 0 {
					return fmt.Errorf("could not find provider for %#v", requirement)
				} else if count > 1 {
					return fmt.Errorf("found multiple providers having same name for %#v", requirement)
				}
				ps.AddProvider(p)
				err = pg.visit(ps, p)
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
				err := pg.visit(ps, provider)
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

func findInterfaceProvider(req *store.ProviderRequirement) (*store.Provider, error) {
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
		implements, err := store.FindImplementsByInterface(req.CmpPkgPath, req.CmpTypName)
		if err != nil {
			return nil, err
		}
		if len(implements) != 1 {
			return nil, fmt.Errorf("could not determine implementation for interface kind. req: %v, implementations: %v", req, implements)
		}
		// find by implementation
		providers1, err1 := store.FindProviderByComponent(implements[0].CmpPkgPath, implements[0].CmpTypName)
		// find by interface directly
		providers2, err2 := store.FindProviderByComponent(req.CmpPkgPath, req.CmpTypName)
		return checked(append(providers1, providers2...), errors.Join(err1, err2))
	}
	providers, err := store.FindProviderByName(req.CmpPvdName)
	cp := make([]store.Provider, 0, 1)
	for _, provider := range providers {
		impl := &store.ImplementStmt{
			CmpPkgPath:   provider.CmpPkgPath,
			CmpTypName:   provider.CmpTypName,
			CmpKind:      provider.CmpKind,
			IfacePkgPath: req.CmpPkgPath,
			IfaceName:    req.CmpTypName,
		}
		implements, err := store.FindImplementsByImpl(impl)
		if err != nil {
			return nil, err
		}
		if len(implements) == 1 {
			cp = append(cp, provider)
		}
	}
	return checked(cp, err)
}
