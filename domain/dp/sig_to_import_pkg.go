package dp

import "errors"

type SigToImportPkg interface {
	SigToImportPkg() string
}

func NewSigToImportPkg(v string) (SigToImportPkg, error) {
	if v == "" {
		return nil, errors.New("empty sig")
	}

	return sigToImportPkg(v), nil
}

type sigToImportPkg string

func (v sigToImportPkg) SigToImportPkg() string {
	return string(v)
}
