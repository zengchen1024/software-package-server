package softwarepkgadapter

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	mongodbCmdLt = "$lt"
)

type dao interface {
	IsDocNotExists(error) bool
	IsDocExists(error) bool

	LikeFilter(v string, caseInsensitive bool) bson.M
	DocIdFilter(s string) (bson.M, error)

	InsertDocIfNotExists(filter, doc bson.M) (string, error)

	Count(filter bson.M) (n int64, err error)
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

type Collections struct {
	SoftwarePkg string `json:"software_pkg" required:"true"`
}
