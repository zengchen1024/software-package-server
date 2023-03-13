package dp

import "errors"

type Language interface {
	Language() string
}

func NewLanguage(v string) (Language, error) {
	if v == "" || !config.isValidLanguage(v) {
		return nil, errors.New("invalid language")
	}

	return language(v), nil
}

type language string

func (v language) Language() string {
	return string(v)
}
