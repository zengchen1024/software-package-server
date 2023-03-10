package main

import "github.com/opensourceways/software-package-server/softwarepkg/app"

type msgToHandleCIChecking struct {
	PkgId      string `json:"pkg_id"`
	RelevantPR string `json:"relevant_pr"`
}

func (msg *msgToHandleCIChecking) toCmd() (cmd app.CmdToHandleCIChecking, err error) {
	return
}
