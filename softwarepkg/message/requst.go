package message

type RequestToHandleCIPassed struct {
	PkgId      string `json:"pkg_id"`
	RelevantPR string `json:"relevant_pr"`
}
