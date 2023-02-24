package repository

import "github.com/opensourceways/software-package-server/softwarepkg/domain"

type SoftwarePkg interface {
	AddSoftwarePkg(*domain.SoftwarePkg) error
}
