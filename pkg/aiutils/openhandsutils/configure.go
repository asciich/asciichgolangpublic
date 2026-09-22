package openhandsutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

const (
	OPENAI_BASE_URL_ENV_VAR = "OPENAI_BASE_URL"
	OPENAI_API_KEY_ENV_VAR  = "OPENAI_API_KEY"
)

type ConfigureLlmProfileOptions struct {
	ProfileName  string
	Model        string
	BaseUrl      string
	ApiKey       string
	ReadFromEnv  bool
}

func ConfigureLlmProfile(ctx context.Context, openHandsUrl string, options *ConfigureLlmProfileOptions) error {
	err := urlsutils.CheckIsUrl(openHandsUrl)
	if err != nil {
		return err
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.ProfileName == "" {
		return tracederrors.TracedErrorEmptyString("options.ProfileName")
	}

	if options.Model == "" {
		return tracederrors.TracedErrorEmptyString("options.Model")
	}

	baseUrl := options.BaseUrl
	apiKey := options.ApiKey

	if options.ReadFromEnv {
		if baseUrl == "" {
			baseUrl, err = environmentvariables.GetEnvValueAsString(ctx, OPENAI_BASE_URL_ENV_VAR, false)
			if err != nil {
				return err
			}
		}

		if apiKey == "" {
			apiKey, err = environmentvariables.GetEnvValueAsString(ctx, OPENAI_API_KEY_ENV_VAR, false)
			if err != nil {
				return err
			}
		}
	}

	if baseUrl == "" {
		return tracederrors.TracedErrorEmptyString("baseUrl or " + OPENAI_BASE_URL_ENV_VAR)
	}

	if apiKey == "" {
		return tracederrors.TracedErrorEmptyString("apiKey or " + OPENAI_API_KEY_ENV_VAR)
	}

	err = urlsutils.CheckIsUrl(baseUrl)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Configure generic LLM profile '%s' on openhands '%s' started.", options.ProfileName, openHandsUrl)

	openHands, err := NewOpenHands(openHandsUrl)
	if err != nil {
		return err
	}

	err = openHands.CreateLlmProfile(ctx, options.ProfileName, &LlmProfileConfig{
		Model:   options.Model,
		ApiKey:  apiKey,
		BaseUrl: baseUrl,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Configure generic LLM profile '%s' on openhands '%s' finished.", options.ProfileName, openHandsUrl)

	return nil
}
