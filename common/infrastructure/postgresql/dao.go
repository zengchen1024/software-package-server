package postgresql

import (
	"errors"
	"strings"
)

var (
	errRowExists    = errors.New("row exists")
	errRowNotExists = errors.New("row doesn't exist")
)

type SortByColumn struct {
	Column string
	Ascend bool
}

func (s SortByColumn) order() string {
	v := " ASC"
	if !s.Ascend {
		v = " DESC"
	}
	return s.Column + v
}

type Pagination struct {
	PageNum      int
	CountPerPage int
}

func (p Pagination) pagination() (limit, offset int) {
	limit = p.CountPerPage
	if limit > 0 {
		if p.PageNum > 0 {
			offset = (p.PageNum - 1) * limit
		}
	}

	return
}

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

func (t dbTable) GetTableRecords(
	filter, result interface{}, p Pagination,
	sort []SortByColumn,
) (err error) {
	query := db.Table(t.name).Where(filter)

	var orders []string
	for _, v := range sort {
		orders = append(orders, v.order())
	}

	if len(orders) >= 0 {
		query.Order(strings.Join(orders, ","))
	}

	if limit, offset := p.pagination(); limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	err = query.Find(result).Error

	return
}

func (t dbTable) Counts(filter interface{}) (int, error) {
	var total int64
	err := db.Table(t.name).Where(filter).Count(&total).Error

	return int(total), err
}

func (t dbTable) IsRowNotExists(err error) bool {
	return errors.Is(err, errRowNotExists)
}

func (t dbTable) IsRowExists(err error) bool {
	return errors.Is(err, errRowExists)
}
