package store

type ImplementStmt struct {
	// component related information
	CmpPkgPath string `db:"cmp_pkg_path"`
	CmpTypName string `db:"cmp_typ_name"`
	// component kind can be found at [types.TypeKind]
	// see /dot/internal/types/util.go
	CmpKind      int    `db:"cmp_kind"`
	IfacePkgPath string `db:"iface_pkg_path"`
	// name is the id field
	IfaceName string `db:"iface_name"`
	Labels    string `db:"labels"`
}

func (*ImplementStmt) TableName() string {
	return "implement_stmts"
}

// TableProvider
// 该表存储Provider的信息和它所能提供的component
const TableImplementStmts = `
create table implement_stmts (
	cmp_pkg_path     text,
	cmp_typ_name     text,
	cmp_kind         integer,
	iface_pkg_path   text,
	iface_name       text,
	labels           text,
	CONSTRAINT UC_Provider UNIQUE(cmp_pkg_path, cmp_typ_name, cmp_kind, iface_pkg_path, iface_name)
)
`

const SqlDeleteImplementStmtById = `
	delete from implement_stmts
	where 1 = 1
		and cmp_pkg_path = :cmp_pkg_path
		and cmp_typ_name = :cmp_typ_name
		and cmp_kind = :cmp_kind
		and iface_pkg_path = :iface_pkg_path
		and iface_name = :iface_name
`

func SaveImplement(impl *ImplementStmt) error {
	err := DeleteImplement(impl)
	if err != nil {
		return err
	}
	return Insert(impl)
}

func DeleteImplement(impl *ImplementStmt) error {
	_, err := db.NamedExec(SqlDeleteImplementStmtById, impl)
	return err
}

func FindImplementsByInterface(interfacePackage string, interfaceName string) ([]ImplementStmt, error) {
	return NamedSelect[ImplementStmt](
		"select * from implement_stmts where iface_pkg_path = :pkg and iface_name = :name",
		Q{"pkg": interfacePackage, "name": interfaceName},
	)
}

const SqlFindImplementStmtById = `
	select * from implement_stmts
	where 1 = 1
		and cmp_pkg_path = :cmp_pkg_path
		and cmp_typ_name = :cmp_typ_name
		and cmp_kind = :cmp_kind
		and iface_pkg_path = :iface_pkg_path
		and iface_name = :iface_name
`

func FindImplementsByImpl(impl *ImplementStmt) ([]ImplementStmt, error) {
	return NamedSelect[ImplementStmt](SqlFindImplementStmtById, impl)
}
