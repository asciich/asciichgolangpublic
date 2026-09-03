package nativeminioclient

import (
	"context"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/asciich/asciichgolangpublic/pkg/environmentvariables"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/s3options"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const ENV_VAR_DEFAULT_ACCESS_KEY = "AWS_ACCESS_KEY"
const ENV_VAR_DEFAULT_SECRET_ACCESS_KEY = "AWS_SECRET_ACCESS_KEY"

func NewClientFromEnvVars(ctx context.Context, endpoint string, options *s3options.NewS3ClientOptions) (*minio.Client, error) {
	if options == nil {
		options = &s3options.NewS3ClientOptions{}
	}

	accessKey, err := environmentvariables.GetEnvValueAsString(ctx, ENV_VAR_DEFAULT_ACCESS_KEY, false)
	if err != nil {
		return nil, err
	}

	secretKey, err := environmentvariables.GetEnvValueAsString(ctx, ENV_VAR_DEFAULT_SECRET_ACCESS_KEY, false)
	if err != nil {
		return nil, err
	}

	return NewClient(endpoint, accessKey, secretKey, options)
}

func NewClient(endoint string, accessKey string, secretKey string, options *s3options.NewS3ClientOptions) (*minio.Client, error) {
	if endoint == "" {
		return nil, tracederrors.TracedErrorEmptyString("endpoint")
	}

	if options == nil {
		options = &s3options.NewS3ClientOptions{}
	}

	client, err := minio.New(endoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: options.UseTLS,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create minio client for '%s': %w", endoint, err)
	}

	return client, err
}

// NewAdminClient creates a new madmin admin client for server management operations
func NewAdminClient(endpoint string, accessKey string, secretKey string, options *s3options.NewS3ClientOptions) (*madmin.AdminClient, error) {
	if endpoint == "" {
		return nil, tracederrors.TracedErrorEmptyString("endpoint")
	}

	if options == nil {
		options = &s3options.NewS3ClientOptions{}
	}

	adminClient, err := madmin.NewWithOptions(endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: options.UseTLS,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create madmin admin client for '%s': %w", endpoint, err)
	}

	return adminClient, nil
}

// NewAdminClientFromEnvVars creates a new admin client using credentials from environment variables
func NewAdminClientFromEnvVars(ctx context.Context, endpoint string, options *s3options.NewS3ClientOptions) (*madmin.AdminClient, error) {
	accessKey, err := environmentvariables.GetEnvValueAsString(ctx, ENV_VAR_DEFAULT_ACCESS_KEY, false)
	if err != nil {
		return nil, err
	}

	secretKey, err := environmentvariables.GetEnvValueAsString(ctx, ENV_VAR_DEFAULT_SECRET_ACCESS_KEY, false)
	if err != nil {
		return nil, err
	}

	return NewAdminClient(endpoint, accessKey, secretKey, options)
}
