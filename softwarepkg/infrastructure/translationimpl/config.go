package translationimpl

import (
	"errors"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/model"
)

type Config struct {
	AccessKey string `json:"access_key"     required:"true"`
	SecretKey string `json:"secret_key"     required:"true"`
	Project   string `json:"project"        required:"true"`
	Region    string `json:"region"         required:"true"`
	Endpoint  string `json:"endpoint"       required:"true"`
}

func getSupportedLanguage(languages []string) (map[string]model.TextTranslationReqTo, error) {
	v := supportedLanguages()

	for _, s := range languages {
		if _, ok := v[s]; !ok {
			return nil, errors.New("unsupported language: " + s)
		}
	}

	return v, nil
}

func supportedLanguages() map[string]model.TextTranslationReqTo {
	t := model.GetTextTranslationReqToEnum()

	return map[string]model.TextTranslationReqTo{
		"chinese": t.ZH,
		"english": t.EN,
		// it can add more here.
	}
}
