package dp

import (
	"errors"
	"fmt"

	"github.com/opensourceways/software-package-server/utils"
)

type PurposeToImportPkg interface {
	PurposeToImportPkg() string
}

func NewPurposeToImportPkg(v string) (PurposeToImportPkg, error) {
	if v == "" {
		return nil, errors.New("empty purpose")
	}

	if max := config.MaxLengthOfPurposeToImportPkg; utils.StrLen(v) > max {
		return nil, fmt.Errorf(
			"the length of purpose should be less than %d", max,
		)
	}

	return purposeToImportPkg(v), nil
}

type purposeToImportPkg string

func (v purposeToImportPkg) PurposeToImportPkg() string {
	return string(v)
}
