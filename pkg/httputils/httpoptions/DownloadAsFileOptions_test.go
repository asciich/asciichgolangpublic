package httpoptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
)

func TestDownloadAsFileOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := httpoptions.NewDownloadAsFileOptions()
		copy := original.GetDeepCopy()

		require.NotNil(t, copy)

		// Modify copy
		copy.OutputPath = "/modified/path"
		copy.OverwriteExisting = true
		copy.Sha256Sum = "abc123"
		copy.UseSudo = true
		copy.PermissionsString = "0644"

		// Original should be unchanged
		require.EqualValues(t, "", original.OutputPath)
		require.EqualValues(t, false, original.OverwriteExisting)
		require.EqualValues(t, "", original.Sha256Sum)
		require.EqualValues(t, false, original.UseSudo)
		require.EqualValues(t, "", original.PermissionsString)
	})

	t.Run("with RequestOptions", func(t *testing.T) {
		original := httpoptions.NewDownloadAsFileOptions()
		original.RequestOptions = httpoptions.NewRequestOptions()
		original.RequestOptions.Header = map[string]string{
			"Accept": "application/json",
		}
		original.RequestOptions.Data = []byte("request")

		copy := original.GetDeepCopy()

		require.NotNil(t, copy.RequestOptions)
		require.EqualValues(t, original.RequestOptions.Header, copy.RequestOptions.Header)
		require.EqualValues(t, original.RequestOptions.Data, copy.RequestOptions.Data)

		// Modify copy's nested RequestOptions
		copy.RequestOptions.Header["Accept"] = "modified"
		copy.RequestOptions.Data[0] = 'X'

		// Original should be unchanged (deep copy of nested object with map and slice)
		require.EqualValues(t, map[string]string{"Accept": "application/json"}, original.RequestOptions.Header)
		require.EqualValues(t, []byte("request"), original.RequestOptions.Data)
	})

	t.Run("with all fields set", func(t *testing.T) {
		original := httpoptions.NewDownloadAsFileOptions()
		original.RequestOptions = httpoptions.NewRequestOptions()
		original.RequestOptions.Url = "https://example.com/file.txt"
		original.OutputPath = "/tmp/file.txt"
		original.OverwriteExisting = true
		original.Sha256Sum = "abc123"
		original.UseSudo = true
		original.PermissionsString = "0644"

		copy := original.GetDeepCopy()

		// Verify all fields are copied
		require.NotNil(t, copy.RequestOptions)
		require.EqualValues(t, original.OutputPath, copy.OutputPath)
		require.EqualValues(t, original.OverwriteExisting, copy.OverwriteExisting)
		require.EqualValues(t, original.Sha256Sum, copy.Sha256Sum)
		require.EqualValues(t, original.UseSudo, copy.UseSudo)
		require.EqualValues(t, original.PermissionsString, copy.PermissionsString)

		// Modify copy
		copy.OutputPath = "/modified/path"
		copy.RequestOptions.Url = "https://modified.com"

		// Original should be unchanged
		require.EqualValues(t, "/tmp/file.txt", original.OutputPath)
		require.EqualValues(t, "https://example.com/file.txt", original.RequestOptions.Url)
	})
}
