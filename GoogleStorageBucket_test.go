package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEST_BUCKET_NAME = "asciich-test-bucket"

func TestGoogleStorageBucketExists(t *testing.T) {

	if ContinuousIntegration().IsRunningInContinuousIntegration() {
		LogInfo("Currently not available in CI/CD.")
		return
	}

	tests := []struct {
		bucketName     string
		expectedExists bool
	}{
		{TEST_BUCKET_NAME, true},
		{TEST_BUCKET_NAME + "-does-not-exist", false},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				var bucket ObjectStoreBucket = MustGetGoogleStorageBucketByName(tt.bucketName)

				assert.EqualValues(tt.expectedExists, bucket.MustExists())
			},
		)
	}
}
