package googleaistudio

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/aiutils/openhandsutils"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

const API_KEY_ENV_VAR_NAME = "GOOGLE_AI_STUDIO_API_KEY"

func AddLlmProfileToOpenhands(ctx context.Context, openHandsUrl string) error {
	err := urlsutils.CheckIsUrl(openHandsUrl)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add LLM profile for Google AI Studio to openhands '%s' started.", openHandsUrl)

	apiKey, err := environmentvariables.GetEnvValueAsString(ctx, API_KEY_ENV_VAR_NAME, false)
	if err != nil {
		return err
	}

	openHands, err := openhandsutils.NewOpenHands(openHandsUrl)
	if err != nil {
		return err
	}

	err = openHands.CreateLlmProfile(ctx, "google-ai-studio", &openhandsutils.LlmProfileConfig{
		Model:   "gemini/gemini-3.5-flash",
		ApiKey:  apiKey,
		BaseUrl: "https://generativelanguage.googleapis.com/v1beta",
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add LLM profile for Google AI Studio to openhands '%s' finished.", openHandsUrl)

	return nil
}
