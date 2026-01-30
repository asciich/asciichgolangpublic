package exoscalenativeclient

import (
	"context"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/credentials"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const ENV_VAR_EXOSCALE_API_KEY = "EXOSCALE_API_KEY"
const ENV_VAR_EXOSCALE_API_SECRET = "EXOSCALE_API_SECRET"

func NewNativeClientFromEnvVars(ctx context.Context) (*v3.Client, error) {
	apiKey, err := environmentvariables.GetEnvValueAsString(ctx, ENV_VAR_EXOSCALE_API_KEY, false)
	if err != nil {
		return nil, err
	}

	secretKey, err := environmentvariables.GetEnvValueAsString(ctx, ENV_VAR_EXOSCALE_API_SECRET, false)
	if err != nil {
		return nil, err
	}

	return NewNativeClient(apiKey, secretKey)
}

func NewNativeClient(apiKey, secretKey string) (*v3.Client, error) {
	if apiKey == "" {
		return nil, tracederrors.TracedErrorEmptyString("apiKey")
	}

	if secretKey == "" {
		return nil, tracederrors.TracedErrorEmptyString("secretKey")
	}

	creds := credentials.NewStaticCredentials(apiKey, secretKey)

	client, err := v3.NewClient(creds)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create native exoscale client: %w", err)
	}

	return client, nil
}
