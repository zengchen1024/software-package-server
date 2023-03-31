package app

import commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"

const (
	errorSoftwarePkgExists          = "software_pkg_exists"
	errorSoftwarePkgNotFound        = "software_pkg_not_found"
	errorSoftwarePkgNoPermission    = "software_pkg_no_permission"
	errorSoftwarePkgCannotComment   = "software_pkg_cannot_comment"
	errorSoftwarePkgCommentIllegal  = "software_pkg_comment_illegal"
	errorSoftwarePkgCommentNotFound = "software_pkg_comment_not_found"
)

func errorCodeForFindingPkg(err error) string {
	if commonrepo.IsErrorResourceNotFound(err) {
		return errorSoftwarePkgNotFound
	}

	return ""
}
