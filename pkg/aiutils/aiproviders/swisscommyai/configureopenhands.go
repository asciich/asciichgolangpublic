package swisscommyai

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/aiutils/openhandsutils"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

const API_KEY_ENV_VAR_NAME = "MYAI_API_KEY"

func AddLlmProfileToOpenhands(ctx context.Context, openHandsUrl string) error {
	err := urlsutils.CheckIsUrl(openHandsUrl)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add LLM profile for myai of Swisscom to openhands '%s' started.", openHandsUrl)

	apiKey, err := environmentvariables.GetEnvValueAsString(ctx, API_KEY_ENV_VAR_NAME, false)
	if err != nil {
		return err
	}

	openHands, err := openhandsutils.NewOpenHands(openHandsUrl)
	if err != nil {
		return err
	}

	err = openHands.CreateLlmProfile(ctx, "myai", &openhandsutils.LlmProfileConfig{
		Model:   "openai/qwen3.5-397b-a17b",
		ApiKey:  apiKey,
		BaseUrl: "https://code.myai.swisscom.ch/v1",
	})

	logging.LogInfoByCtxf(ctx, "Add LLM profile for myai of Swisscom to openhands '%s' finished.", openHandsUrl)

	return nil
}
