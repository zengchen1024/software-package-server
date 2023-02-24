package dp

import "errors"

type ImportingPkgSig interface {
	ImportingPkgSig() string
}

func NewImportingPkgSig(v string) (ImportingPkgSig, error) {
	if v == "" {
		return nil, errors.New("empty sig")
	}

	return importingPkgSig(v), nil
}

type importingPkgSig string

func (v importingPkgSig) ImportingPkgSig() string {
	return string(v)
}
