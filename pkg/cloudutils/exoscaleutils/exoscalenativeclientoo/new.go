package exoscalenativeclientoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/cloudutils/exoscaleutils/exoscalenativeclient"
)

func NewNativeClient(apiKey, secretKey string) (*ExoscaleClient, error) {
	nativeClient, err := exoscalenativeclient.NewNativeClient(apiKey, secretKey)
	if err != nil {
		return nil, err
	}

	client := &ExoscaleClient{
		client: nativeClient,
	}

	return client, nil
}

func NewNativeClientFromEnvVars(ctx context.Context) (*ExoscaleClient, error) {
	nativeClient, err := exoscalenativeclient.NewNativeClientFromEnvVars(ctx)
	if err != nil {
		return nil, err
	}

	client := &ExoscaleClient{
		client: nativeClient,
	}

	return client, nil
}
