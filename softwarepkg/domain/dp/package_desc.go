package dp

import (
	"errors"
	"fmt"

	"github.com/opensourceways/software-package-server/utils"
)

type PackageDesc interface {
	PackageDesc() string
}

func NewPackageDesc(v string) (PackageDesc, error) {
	if v == "" {
		return nil, errors.New("empty package desc")
	}

	if max := config.MaxLengthOfPackageDesc; utils.StrLen(v) > max {
		return nil, fmt.Errorf(
			"the length of package desc should be less than %d", max,
		)
	}

	return packageDesc(v), nil
}

type packageDesc string

func (v packageDesc) PackageDesc() string {
	return string(v)
}
