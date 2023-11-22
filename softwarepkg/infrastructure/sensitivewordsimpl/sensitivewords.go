package sensitivewordsimpl

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/moderation/v3/model"

	"github.com/opensourceways/software-package-server/allerror"
)

var instance *service

type service struct {
	cli *v3.ModerationClient
}

func Init(cfg *Config) error {
	auth := basic.NewCredentialsBuilder().
		WithAk(cfg.AccessKey).
		WithSk(cfg.SecretKey).
		WithIamEndpointOverride(cfg.IAMEndpint).
		Build()

	cli := v3.NewModerationClient(
		v3.ModerationClientBuilder().
			WithRegion(region.NewRegion(cfg.Region, cfg.Endpoint)).
			WithCredential(auth).
			Build(),
	)

	instance = &service{cli: cli}

	return nil
}

func Sensitive() *service {
	return instance
}

func (s *service) CheckSensitiveWords(content string) error {
	request := &model.RunTextModerationRequest{
		Body: &model.TextDetectionReq{
			Data: &model.TextDetectionDataReq{
				Text: content,
			},
			EventType: "comment",
		},
	}

	resp, err := s.cli.RunTextModeration(request)
	if err != nil {
		return err
	}

	if *resp.Result.Suggestion != "pass" {
		return allerror.New(allerror.ErrorCodeSensitiveContent, "invalid text")
	}

	return nil
}
