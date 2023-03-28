package clavalidator

import "github.com/opensourceways/software-package-server/softwarepkg/domain/dp"

type ClaValidator interface {
	HasSignedCLA(dp.Email) (bool, error)
}
