package postgresql

import "errors"

var (
	errRowExists    = errors.New("row exists")
	errRowNotExists = errors.New("row doesn't exist")
)

type dbTable struct {
	name string
}

func NewDBTable(name string) dbTable {
	return dbTable{name: name}
}

func (t dbTable) Insert(filter, result interface{}) error {
	query := db.Table(t.name).Where(filter).FirstOrCreate(result)

	if err := query.Error; err != nil {
		return err
	}

	if query.RowsAffected == 0 {
		return errRowExists
	}

	return nil
}

func (t dbTable) IsRowNotExists(err error) bool {
	return errors.Is(err, errRowNotExists)
}

func (t dbTable) IsRowExists(err error) bool {
	return errors.Is(err, errRowExists)
}
