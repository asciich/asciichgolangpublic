package httpoptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
)

func TestRequestOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := httpoptions.NewRequestOptions()
		copy := original.GetDeepCopy()

		require.NotNil(t, copy)

		// Modify copy
		copy.Url = "https://modified.com"
		copy.Path = "/modified"
		copy.Port = 9999
		copy.Method = "POST"
		copy.SkipTLSvalidation = true
		copy.Header = map[string]string{"modified": "header"}
		copy.Data = []byte("modified")

		// Original should be unchanged
		require.EqualValues(t, "", original.Url)
		require.EqualValues(t, "", original.Path)
		require.EqualValues(t, 0, original.Port)
		require.EqualValues(t, "", original.Method)
		require.EqualValues(t, false, original.SkipTLSvalidation)
		require.Nil(t, original.Header)
		require.Nil(t, original.Data)
	})

	t.Run("with Header map", func(t *testing.T) {
		original := httpoptions.NewRequestOptions()
		original.Header = map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token",
		}

		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Header, copy.Header)

		// Modify copy's map
		copy.Header["Content-Type"] = "modified"
		copy.Header["X-Custom"] = "new-header"

		// Original should be unchanged (deep copy of map)
		require.EqualValues(t, map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token",
		}, original.Header)
		require.EqualValues(t, map[string]string{
			"Content-Type":  "modified",
			"Authorization": "Bearer token",
			"X-Custom":      "new-header",
		}, copy.Header)
	})

	t.Run("with Data byte slice", func(t *testing.T) {
		original := httpoptions.NewRequestOptions()
		original.Data = []byte("request body")

		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Data, copy.Data)

		// Modify copy's byte slice
		copy.Data[0] = 'X'
		copy.Data = append(copy.Data, byte('!'))

		// Original should be unchanged (deep copy of byte slice)
		require.EqualValues(t, []byte("request body"), original.Data)
		require.NotEqualValues(t, original.Data, copy.Data)
	})

	t.Run("with all fields set", func(t *testing.T) {
		original := httpoptions.NewRequestOptions()
		original.Url = "https://example.com"
		original.Path = "/api/v1"
		original.Port = 8443
		original.Method = "POST"
		original.SkipTLSvalidation = true
		original.Header = map[string]string{
			"Content-Type": "application/json",
		}
		original.Data = []byte("payload")

		copy := original.GetDeepCopy()

		// Verify all fields are copied
		require.EqualValues(t, original.Url, copy.Url)
		require.EqualValues(t, original.Path, copy.Path)
		require.EqualValues(t, original.Port, copy.Port)
		require.EqualValues(t, original.Method, copy.Method)
		require.EqualValues(t, original.SkipTLSvalidation, copy.SkipTLSvalidation)
		require.EqualValues(t, original.Header, copy.Header)
		require.EqualValues(t, original.Data, copy.Data)

		// Modify copy
		copy.Url = "https://modified.com"
		copy.Header["Content-Type"] = "modified"
		copy.Data[0] = 'X'

		// Original should be unchanged
		require.EqualValues(t, "https://example.com", original.Url)
		require.EqualValues(t, map[string]string{"Content-Type": "application/json"}, original.Header)
		require.EqualValues(t, []byte("payload"), original.Data)
	})
}
