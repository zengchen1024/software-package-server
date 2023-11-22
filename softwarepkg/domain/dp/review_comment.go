package dp

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/opensourceways/software-package-server/utils"
)

const (
	cmdAPPROVE = "APPROVE"
	cmdReject  = "REJECT"
)

var (
	validCmds = map[string]bool{
		cmdAPPROVE: true,
		cmdReject:  true,
	}

	commandRegex = regexp.MustCompile(`(?m)^/([^\s]+)[\t ]*([^\n\r]*)`)

	sensitiveWords sensitiveWordsValidator
)

type sensitiveWordsValidator interface {
	CheckSensitiveWords(string) error
}

type ReviewComment interface {
	ReviewComment() string
	ParseReviewComment() (isCmd, isApprove bool)
}

func NewReviewComment(v string) (ReviewComment, error) {
	if err := checkReviewComment(v); err != nil {
		return nil, err
	}

	if err := sensitiveWords.CheckSensitiveWords(v); err != nil {
		return nil, err
	}

	return reviewComment(v), nil
}

func NewReviewCommentInternal(v string) (ReviewComment, error) {
	if err := checkReviewComment(v); err != nil {
		return nil, err
	}

	return reviewComment(v), nil
}

func checkReviewComment(v string) error {
	if v == "" {
		return errors.New("empty review comment")
	}

	if max := config.MaxLengthOfReviewComment; utils.StrLen(v) > max {
		return fmt.Errorf(
			"the length of review comment should be less than %d", max,
		)
	}

	return nil
}

type reviewComment string

func (v reviewComment) ReviewComment() string {
	return string(v)
}

func (v reviewComment) ParseReviewComment() (isCmd, isApprove bool) {
	if cmd := parseReviewCommand(string(v)); cmd != "" {
		isCmd = true
		isApprove = cmd == cmdAPPROVE
	}

	return
}

func parseReviewCommand(comment string) string {
	v := parseCommentCommands(comment)
	n := len(v)
	if n == 0 {
		return ""
	}

	for i := n - 1; i >= 0; i-- {
		if validCmds[v[i]] {
			return v[i]
		}
	}

	return ""
}

func parseCommentCommands(comment string) (r []string) {
	items := commandRegex.FindAllStringSubmatch(comment, -1)
	for i := range items {
		r = append(r, strings.ToUpper(items[i][1]))
	}

	return
}
