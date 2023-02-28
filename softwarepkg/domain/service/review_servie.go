package service

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/message"
)

type SoftwarePkgReviewService interface {
	ApprovePkg(pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account) bool

	RejectPkg(pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account) bool
}

type reviewService struct {
	message message.SoftwarePkgMessage
}

func (s *reviewService) ApprovePkg(
	pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account,
) bool {
	changed, approved := pkg.ApproveBy(user)
	if approved {
		_ = s.message.NotifyPkgApproved(&domain.SoftwarePkgApprovedEvent{})
		// TODO if failed , log it
		// Event handler should check if the pkg is approved actualy.
	}

	return changed
}

func (s *reviewService) RejectPkg(
	pkg *domain.SoftwarePkgBasicInfo, version int, user dp.Account,
) bool {
	changed, rejected := pkg.RejectBy(user)
	if rejected {
		s.message.NotifyPkgRejected(&domain.SoftwarePkgRejectedEvent{})
		// TODO
	}

	return changed
}
