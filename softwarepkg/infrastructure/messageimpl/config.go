package messageimpl

import "github.com/opensourceways/software-package-server/common/infrastructure/kafka"

type Config struct {
	kafka.Config

	Topics Topics `json:"topics"  required:"true"`
}

type Topics struct {
	ApplyingSoftwarePkg           string `json:"applying_software_pkg"              required:"true"`
	ApprovedSoftwarePkg           string `json:"approved_software_pkg"              required:"true"`
	RejectedSoftwarePkg           string `json:"rejected_software_pkg"              required:"true"`
	AbandonedSoftwarePkg          string `json:"abandoned_software_pkg"             required:"true"`
	AlreadyClosedSoftwarePkg      string `json:"already_closed_software_pkg"        required:"true"`
	IndirectlyApprovedSoftwarePkg string `json:"indirectly_approved_software_pkg"   required:"true"`
}
