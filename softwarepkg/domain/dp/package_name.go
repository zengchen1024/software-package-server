package dp

import (
	"errors"
	"fmt"

	"github.com/opensourceways/software-package-server/utils"
)

type PackageName interface {
	PackageName() string
}

func NewPackageName(v string) (PackageName, error) {
	if v == "" || !reName.MatchString(v) {
		return nil, errors.New("invalid package name")
	}

	if max := config.MaxLengthOfPackageName; utils.StrLen(v) > max {
		return nil, fmt.Errorf(
			"the length of package name should be less than %d", max,
		)
	}

	return packageName(v), nil
}

type packageName string

func (v packageName) PackageName() string {
	return string(v)
}

func IsSamePkgName(n1, n2 PackageName) bool {
	return n1 != nil && n2 != nil && n1.PackageName() == n2.PackageName()
}
