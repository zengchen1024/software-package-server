package domain

import (
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

// SoftwarePkgOperationLog
type SoftwarePkgOperationLog struct {
	Time   int64
	User   dp.Account
	Action dp.PackageOperationLogAction
}

func (log *SoftwarePkgOperationLog) String() string {
	return fmt.Sprintf(
		"%s %s at %s",
		log.User.Account(),
		log.Action.PackageOperationLogAction(),
		utils.ToDateTime(log.Time),
	)
}

func NewSoftwarePkgOperationLog(
	user dp.Account, action dp.PackageOperationLogAction,
) SoftwarePkgOperationLog {
	return SoftwarePkgOperationLog{
		Time:   utils.Now(),
		User:   user,
		Action: action,
	}
}
