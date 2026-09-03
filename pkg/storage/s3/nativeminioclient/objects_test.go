package nativeminioclient_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/storage/s3/nativeminioclient"
)

func TestGetDownloadUrl(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		client      *minio.Client
		bucketName  string
		objectKey   string
		expectedUrl string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil client returns error",
			client:      nil,
			bucketName:  "my-bucket",
			objectKey:   "my-object.txt",
			expectError: true,
			errorMsg:    "client",
		},
		{
			name:        "empty bucketName returns error",
			client:      newTestClient(t, "https://s3.example.com"),
			bucketName:  "",
			objectKey:   "my-object.txt",
			expectError: true,
			errorMsg:    "bucketName",
		},
		{
			name:        "empty objectKey returns error",
			client:      newTestClient(t, "https://s3.example.com"),
			bucketName:  "my-bucket",
			objectKey:   "",
			expectError: true,
			errorMsg:    "objectKey",
		},
		{
			name:        "valid inputs return correct URL",
			client:      newTestClient(t, "https://s3.example.com"),
			bucketName:  "my-bucket",
			objectKey:   "my-object.txt",
			expectedUrl: "https://s3.example.com/my-bucket/my-object.txt",
			expectError: false,
		},
		{
			name:        "valid inputs with nested object key",
			client:      newTestClient(t, "https://s3.example.com"),
			bucketName:  "my-bucket",
			objectKey:   "path/to/my-object.txt",
			expectedUrl: "https://s3.example.com/my-bucket/path/to/my-object.txt",
			expectError: false,
		},
		{
			name:        "endpoint with trailing slash",
			client:      newTestClient(t, "https://s3.example.com/"),
			bucketName:  "bucket",
			objectKey:   "file.bin",
			expectedUrl: "https://s3.example.com/bucket/file.bin",
			expectError: false,
		},
		{
			name:        "endpoint with port",
			client:      newTestClient(t, "https://s3.example.com:9000"),
			bucketName:  "my-bucket",
			objectKey:   "data.json",
			expectedUrl: "https://s3.example.com:9000/my-bucket/data.json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := nativeminioclient.GetDownloadUrl(ctx, tt.client, tt.bucketName, tt.objectKey)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUrl, result)
			}
		})
	}
}

func TestGetDownloadUrlByEndpoint(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		endpoint    string
		bucketName  string
		objectKey   string
		expectedUrl string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty endpoint returns error",
			endpoint:    "",
			bucketName:  "my-bucket",
			objectKey:   "my-object.txt",
			expectError: true,
			errorMsg:    "endpoint",
		},
		{
			name:        "empty bucketName returns error",
			endpoint:    "s3.example.com",
			bucketName:  "",
			objectKey:   "my-object.txt",
			expectError: true,
			errorMsg:    "bucketName",
		},
		{
			name:        "empty objectKey returns error",
			endpoint:    "s3.example.com",
			bucketName:  "my-bucket",
			objectKey:   "",
			expectError: true,
			errorMsg:    "objectKey",
		},
		{
			name:        "endpoint without scheme gets https prefix",
			endpoint:    "s3.example.com",
			bucketName:  "my-bucket",
			objectKey:   "my-object.txt",
			expectedUrl: "https://s3.example.com/my-bucket/my-object.txt",
			expectError: false,
		},
		{
			name:        "endpoint with https scheme",
			endpoint:    "https://s3.example.com",
			bucketName:  "my-bucket",
			objectKey:   "my-object.txt",
			expectedUrl: "https://s3.example.com/my-bucket/my-object.txt",
			expectError: false,
		},
		{
			name:        "endpoint with http scheme",
			endpoint:    "http://s3.example.com",
			bucketName:  "my-bucket",
			objectKey:   "my-object.txt",
			expectedUrl: "http://s3.example.com/my-bucket/my-object.txt",
			expectError: false,
		},
		{
			name:        "endpoint with trailing slash",
			endpoint:    "https://s3.example.com/",
			bucketName:  "my-bucket",
			objectKey:   "my-object.txt",
			expectedUrl: "https://s3.example.com/my-bucket/my-object.txt",
			expectError: false,
		},
		{
			name:        "nested object key",
			endpoint:    "s3.example.com",
			bucketName:  "my-bucket",
			objectKey:   "path/to/my-object.txt",
			expectedUrl: "https://s3.example.com/my-bucket/path/to/my-object.txt",
			expectError: false,
		},
		{
			name:        "bucket with leading slash is trimmed",
			endpoint:    "https://s3.example.com",
			bucketName:  "/my-bucket",
			objectKey:   "my-object.txt",
			expectedUrl: "https://s3.example.com/my-bucket/my-object.txt",
			expectError: false,
		},
		{
			name:        "object key with leading slash is trimmed",
			endpoint:    "https://s3.example.com",
			bucketName:  "my-bucket",
			objectKey:   "/my-object.txt",
			expectedUrl: "https://s3.example.com/my-bucket/my-object.txt",
			expectError: false,
		},
		{
			name:        "endpoint with port and no scheme",
			endpoint:    "s3.example.com:9000",
			bucketName:  "my-bucket",
			objectKey:   "file.bin",
			expectedUrl: "https://s3.example.com:9000/my-bucket/file.bin",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := nativeminioclient.GetDownloadUrlByEndpoint(ctx, tt.endpoint, tt.bucketName, tt.objectKey)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedUrl, result)
			}
		})
	}
}

// newTestClient creates a minio.Client with the given endpoint URL for testing purposes.
// This client is not connected to any real server and is only used to test URL construction.
func newTestClient(t *testing.T, endpoint string) *minio.Client {
	t.Helper()

	parsedURL, err := url.Parse(endpoint)
	require.NoError(t, err)

	client, err := minio.New(parsedURL.Host, &minio.Options{
		Secure: parsedURL.Scheme == "https",
	})
	require.NoError(t, err)

	return client
}
