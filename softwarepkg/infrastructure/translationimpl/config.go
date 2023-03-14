package translationimpl

import "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nlp/v2/model"

type Config struct {
	AccessKey    string `json:"access_key"     required:"true"`
	SecretKey    string `json:"secret_key"     required:"true"`
	Project      string `json:"project"        required:"true"`
	Region       string `json:"region"         required:"true"`
	AuthEndpoint string `json:"auth_endpoint"  required:"true"`
}

var textTranslationTo map[string]model.TextTranslationReqTo

func initMap(languages []string) {
	textTranslationTo = make(map[string]model.TextTranslationReqTo, len(languages))
	t := model.GetTextTranslationReqToEnum()

	for _, language := range languages {
		if l, ok := reqTo(language, t); ok {
			textTranslationTo[language] = l
		}
	}
}

func reqTo(l string, t model.TextTranslationReqToEnum) (model.TextTranslationReqTo, bool) {
	switch l {
	case "chinese":
		return t.ZH, true
	case "english":
		return t.EN, true
	}

	return model.TextTranslationReqTo{}, false
}
