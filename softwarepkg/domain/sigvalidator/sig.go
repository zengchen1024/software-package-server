package sigvalidator

type Sig struct {
	EnFeature string `json:"en_feature"`
	Feature   string `json:"feature"`
	EnGroup   string `json:"en_group"`
	SigNames  string `json:"sig_names"`
	Group     string `json:"group"`
}

type SigValidator interface {
	IsValidSig(sig string) bool
	GetAll() []Sig
}
