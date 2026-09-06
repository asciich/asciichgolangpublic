package httpoptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
)

func TestDownloadAsTemporaryFileOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := &httpoptions.DownloadAsTemporaryFileOptions{}
		copy := original.GetDeepCopy()

		require.NotNil(t, copy)

		// Modify copy
		copy.Sha256Sum = "abc123"

		// Original should be unchanged
		require.EqualValues(t, "", original.Sha256Sum)
	})

	t.Run("with RequestOptions", func(t *testing.T) {
		original := &httpoptions.DownloadAsTemporaryFileOptions{}
		original.RequestOptions = httpoptions.NewRequestOptions()
		original.RequestOptions.Header = map[string]string{
			"Accept": "application/octet-stream",
		}
		original.RequestOptions.Data = []byte("download request")

		copy := original.GetDeepCopy()

		require.NotNil(t, copy.RequestOptions)
		require.EqualValues(t, original.RequestOptions.Header, copy.RequestOptions.Header)
		require.EqualValues(t, original.RequestOptions.Data, copy.RequestOptions.Data)

		// Modify copy's nested RequestOptions
		copy.RequestOptions.Header["Accept"] = "modified"
		copy.RequestOptions.Data[0] = 'X'

		// Original should be unchanged (deep copy of nested object with map and slice)
		require.EqualValues(t, map[string]string{"Accept": "application/octet-stream"}, original.RequestOptions.Header)
		require.EqualValues(t, []byte("download request"), original.RequestOptions.Data)
	})

	t.Run("with Sha256Sum set", func(t *testing.T) {
		original := &httpoptions.DownloadAsTemporaryFileOptions{
			Sha256Sum: "expected_hash_value",
		}

		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Sha256Sum, copy.Sha256Sum)

		// Modify copy
		copy.Sha256Sum = "modified_hash"

		// Original should be unchanged
		require.EqualValues(t, "expected_hash_value", original.Sha256Sum)
		require.EqualValues(t, "modified_hash", copy.Sha256Sum)
	})
}
