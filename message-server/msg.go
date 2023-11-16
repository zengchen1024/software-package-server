package main

import (
	"github.com/opensourceways/software-package-server/softwarepkg/app"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

func cmdToDownloadPkgCode(data []byte) (cmd app.CmdToDownloadPkgCode, err error) {
	v, err := domain.UnmarshalToSoftwarePkgAppliedEvent(data)
	if err == nil {
		cmd.PkgId = v.PkgId
	}

	return
}

func cmdToStartCI(data []byte) (cmd app.CmdToStartCI, err error) {
	v, err := domain.UnmarshalToSoftwarePkgAppliedEvent(data)
	if err == nil {
		cmd.PkgId = v.PkgId
	}

	return
}

// msgToHandlePkgCIDone
type msgToHandlePkgCIDone struct {
	PkgId    string `json:"pkg_id"`
	Detail   string `json:"detail"`
	PRNumber int    `json:"number"`
	Success  bool   `json:"success"`
}

func (msg *msgToHandlePkgCIDone) toCmd() app.CmdToHandlePkgCIDone {
	return app.CmdToHandlePkgCIDone{
		PkgId:    msg.PkgId,
		Detail:   msg.Detail,
		Success:  msg.Success,
		PRNumber: msg.PRNumber,
	}
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

// msgToHandlePkgRepoCodePushed
type msgToHandlePkgRepoCodePushed struct {
	PkgId string `json:"pkg_id"`
}

func (msg *msgToHandlePkgRepoCodePushed) toCmd() (cmd app.CmdToHandlePkgRepoCodePushed, err error) {
	cmd.PkgId = msg.PkgId

	return
}

// cmdToHandlePkgAlreadyExisted
func cmdToHandlePkgAlreadyExisted(data []byte) (cmd app.CmdToHandlePkgAlreadyExisted, err error) {
	v, err := domain.UnmarshalToSoftwarePkgAlreadyExistEvent(data)
	if err != nil {
		return
	}

	cmd.PkgName, err = dp.NewPackageName(v.PkgName)

	return
}
