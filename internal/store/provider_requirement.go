package store

import (
	"fmt"
)

// ProviderRequirement
type ProviderRequirement struct {
	PvdPkgPath  string `db:"pvd_pkg_path"`
	PvdPkgName  string `db:"pvd_pkg_name"`
	PvdFuncName string `db:"pvd_func_name"`
	PvdName     string `db:"pvd_name"`
	PvdKind     string `db:"pvd_kind"`
	CmpPkgPath  string `db:"cmp_pkg_path"`
	CmpPkgName  string `db:"cmp_pkg_name"`
	CmpTypName  string `db:"cmp_typ_name"`
	CmpKind     int    `db:"cmp_kind"`
	// go:ioc --param name.provider="NewLiu"
	CmpPvdName string `db:"cmp_pvd_name"`
	// go:ioc --param age.ident=123
	CmpIdentValue string `db:"cmp_ident_value"`
}

// TableProviderRequirement
// 该表存储Provider的信息和它所需要的所有的components
const TableProviderRequirements = `
create table provider_requirements (
	pvd_pkg_path 	text,
	pvd_pkg_name 	text,	
	pvd_func_name   text,
	pvd_name     	text,
	pvd_kind 	 	text,
	cmp_pkg_path 	text,
	cmp_pkg_name 	text,
	cmp_typ_name 	text,
	cmp_kind     	integer,
	cmp_pvd_name 	text,
	CONSTRAINT UC_Provider_Requirements UNIQUE(pvd_pkg_path, pvd_func_name, pvd_kind, cmp_pkg_path, cmp_typ_name, cmp_kind, cmp_pvd_name)
)
`

const SqlDeleteProviderRequirementById = `
	delete from provider_requirements
	where 1 = 1
		and pvd_pkg_path = ?
		and pvd_name = ?
		and pvd_kind = ?
		and cmp_pkg_path = ? 
		and cmp_typ_name = ?
		and cmp_kind = ?
		and cmp_pvd_name = ?
`

const SqlInsertProviderRequirement = `
	insert into provider_requirements(pvd_pkg_path, pvd_pkg_name, pvd_func_name, pvd_name, pvd_kind, cmp_pkg_path, cmp_pkg_name, cmp_typ_name, cmp_kind, cmp_pvd_name)
	values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

func SaveProviderRequirement(c *ProviderRequirement) error {
	res, err := db.Exec(SqlDeleteProviderRequirementById, c.PvdPkgPath, c.PvdName, c.PvdKind, c.CmpPkgPath, c.CmpTypName, c.CmpKind, c.CmpPvdName)
	if err != nil {
		return err
	}
	_ = res

	res, err = db.Exec(
		SqlInsertProviderRequirement,
		c.PvdPkgPath, c.PvdPkgName, c.PvdFuncName, c.PvdName, c.PvdKind,
		c.CmpPkgPath, c.CmpPkgName, c.CmpTypName, c.CmpKind, c.CmpPvdName,
	)
	if err != nil {
		return fmt.Errorf("err: %w, record: %#v", err, c)
	}
	_ = res
	return nil
}

const SqlFindProviderRequirementByCmpType = `
	select * from provider_requirements t 
	where 1 = 1
		and t.pvd_pkg_path = ?
		and t.pvd_name = ?
		and t.pvd_kind = ?
`

// FindProviderRequirements
//
//	select * from provider_requirements t
//	where 1 = 1
//		and t.pvd_pkg_path = ?
//		and t.pvd_func_name = ?
//		and t.pvd_kind = ?
func FindProviderRequirements(c *Provider) ([]ProviderRequirement, error) {
	return Select[ProviderRequirement](SqlFindProviderRequirementByCmpType, c.PvdPkgPath, c.PvdFuncName, c.PvdKind)
}