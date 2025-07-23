package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

const TEST_BUCKET_NAME = "asciich-test-bucket"

func TestGoogleStorageBucketExists(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		bucketName     string
		expectedExists bool
	}{
		{TEST_BUCKET_NAME, true},
		{TEST_BUCKET_NAME + "-does-not-exist", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				var bucket ObjectStoreBucket = MustGetGoogleStorageBucketByName(tt.bucketName)

				require.EqualValues(tt.expectedExists, bucket.MustExists())
			},
		)
	}
}
