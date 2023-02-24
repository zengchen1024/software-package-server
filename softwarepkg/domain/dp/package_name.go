package dp

import (
	"errors"
	"fmt"

	"github.com/opensourceways/src-package-server/utils"
)

type PackageName interface {
	PackageName() string
}

func NewPackageName(v string) (PackageName, error) {
	if v == "" {
		return nil, errors.New("empty package name")
	}

	// TODO: check by regexp

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
