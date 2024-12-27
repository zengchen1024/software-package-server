package domain

const (
	PkgStatusInitialized = "initialized"
	PkgStatusPRCreated   = "pr_created"
	PkgStatusPRMerged    = "pr_merged"
	PkgStatusDone        = "done"
	PkgStatusException   = "exception" // more information in the email of maintainer
)

var PkgStatusNeedToHandle = []string{PkgStatusInitialized, PkgStatusPRCreated, PkgStatusPRMerged}

type PkgWatch struct {
	Id     string
	Status string
	PR     PullRequest
}

type PullRequest struct {
	Num  int
	Link string
}

func (r *PkgWatch) SetPkgStatusPRCreated() {
	r.Status = PkgStatusPRCreated
}

func (r *PkgWatch) SetPkgStatusPRMerged() {
	r.Status = PkgStatusPRMerged
}

func (r *PkgWatch) SetPkgStatusDone() {
	r.Status = PkgStatusDone
}

func (r *PkgWatch) SetPkgStatusException() {
	r.Status = PkgStatusException
}

func (r *PkgWatch) IsPkgStatusMerged() bool {
	return r.Status == PkgStatusPRMerged
}

func (r *PkgWatch) IsPkgStatusInitialized() bool {
	return r.Status == PkgStatusInitialized
}

func (r *PkgWatch) IsPkgInProgress() bool {
	return r.Status == PkgStatusPRCreated || r.Status == PkgStatusPRMerged
}
