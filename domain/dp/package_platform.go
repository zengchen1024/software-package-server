package dp

import "errors"

type PackagePlatform interface {
	PackagePlatform() string
}

func NewPackagePlatform(v string) (PackagePlatform, error) {
	if v == "" {
		return nil, errors.New("empty package license")
	}

	return packagePlatform(v), nil
}

type packagePlatform string

func (r packagePlatform) PackagePlatform() string {
	return string(r)
}
