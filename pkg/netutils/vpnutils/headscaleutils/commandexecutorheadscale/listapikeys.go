package commandexecutorheadscale

import (
	"context"
	"encoding/json"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ListApiKeyEntry is one entry of `headscale apikeys list -o json`.
type ListApiKeyEntry struct {
	ID         int               `json:"id"`
	Prefix     string            `json:"prefix"`
	Expiration ListNodeTimestamp `json:"expiration"`
	CreatedAt  ListNodeTimestamp `json:"created_at"`
}

// ToApiKey converts the entry into a headscalegeneric.ApiKey.
func (l *ListApiKeyEntry) ToApiKey() (*headscalegeneric.ApiKey, error) {
	if l == nil {
		return nil, tracederrors.TracedErrorNil("listApiKeyEntry")
	}

	if l.ID <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid api key id '%d'.", l.ID)
	}

	created := l.CreatedAt.ToTime()
	if created == nil {
		return nil, tracederrors.TracedErrorf("Created time of api key with id '%d' is not set.", l.ID)
	}

	ret := headscalegeneric.NewApiKey()

	steps := []error{
		ret.SetID(uint64(l.ID)),
		ret.SetPrefix(l.Prefix),
		ret.SetExpiration(l.Expiration.ToTime()),
		ret.SetCreated(*created),
	}

	for _, err := range steps {
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to convert headscale api key with id '%d': %w", l.ID, err)
		}
	}

	return ret, nil
}

// ListApiKeys returns all headscale API keys with details.
func ListApiKeys(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]headscaleinterfaces.ApiKey, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List headscale api keys on '%s' started.", hostDescription)

	output, err := commandExecutor.RunCommandAndGetStdoutAsBytes(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"headscale", "apikeys", "list", "-ojson"},
	})
	if err != nil {
		return nil, err
	}

	entries := []ListApiKeyEntry{}
	err = json.Unmarshal(output, &entries) // An empty list can be "null" or "[]". Both are fine.
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse JSON of listed headscale api keys: %w", err)
	}

	ret := make([]headscaleinterfaces.ApiKey, 0, len(entries))
	for i := range entries {
		key, err := entries[i].ToApiKey()
		if err != nil {
			return nil, err
		}

		ret = append(ret, key)
	}

	logging.LogInfoByCtxf(ctx, "List headscale api keys on '%s' finished. Found '%d' api keys.", hostDescription, len(ret))

	return ret, nil
}
