package pkgciimpl

type Config struct {
	CIEndpoint string `json:"ci_endpoint" required:"true"`
	CIService  string `json:"ci_service"  required:"true"`
}
