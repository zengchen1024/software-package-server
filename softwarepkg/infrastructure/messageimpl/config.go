package messageimpl

type Topics struct {
	SoftwarePkgClosed         string `json:"software_pkg_closed"           required:"true"`
	SoftwarePkgApplied        string `json:"software_pkg_applied"          required:"true"`
	SoftwarePkgRetested       string `json:"software_pkg_retested"         required:"true"`
	SoftwarePkgAlreadyExisted string `json:"software_pkg_already_existed"  required:"true"`

	// importer edited the pkg and want to reload code
	SoftwarePkgCodeUpdated string `json:"software_pkg_code_updated"        required:"true"`
}
