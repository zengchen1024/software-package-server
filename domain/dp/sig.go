package dp

import "errors"

type Sig interface {
	Sig() string
}

func NewSig(v string) (Sig, error) {
	if v == "" {
		return nil, errors.New("empty package license")
	}

	return dpSig(v), nil
}

type dpSig string

func (r dpSig) Sig() string {
	return string(r)
}
