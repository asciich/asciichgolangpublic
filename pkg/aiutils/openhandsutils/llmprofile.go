package openhandsutils

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type LlmProfileConfig struct {
	Model   string `json:"model"`
	ApiKey  string `json:"api_key"`
	BaseUrl string `json:"base_url"`
}

func (o *Openhands) DeleteAllLlmProfiles(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Delete all openhands LLM profiles started.")

	type Profile struct {
		Name      string `json:"name"`
		Model     string `json:"model"`
		BaseUrl   string `json:"base_url"`
		ApiKeySet bool   `json:"api_key_set"`
	}

	type ProfilesResponse struct {
		Profiles      []Profile `json:"profiles"`
		ActiveProfile *string   `json:"active_profile"`
	}

	body, err := o.DoApiRequest(ctx, "GET", "/api/profiles", nil)
	if err != nil {
		return err
	}

	var response ProfilesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to parse profiles response: %w", err)
	}

	for _, profile := range response.Profiles {
		logging.LogInfoByCtxf(ctx, "Deleting openhands LLM profile '%s'.", profile.Name)

		_, err := o.DoApiRequest(
			ctx,
			"DELETE",
			"/api/profiles/"+profile.Name,
			nil,
		)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete LLM profile '%s': %w", profile.Name, err)
		}

		logging.LogChangedByCtxf(ctx, "Deleted LLM profile '%s'.", profile.Name)
	}

	logging.LogInfoByCtxf(ctx, "Delete all openhands LLM profiles finished.")

	return nil
}

func (o *Openhands) ListLlmProfileNames(ctx context.Context) ([]string, error) {
	logging.LogInfoByCtxf(ctx, "List openhands LLM profile names started.")

	type Profile struct {
		Name      string `json:"name"`
		Model     string `json:"model"`
		BaseUrl   string `json:"base_url"`
		ApiKeySet bool   `json:"api_key_set"`
	}

	type ProfilesResponse struct {
		Profiles      []Profile `json:"profiles"`
		ActiveProfile *string   `json:"active_profile"`
	}

	body, err := o.DoApiRequest(
		ctx,
		"GET",
		"/api/profiles",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var response ProfilesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to parse profiles response: %w", err)
	}

	names := []string{}
	for _, profile := range response.Profiles {
		names = append(names, profile.Name)
	}

	sort.Strings(names)

	logging.LogInfoByCtxf(ctx, "List openhands LLM profile names finished.")

	return names, nil
}

func (o *Openhands) LlmProfileExists(ctx context.Context, profileName string) (bool, error) {
	if profileName == "" {
		return false, tracederrors.TracedErrorEmptyString("profileName")
	}

	existingNames, err := o.ListLlmProfileNames(ctx)
	if err != nil {
		return false, err
	}

	for _, name := range existingNames {
		if name == profileName {
			return true, nil
		}
	}

	return false, nil
}

func (o *Openhands) CreateLlmProfile(ctx context.Context, profileName string, config *LlmProfileConfig) error {
	if profileName == "" {
		return tracederrors.TracedErrorEmptyString("profileName")
	}

	if config == nil {
		return tracederrors.TracedErrorNil("config")
	}

	logging.LogInfoByCtxf(ctx, "Create openhands LLM profile '%s' started.", profileName)

	exists, err := o.LlmProfileExists(ctx, profileName)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Openhands LLM profile '%s' already exists.", profileName)
		return nil
	}

	type LlmPayload struct {
		Model              string  `json:"model"`
		ApiKey             string  `json:"api_key"`
		BaseUrl            string  `json:"base_url"`
		AuthType           string  `json:"auth_type"`
		SubscriptionVendor *string `json:"subscription_vendor"`
	}

	type ProfilePayload struct {
		Llm            LlmPayload `json:"llm"`
		IncludeSecrets bool       `json:"include_secrets"`
	}

	payloadData := ProfilePayload{
		Llm: LlmPayload{
			Model:              config.Model,
			ApiKey:             config.ApiKey,
			BaseUrl:            config.BaseUrl,
			AuthType:           "api_key",
			SubscriptionVendor: nil,
		},
		IncludeSecrets: true,
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to prepare payload to create LLM profile: %w", err)
	}

	_, err = o.DoApiRequest(ctx, "POST", "/api/profiles/"+profileName, payload)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Created LLM profile '%s'.", profileName)

	logging.LogInfoByCtxf(ctx, "Create openhands LLM profile '%s' finished.", profileName)

	return nil
}

func (o *Openhands) ActivateLlmProfile(ctx context.Context, profileName string) error {
	if profileName == "" {
		return tracederrors.TracedErrorEmptyString("profileName")
	}

	logging.LogInfoByCtxf(ctx, "Activate openhands LLM profile '%s' started.", profileName)

	exists, err := o.LlmProfileExists(ctx, profileName)
	if err != nil {
		return err
	}

	if !exists {
		return tracederrors.TracedErrorf("LLM profile '%s' does not exist", profileName)
	}

	_, err = o.DoApiRequest(ctx, "POST", "/api/profiles/"+profileName+"/activate", []byte("{}"))
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Activated LLM profile '%s'.", profileName)

	logging.LogInfoByCtxf(ctx, "Activate openhands LLM profile '%s' finished.", profileName)

	return nil
}

func (o *Openhands) GetActiveLlmProfileName(ctx context.Context) (string, error) {
	logging.LogInfoByCtxf(ctx, "Get active openhands LLM profile name started.")

	type Profile struct {
		Name      string `json:"name"`
		Model     string `json:"model"`
		BaseUrl   string `json:"base_url"`
		ApiKeySet bool   `json:"api_key_set"`
	}

	type ProfilesResponse struct {
		Profiles      []Profile `json:"profiles"`
		ActiveProfile *string   `json:"active_profile"`
	}

	body, err := o.DoApiRequest(ctx, "GET", "/api/profiles", nil)
	if err != nil {
		return "", err
	}

	var response ProfilesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to parse profiles response: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Get active openhands LLM profile name finished.")

	if response.ActiveProfile == nil {
		return "", nil
	}

	return *response.ActiveProfile, nil
}
