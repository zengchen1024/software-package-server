package repositoryimpl

import "github.com/opensourceways/software-package-server/common/infrastructure/postgresql"

type dbClient interface {
	Insert(filter, result interface{}) error
	Counts(filter interface{}) (int, error)
	GetTableRecords(filter, result interface{}, p postgresql.Pagination, sort []postgresql.SortByColumn) (err error)
	GetTableRecord(filter, result interface{}) error

	IsRowNotExists(err error) bool
	IsRowExists(err error) bool
}
