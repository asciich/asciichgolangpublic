package openhandsutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CreateDeleteAndListLlmProfile(t *testing.T) {
	ctx := getCtx()
	openHands := startTestOpenHands(ctx, t)

	const profileName = "testai"

	err := openHands.DeleteAllLlmProfiles(ctx)
	require.NoError(t, err)

	exists, err := openHands.LlmProfileExists(ctx, profileName)
	require.NoError(t, err)
	require.False(t, exists)

	profileNames, err := openHands.ListLlmProfileNames(ctx)
	require.NoError(t, err)
	require.Empty(t, profileNames)

	activeName, err := openHands.GetActiveLlmProfileName(ctx)
	require.NoError(t, err)
	require.Empty(t, activeName)

	config := &LlmProfileConfig{
		Model:   "openai/qwen3.5-397b-a17b",
		ApiKey:  "mysecretapikey",
		BaseUrl: "https://code.myai.swisscom.ch/v1",
	}

	for range 2 { // run twice to check idempotence
		err = openHands.CreateLlmProfile(ctx, profileName, config)
		require.NoError(t, err)

		profileNames, err = openHands.ListLlmProfileNames(ctx)
		require.NoError(t, err)
		require.EqualValues(t, []string{profileName}, profileNames)

		exists, err = openHands.LlmProfileExists(ctx, profileName)
		require.NoError(t, err)
		require.True(t, exists)
	}

	err = openHands.ActivateLlmProfile(ctx, profileName)
	require.NoError(t, err)

	activeName, err = openHands.GetActiveLlmProfileName(ctx)
	require.NoError(t, err)
	require.EqualValues(t, profileName, activeName)
}
