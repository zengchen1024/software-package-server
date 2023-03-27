package messageimpl

import "github.com/opensourceways/software-package-server/common/infrastructure/kafka"

type Config struct {
	kafka.Config

	Topics Topics `json:"topics"  required:"true"`
}

type Topics struct {
	NewSoftwarePkg            string `json:"new_software_pkg"              required:"true"`
	ApprovedSoftwarePkg       string `json:"approved_software_pkg"         required:"true"`
	RejectedSoftwarePkg       string `json:"rejected_software_pkg"         required:"true"`
	AbandonedSoftwarePkg      string `json:"abandoned_software_pkg"        required:"true"`
	AlreadyExistedSoftwarePkg string `json:"already_existed_software_pkg"    required:"true"`
}
