package app

import (
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

// CmdToHandlePkgInitialized
type CmdToHandlePkgInitialized struct {
	PkgId      string
	RelevantPR dp.URL
	PRNum      int
	// RepoLink is the one of already existed pkg
	RepoLink    string
	FiledReason string
}

func (cmd *CmdToHandlePkgInitialized) isSuccess() bool {
	return cmd.FiledReason == "" && cmd.RepoLink == ""
}

func (cmd *CmdToHandlePkgInitialized) isPkgAreadyExisted() bool {
	return cmd.RepoLink != ""
}

func (cmd *CmdToHandlePkgInitialized) logString() string {
	return fmt.Sprintf(
		"handling pkg init done, pkgid:%s, pr:%s",
		cmd.PkgId, cmd.RelevantPR.URL(),
	)
}

// CmdToHandlePkgRepoCreated
type CmdToHandlePkgRepoCreated struct {
	PkgId       string
	FiledReason string

	domain.RepoCreatedInfo
}

func (cmd *CmdToHandlePkgRepoCreated) isSuccess() bool {
	return cmd.FiledReason == ""
}

func (cmd *CmdToHandlePkgRepoCreated) logString() string {
	if !cmd.isSuccess() {
		return ""
	}

	return fmt.Sprintf(
		"handling pkg repo created, pkgid:%s, platform:%s, repo:%s",
		cmd.PkgId, cmd.Platform.PackagePlatform(), cmd.RepoLink.URL(),
	)
}
