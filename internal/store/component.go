package store

type Component struct {
	CmpPkgPath string
	CmpPkgName string
	// the type name eg: Book
	CmpTypName string
	CmpName    string
	// struct or interface or primitive types
	CmpKind int
}

func (c *Component) TableName() string {
	return "components"
}

// TableComponents
// 该表存储每个Component的信息
const TableComponents = `
create table components (
	cmp_pkg_path text,
	cmp_pkg_name text,
	cmp_typ_name text,
	cmp_name     text,
	cmp_kind     integer,
  	CONSTRAINT UC_Component UNIQUE(cmp_pkg_path, cmp_typ_name, cmp_name)
)
`

const SqlFindComponentById = `
	select * from components t 
	where 1 = 1
		and t.cmp_pkg_path = ? 
		and t.cmp_typ_name = ? 
		and t.cmp_name = ?
`

const SqlCountComponentById = `
	select count(1) from components t 
	where 1 = 1
		and t.cmp_pkg_path = ? 
		and t.cmp_typ_name = ? 
		and t.cmp_name = ?
`

func SaveComponent(c *Component) error {
	row := db.QueryRow(SqlCountComponentById, c.CmpPkgPath, c.CmpTypName, c.CmpName)
	if err := row.Err(); err != nil {
		return err
	}
	n := 0
	err := row.Scan(&n)
	if err != nil {
		return err
	}
	if n < 1 {
		return Insert(c)
	}
	return nil
}
