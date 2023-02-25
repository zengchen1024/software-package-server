package message

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

type SoftwarePkgMessage interface {
	NotifyPkgApproved(*domain.SoftwarePkgApprovedEvent) error
	NotifyPkgRejected(*domain.SoftwarePkgRejectedEvent) error
}
