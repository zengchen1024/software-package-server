package postgresql

type dbTable struct {
	name string
}

func NewDBTable(name string) dbTable {
	return dbTable{name: name}
}

func (d dbTable) Insert(filter, result interface{}) (rows int64, err error) {
	query := db.Table(d.name).Where(filter).FirstOrCreate(result)

	err = query.Error
	rows = query.RowsAffected

	return
}
