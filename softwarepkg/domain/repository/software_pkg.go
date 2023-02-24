package repository

import "github.com/opensourceways/software-package-server/domain"

type SoftwarePkg interface {
	AddSoftwarePkg(*domain.SoftwarePkg) error
}
