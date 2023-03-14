package app

import (
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

// CmdToHandleCIChecking
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

// CmdToHandleRepoCreated
type CmdToHandleRepoCreated struct {
	PkgId       string
	FiledReason string

	domain.RepoCreatedInfo
}

func (cmd *CmdToHandleRepoCreated) isSuccess() bool {
	return cmd.FiledReason == ""
}

func (cmd *CmdToHandleRepoCreated) logString() string {
	if !cmd.isSuccess() {
		return ""
	}

	return fmt.Sprintf(
		"handling repo created, pkgid:%s, platform:%s, repo:%s",
		cmd.PkgId, cmd.Platform.PackagePlatform(), cmd.RepoLink.URL(),
	)
}

// CmdToHandlePkgRejected
type CmdToHandlePkgRejected struct {
	PkgId      string
	Reason     string
	RejectedBy string
}

func (cmd *CmdToHandlePkgRejected) logString() string {
	return fmt.Sprintf(
		"handling pkg rejected, pkgid:%s, reason:%s, rejected by:%s",
		cmd.PkgId, cmd.Reason, cmd.RejectedBy,
	)
}
