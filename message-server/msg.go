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
