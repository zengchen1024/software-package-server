package dp

const (
	checkItemReviewResultPass    = "pass"
	checkItemReviewResultNoIdea  = "no_idea"
	checkItemReviewResultNotPass = "not_pass"
)

var (
	CheckItemPass    = checkItemReviewResult(checkItemReviewResultPass)
	CheckItemNoIdea  = checkItemReviewResult(checkItemReviewResultNoIdea)
	CheckItemNotPass = checkItemReviewResult(checkItemReviewResultNotPass)
)

type CheckItemReviewResult interface {
	CheckItemReviewResult() string
}

type checkItemReviewResult string

func (v checkItemReviewResult) CheckItemReviewResult() string {
	return string(v)
}

func IsCheckItemPass(v CheckItemReviewResult) bool {
	return v != nil && v.CheckItemReviewResult() == checkItemReviewResultPass
}
