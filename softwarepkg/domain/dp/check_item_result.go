package dp

const (
	checkItemResultPass    = "pass"
	checkItemResultNoIdea  = "no_idea"
	checkItemResultNotPass = "not_pass"
)

var (
	CheckItemPass    = checkItemResult(checkItemResultPass)
	CheckItemNoIdea  = checkItemResult(checkItemResultNoIdea)
	CheckItemNotPass = checkItemResult(checkItemResultNotPass)
)

type CheckItemResult interface {
	CheckItemResult() string
}

type checkItemResult string

func (v checkItemResult) CheckItemResult() string {
	return string(v)
}

func IsCheckItemPass(v CheckItemResult) bool {
	return v != nil && v.CheckItemResult() == checkItemResultPass
}
