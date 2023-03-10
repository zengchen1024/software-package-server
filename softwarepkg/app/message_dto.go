package app

import "github.com/opensourceways/software-package-server/softwarepkg/domain/dp"

type CmdToHandleCIChecking struct {
	PkgId       string
	RelevantPR  dp.URL
	FiledReason string
}

func (cmd *CmdToHandleCIChecking) isSuccess() bool {
	return cmd.FiledReason == ""
}
