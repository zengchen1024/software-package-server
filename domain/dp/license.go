package dp

import "errors"

type License interface {
	License() string
}

func NewLicense(v string) (License, error) {
	if v == "" {
		return nil, errors.New("empty  license")
	}

	return license(v), nil
}

type license string

func (v license) License() string {
	return string(v)
}
