package dp

import "errors"

type PackageLicense interface {
	PackageLicense() string
}

func NewPackageLicense(v string) (PackageLicense, error) {
	if v == "" {
		return nil, errors.New("empty package license")
	}

	return packageLicense(v), nil
}

type packageLicense string

func (r packageLicense) PackageLicense() string {
	return string(r)
}
