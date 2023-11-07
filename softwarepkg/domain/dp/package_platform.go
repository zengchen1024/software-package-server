package dp

import "errors"

type PackagePlatform interface {
	PackagePlatform() string
	IsLocalPlatform() bool
	RepoLink(PackageName) URL
}

func NewPackagePlatform(v string) (PackagePlatform, error) {
	if !config.isValidPlatform(v) {
		return nil, errors.New("invalid package platform")
	}

	return packagePlatform(v), nil
}

type packagePlatform string

func (v packagePlatform) PackagePlatform() string {
	return string(v)
}

func (v packagePlatform) IsLocalPlatform() bool {
	return string(v) == config.LocalPlatform
}

func (v packagePlatform) RepoLink(name PackageName) URL {
	return dpURL(config.SupportedPlatforms[string(v)] + name.PackageName())
}

func IsSamePlatform(a, b PackagePlatform) bool {
	return a != nil && b != nil && a.PackagePlatform() == b.PackagePlatform()
}
