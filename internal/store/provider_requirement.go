package store

// ProviderRequirement
type ProviderRequirement struct {
	PvdPkgPath string `db:"pvd_pkg_path"`
	PvdPkgName string `db:"pvd_pkg_name"`
	PvdOriName string `db:"pvd_ori_name"`
	PvdName    string `db:"pvd_name"`
	PvdKind    string `db:"pvd_kind"`
	CmpPkgPath string `db:"cmp_pkg_path"`
	CmpPkgName string `db:"cmp_pkg_name"`
	CmpTypName string `db:"cmp_typ_name"`
	// multiple params may have same type, thus use param name to distinct each other
	CmpName string `db:"cmp_name"`
	CmpKind int    `db:"cmp_kind"`
	// go:ioc --param name.provider="NewLiu"
	CmpPvdName string `db:"cmp_pvd_name"`
	Labels     string `db:"labels"`
}

func (*ProviderRequirement) TableName() string {
	return "provider_requirements"
}

// TableProviderRequirement
// 该表存储Provider的信息和它所需要的所有的components
const TableProviderRequirements = `
create table provider_requirements (
	pvd_pkg_path 	text,
	pvd_pkg_name 	text,	
	pvd_ori_name   text,
	pvd_name     	text,
	pvd_kind 	 	text,
	cmp_pkg_path 	text,
	cmp_pkg_name 	text,
	cmp_typ_name 	text,
	cmp_name 	    text,
	cmp_kind     	integer,
	cmp_pvd_name 	text,
	labels          text,
	CONSTRAINT UC_Provider_Requirements UNIQUE(pvd_pkg_path, pvd_ori_name, pvd_kind, cmp_pkg_path, cmp_typ_name, cmp_kind, cmp_name)
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
		and cmp_name = ?
`

func SaveProviderRequirement(c *ProviderRequirement) error {
	res, err := db.Exec(SqlDeleteProviderRequirementById, c.PvdPkgPath, c.PvdName, c.PvdKind, c.CmpPkgPath, c.CmpTypName, c.CmpKind, c.CmpName)
	if err != nil {
		return err
	}
	_ = res
	return Insert(c)
}

const SqlFindProviderRequirementByCmpType = `
	select * from provider_requirements t 
	where 1 = 1
		and t.pvd_pkg_path = ?
		and t.pvd_ori_name = ?
		and t.pvd_kind = ?
`

// FindProviderRequirements
//
//	select * from provider_requirements t
//	where 1 = 1
//		and t.pvd_pkg_path = ?
//		and t.pvd_ori_name = ?
//		and t.pvd_kind = ?
func FindProviderRequirements(c *Provider) ([]ProviderRequirement, error) {
	return Select[ProviderRequirement](SqlFindProviderRequirementByCmpType, c.PvdPkgPath, c.PvdOriName, c.PvdKind)
}

func DeleteProviderRequirementByPkg(pkgPath string) error {
	_, err := db.Exec("delete from provider_requirements where pvd_pkg_path = ?", pkgPath)
	return err
}
