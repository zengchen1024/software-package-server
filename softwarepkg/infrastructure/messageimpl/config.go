package messageimpl

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
)

type Config struct {
	kfklib.Config

	Topics Topics `json:"topics"  required:"true"`
}

type Topics struct {
	//ApprovedSoftwarePkg       string `json:"approved_software_pkg"         required:"true"`

	SoftwarePkgApplied        string `json:"software_pkg_applied"          required:"true"`
	SoftwarePkgRetested       string `json:"software_pkg_retested"         required:"true"`
	SoftwarePkgAlreadyExisted string `json:"software_pkg_already_existed"  required:"true"`

	// importer edited the pkg and want to reload code
	SoftwarePkgCodeUpdated string `json:"software_pkg_code_updated"        required:"true"`

	SoftwarePkgApproved  string
	SoftwarePkgRejected  string
	SoftwarePkgAbandoned string
}
