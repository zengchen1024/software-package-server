package main

import (
	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	pkgPhasePR       = "pr"
	pkgPhaseRepo     = "repo"
	pkgPhaseInitCode = "init_code"
)

// msgToHandlePkg
type msgToHandlePkg struct {
	Phase    string           `json:"phase"`
	PkgId    string           `json:"pkg_id"`
	PR       msgOfPkgPR       `json:"pr"`
	Repo     msgOfPkgRepo     `json:"repo"`
	InitCode msgOfPkgInitCode `json:"init_code"`
}

// msgOfPkgPR
type msgOfPkgPR struct {
	RelevantPR    string `json:"relevant_pr"`
	Merged        bool   `json:"merged"`
	DuplicatedPkg bool   `json:"duplicated_pkg"`
	OtherReason   string `json:"other_reason"`
}

// msgOfPkgRepo
type msgOfPkgRepo struct {
	Platform     string `json:"platform"`
	RepoLink     string `json:"repo_link"`
	FailedReason string `json:"failed_reason"`
}

// msgOfPkgInitCode
type msgOfPkgInitCode = msgOfPkgRepo

// msgToHandlePkgPRCIChecked
type msgToHandlePkgPRCIChecked struct {
	PkgId        string `json:"pkg_id"`
	RelevantPR   string `json:"relevant_pr"`
	PRNum        int    `json:"pr_num"`
	FailedReason string `json:"failed_reason"`
}

func (msg *msgToHandlePkgPRCIChecked) toCmd() (cmd app.CmdToHandlePkgPRCIChecked, err error) {
	if cmd.RelevantPR, err = dp.NewURL(msg.RelevantPR); err != nil {
		return
	}

	cmd.PkgId = msg.PkgId
	cmd.FiledReason = msg.FailedReason
	cmd.PRNum = msg.PRNum

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

// msgToHandlePkgPRClosed
type msgToHandlePkgPRClosed struct {
	PkgId      string `json:"pkg_id"`
	Reason     string `json:"reason"`
	RejectedBy string `json:"rejected_by"`
}

func (msg *msgToHandlePkgPRClosed) toCmd() app.CmdToHandlePkgPRClosed {
	return app.CmdToHandlePkgPRClosed{
		PkgId:      msg.PkgId,
		Reason:     msg.Reason,
		RejectedBy: msg.RejectedBy,
	}
}

// msgToHandlePkgPRMerged
type msgToHandlePkgPRMerged struct {
	PkgId      string   `json:"pkg_id"`
	ApprovedBy []string `json:"approved_by"`
}

func (msg *msgToHandlePkgPRMerged) toCmd() app.CmdToHandlePkgPRMerged {
	return app.CmdToHandlePkgPRMerged{
		PkgId:      msg.PkgId,
		ApprovedBy: msg.ApprovedBy,
	}
}
