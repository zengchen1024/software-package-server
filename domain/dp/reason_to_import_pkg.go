package dp

import (
	"errors"
	"fmt"
)

type ReasonToImportPkg interface {
	ReasonToImportPkg() string
}

func NewReasonToImportPkg(v string) (ReasonToImportPkg, error) {
	if v == "" {
		return nil, errors.New("empty reason")
	}

	if max := config.MaxLengthOfReasonToImportPkg; utils.StrLen(v) > max {
		return nil, fmt.Errorf(
			"the length of package desc should be less than %d", max,
		)
	}

	return reasonToImportPkg(v), nil
}

type reasonToImportPkg string

func (r reasonToImportPkg) ReasonToImportPkg() string {
	return string(r)
}
