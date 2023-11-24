package dp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opensourceways/software-package-server/utils"
)

var sensitiveWords sensitiveWordsValidator

type sensitiveWordsValidator interface {
	CheckSensitiveWords(string) error
}

type ReviewComment interface {
	ReviewComment() string
}

func NewReviewComment(v string) (ReviewComment, error) {
	return newReviewComment(v, config.MaxLengthOfReviewComment)
}

func NewReviewCommentInternal(v string) (ReviewComment, error) {
	return newReviewCommentInternal(v, config.MaxLengthOfReviewComment)
}

func CheckMultiComments(cs []string) error {
	maxLen := config.MaxLengthOfReviewComment

	s := strings.Join(cs, ".")
	n := utils.StrLen(s)
	if n == 0 {
		return nil
	}
	if n <= maxLen {
		return sensitiveWords.CheckSensitiveWords(s)
	}

	s = ""
	n = 0
	for _, c := range cs {
		n1 := utils.StrLen(c)

		if n2 := n1 + n; n2 > maxLen {
			if err := sensitiveWords.CheckSensitiveWords(s); err != nil {
				return err
			}

			s = c
			n = n1
		} else {
			s += c
			n = n2
		}
	}

	if s != "" {
		return sensitiveWords.CheckSensitiveWords(s)
	}

	return nil
}

func NewCheckItemComment(v string) (ReviewComment, error) {
	if v == "" {
		return nil, nil
	}

	return newReviewCommentInternal(v, config.MaxLengthOfCheckItemComment)
}

func newReviewComment(v string, maxLen int) (ReviewComment, error) {
	if err := checkReviewComment(v, maxLen); err != nil {
		return nil, err
	}

	if err := sensitiveWords.CheckSensitiveWords(v); err != nil {
		return nil, err
	}

	return reviewComment(v), nil
}

func newReviewCommentInternal(v string, maxLen int) (ReviewComment, error) {
	if err := checkReviewComment(v, maxLen); err != nil {
		return nil, err
	}

	return reviewComment(v), nil
}

func checkReviewComment(v string, maxLen int) error {
	if v == "" {
		return errors.New("empty review comment")
	}

	if utils.StrLen(v) > maxLen {
		return fmt.Errorf(
			"the length of review comment should be less than %d", maxLen,
		)
	}

	return nil
}

// reviewComment
type reviewComment string

func (v reviewComment) ReviewComment() string {
	return string(v)
}
