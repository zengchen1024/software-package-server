package dp

import "errors"

type PackagePlatform interface {
	PackagePlatform() string
}

func NewPackagePlatform(v string) (PackagePlatform, error) {
	if v == "" {
		return nil, errors.New("empty package platform")
	}

	return packagePlatform(v), nil
}

type packagePlatform string

func (v packagePlatform) PackagePlatform() string {
	return string(v)
}
