package app

import (
	"github.com/opensourceways/software-package-server/allerror"
	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
)

var errorSoftwarePkgExists = allerror.New(allerror.ErrorCodePkgExists, "pkg exists")

func parseErrorForFindingPkg(err error) error {
	if commonrepo.IsErrorResourceNotFound(err) {
		return allerror.NewNotFound(allerror.ErrorCodePkgNotFound, err.Error())
	}

	return err
}
