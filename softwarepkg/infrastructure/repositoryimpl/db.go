package repositoryimpl

type dbClient interface {
	Insert(filter, result interface{}) error

	IsRowNotExists(err error) bool
	IsRowExists(err error) bool
}
