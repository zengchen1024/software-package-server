package repositoryimpl

type dbClient interface {
	Insert(filter, result interface{}) (rows int64, err error)
}
