package infomaniak

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/aiutils/openhandsutils"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

const API_KEY_ENV_VAR_NAME = "INFOMANIAK_API_KEY"

func AddLlmProfileToOpenhands(ctx context.Context, openHandsUrl string, productId string) error {
	if productId == "" {
		return tracederrors.TracedErrorEmptyString("productId")
	}

	err := AddLlmProfile397bToOpenhands(ctx, openHandsUrl, productId)
	if err != nil {
		return err
	}

	err = AddLlmProfile122bToOpenhands(ctx, openHandsUrl, productId)
	if err != nil {
		return err
	}

	return nil
}

func AddLlmProfile397bToOpenhands(ctx context.Context, openHandsUrl string, productId string) error {
	if productId == "" {
		return tracederrors.TracedErrorEmptyString("productId")
	}

	err := urlsutils.CheckIsUrl(openHandsUrl)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add LLM profile for infomaniak-397b to openhands '%s' started.", openHandsUrl)

	apiKey, err := environmentvariables.GetEnvValueAsString(ctx, API_KEY_ENV_VAR_NAME, false)
	if err != nil {
		return err
	}

	openHands, err := openhandsutils.NewOpenHands(openHandsUrl)
	if err != nil {
		return err
	}

	err = openHands.CreateLlmProfile(ctx, "infomaniak-397b", &openhandsutils.LlmProfileConfig{
		Model:   "openai/Qwen/Qwen3.5-397B-A17B-FP8",
		ApiKey:  apiKey,
		BaseUrl: fmt.Sprintf("https://api.infomaniak.com/2/ai/%s/openai/v1", productId),
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add LLM profile for infomaniak-397b to openhands '%s' finished.", openHandsUrl)

	return nil
}

func AddLlmProfile122bToOpenhands(ctx context.Context, openHandsUrl string, productId string) error {
	if productId == "" {
		return tracederrors.TracedErrorEmptyString("productId")
	}

	err := urlsutils.CheckIsUrl(openHandsUrl)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add LLM profile for infomaniak-122b to openhands '%s' started.", openHandsUrl)

	apiKey, err := environmentvariables.GetEnvValueAsString(ctx, API_KEY_ENV_VAR_NAME, false)
	if err != nil {
		return err
	}

	openHands, err := openhandsutils.NewOpenHands(openHandsUrl)
	if err != nil {
		return err
	}

	err = openHands.CreateLlmProfile(ctx, "infomaniak-122b", &openhandsutils.LlmProfileConfig{
		Model:   "openai/Qwen/Qwen3.5-122B-A10B-FP8",
		ApiKey:  apiKey,
		BaseUrl: fmt.Sprintf("https://api.infomaniak.com/2/ai/%s/openai/v1", productId),
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Add LLM profile for infomaniak-122b to openhands '%s' finished.", openHandsUrl)

	return nil
}
