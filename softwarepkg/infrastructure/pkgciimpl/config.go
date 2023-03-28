package pkgciimpl

type Config struct {
	CIEndpoint         string `json:"ci_endpoint"             required:"true"`
	CommonTestBodyFile string `json:"common_test_body_file"   required:"true"`
}
