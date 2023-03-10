package app

import (
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

type CmdToHandleCIChecking struct {
	PkgId       string
	RelevantPR  dp.URL
	FiledReason string
}

func (cmd *CmdToHandleCIChecking) isSuccess() bool {
	return cmd.FiledReason == ""
}

func (cmd *CmdToHandleCIChecking) logString() string {
	return fmt.Sprintf(
		"handling ci checking, pkgid:%s, pr:%s",
		cmd.PkgId, cmd.RelevantPR.URL(),
	)
}
