package clavalidator

import "github.com/opensourceways/software-package-server/softwarepkg/domain/dp"

type CLA struct {
	Signed bool `json:"signed"`
}

type ClaValidator interface {
	HasSignedCLA(dp.Email) (bool, error)
}
