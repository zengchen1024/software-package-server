package softwarepkgadapter

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	mongodbCmdOr        = "$or"
	mongodbCmdIn        = "$in"
	mongodbCmdRegex     = "$regex"
	mongodbCmdElemMatch = "$elemMatch"

	mongodbCmdLt = "$lt"
)

type dao interface {
	InsertDocIfNotExists(filter, doc bson.M) (string, error)
	IsDocNotExists(error) bool
	IsDocExists(error) bool
	DocIdFilter(s string) (bson.M, error)
	GetDoc(filter, project bson.M, result interface{}) error
	UpdateDoc(filter bson.M, doc bson.M, version int) error
	Paginate(filter, project, sortBy bson.M, pageNum, countPerPage int64, result interface{}) error
}

func genDoc(doc interface{}) (m bson.M, err error) {
	v, err := json.Marshal(doc)
	if err != nil {
		return
	}

	err = json.Unmarshal(v, &m)

	return
}

type Config struct {
	SoftwarePkg string `json:"software_pkg" required:"true"`
}
