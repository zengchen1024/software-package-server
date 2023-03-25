package main

import (
	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

// msgToHandlePkgCIChecked
type msgToHandlePkgCIChecked struct {
	PkgId        string `json:"pkg_id"`
	RelevantPR   string `json:"relevant_pr"`
	PRNum        int    `json:"pr_num"`
	FailedReason string `json:"failed_reason"`
}

func (msg *msgToHandlePkgCIChecked) toCmd() (cmd app.CmdToHandlePkgCIChecked, err error) {
	if cmd.RelevantPR, err = dp.NewURL(msg.RelevantPR); err != nil {
		return
	}

	cmd.PkgId = msg.PkgId
	cmd.FiledReason = msg.FailedReason
	cmd.PRNum = msg.PRNum

	return
}

// msgToHandlePkgInitialized
type msgToHandlePkgInitialized struct {
	PkgId      string `json:"pkg_id"`
	RelevantPR string `json:"relevant_pr"`
	// RepoLink is the one of already existed pkg
	RepoLink     string `json:"repo_link"`
	FailedReason string `json:"failed_reason"`
}

func (msg *msgToHandlePkgInitialized) toCmd() (cmd app.CmdToHandlePkgInitialized, err error) {
	cmd.PkgId = msg.PkgId
	cmd.FiledReason = msg.FailedReason

	if msg.RelevantPR != "" {
		if cmd.RelevantPR, err = dp.NewURL(msg.RelevantPR); err != nil {
			return
		}
	}

	if msg.RepoLink != "" {
		cmd.RepoLink, err = dp.NewURL(msg.RepoLink)
	}

	return
}

// msgToHandlePkgRepoCreated
type msgToHandlePkgRepoCreated struct {
	PkgId        string `json:"pkg_id"`
	Platform     string `json:"platform"`
	RepoLink     string `json:"repo_link"`
	FailedReason string `json:"failed_reason"`
}

func (msg *msgToHandlePkgRepoCreated) toCmd() (cmd app.CmdToHandlePkgRepoCreated, err error) {
	cmd.PkgId = msg.PkgId
	cmd.FiledReason = msg.FailedReason

	if cmd.Platform, err = dp.NewPackagePlatform(msg.Platform); err != nil {
		return
	}

	if msg.RepoLink != "" {
		cmd.RepoLink, err = dp.NewURL(msg.RepoLink)
	}

	return
}

// msgToHandlePkgCodeSaved
type msgToHandlePkgCodeSaved = msgToHandlePkgRepoCreated
