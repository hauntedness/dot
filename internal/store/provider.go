package store

import "cmp"

type Provider struct {
	//
	PvdPkgPath string `db:"pvd_pkg_path"`
	// note that here the package name might be rewrote
	PvdPkgName string `db:"pvd_pkg_name"`
	// the original provider name. typically func name or struct name
	// func named is used to be write to the wire.Set
	PvdOriName string `db:"pvd_ori_name"`
	// name is the id field
	PvdName string `db:"pvd_name"`
	// provider kind, eg: use function or direct a variable
	PvdKind string `db:"pvd_kind"`
	// 1 represent that results contains error
	PvdError int `db:"pvd_error"`
	// component related information
	CmpPkgPath string `db:"cmp_pkg_path"`
	CmpPkgName string `db:"cmp_pkg_name"`
	CmpTypName string `db:"cmp_typ_name"`
	// component kind can be found at [types.TypeKind]
	// see /dot/internal/types/util.go
	CmpKind int    `db:"cmp_kind"`
	Labels  string `db:"labels"`
}

func (*Provider) TableName() string {
	return "providers"
}

func (a *Provider) Compare(b *Provider) int {
	return cmp.Compare(a.PvdPkgPath+a.PvdOriName, b.PvdPkgPath+b.PvdOriName)
}

// TableProvider
// 该表存储Provider的信息和它所能提供的component
const TableProviders = `
create table providers (
	pvd_pkg_path  text,
	pvd_pkg_name  text,
	pvd_ori_name  text,
	pvd_name      text,
	pvd_kind      text,
	pvd_error     text,
	cmp_pkg_path  text,
	cmp_pkg_name  text,
	cmp_typ_name  text,
	cmp_kind      integer,
	labels        text,
	CONSTRAINT UC_Provider UNIQUE(pvd_pkg_path, pvd_ori_name)
)
`

const SqlDeleteProviderById = `
	delete from providers
	where 1 = 1
		and pvd_pkg_path = ?
		and pvd_ori_name = ?
`

func SaveProvider(c *Provider) error {
	_, err := db.Exec(SqlDeleteProviderById, c.PvdPkgPath, c.PvdOriName)
	if err != nil {
		return err
	}
	return Insert(c)
}

// FindProviderByCmp
//
// find suitable providers by component
func FindProviderByCmp(pkg string, typ string, kind int) ([]Provider, error) {
	sql := `select * from providers t where t.cmp_pkg_path = ? and t.cmp_typ_name = ? and t.cmp_kind = ?`
	return Select[Provider](sql, pkg, typ, kind)
}

// FindProviderByCmpName
//
// find suitable provider by component
func FindProviderByName(component string) ([]Provider, error) {
	sql := `select * from providers t where t.pvd_name = ?`
	return Select[Provider](sql, component)
}

// FindProviderByPkg
//
//	select * from providers t where t.pvd_pkg_path = ?
func FindProviderByPkg(pkg string) ([]Provider, error) {
	return Select[Provider]("select * from providers t where t.pvd_pkg_path = ?", pkg)
}

func DeleteProviderByPkg(pkg string) error {
	_, err := db.Exec("delete from providers where pvd_pkg_path = ?", pkg)
	return err
}
