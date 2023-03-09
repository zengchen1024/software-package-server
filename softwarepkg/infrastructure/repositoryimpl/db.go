package repositoryimpl

import "github.com/opensourceways/software-package-server/common/infrastructure/postgresql"

type dbClient interface {
	Insert(filter, result interface{}) error
	Count(filter interface{}) (int, error)
	GetRecords(filter, result interface{}, p postgresql.Pagination, sort []postgresql.SortByColumn) error
	GetRecord(filter, result interface{}) error
	UpdateRecord(filter, update interface{}) error

	IsRowNotFound(err error) bool
	IsRowExists(err error) bool
}
