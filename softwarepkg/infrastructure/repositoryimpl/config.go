package repositoryimpl

type Config struct {
	Table Table `json:"table" required:"true"`
}

type Table struct {
	SoftwarePkg   string `json:"software_pkg"    required:"true"`
	ReviewComment string `json:"review_comment"  required:"true"`
}
