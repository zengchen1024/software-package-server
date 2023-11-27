package app

import (
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

// CmdToDownloadPkgCode
type CmdToDownloadPkgCode struct {
	PkgId string
}

func (cmd *CmdToDownloadPkgCode) logString() string {
	return fmt.Sprintf(
		"downloading pkg code, pkgid:%s", cmd.PkgId,
	)
}

// CmdToStartCI
type CmdToStartCI struct {
	PkgId string
}

func (cmd *CmdToStartCI) logString() string {
	return fmt.Sprintf(
		"starting pkg ci, pkgid:%s", cmd.PkgId,
	)
}

// CmdToHandlePkgCIDone
type CmdToHandlePkgCIDone struct {
	PkgId    string
	Detail   string
	Success  bool
	PRNumber int
}

func (cmd *CmdToHandlePkgCIDone) logString() string {
	return fmt.Sprintf(
		"handling pkg ci done, pkgid:%s, pr number:%d",
		cmd.PkgId, cmd.PRNumber,
	)
}

// CmdToHandlePkgRepoCodePushed
type CmdToHandlePkgRepoCodePushed struct {
	PkgId string
}

func (cmd *CmdToHandlePkgRepoCodePushed) logString() string {
	return fmt.Sprintf(
		"handling pkg repo code pushed, pkgid:%s", cmd.PkgId,
	)
}

type CmdToHandlePkgAlreadyExisted struct {
	PkgName dp.PackageName
}

func (cmd *CmdToHandlePkgAlreadyExisted) logString() string {
	return fmt.Sprintf(
		"importing existed pkg, pkg name:%s", cmd.PkgName.PackageName(),
	)
}
