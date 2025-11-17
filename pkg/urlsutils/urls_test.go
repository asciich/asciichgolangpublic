package urlsutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/urlsutils"
)

func Test_GetBaseUrl(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		url, err := urlsutils.GetBaseUrl("")
		require.Error(t, err)
		require.Empty(t, url)
	})

	t.Run("already base url", func(t *testing.T) {
		url, err := urlsutils.GetBaseUrl("https://asciich.ch")
		require.NoError(t, err)
		require.EqualValues(t, "https://asciich.ch", url)
	})

	t.Run("already base url with slash", func(t *testing.T) {
		url, err := urlsutils.GetBaseUrl("https://asciich.ch/")
		require.NoError(t, err)
		require.EqualValues(t, "https://asciich.ch", url)
	})

	t.Run("with path", func(t *testing.T) {
		url, err := urlsutils.GetBaseUrl("https://asciich.ch/a/random/path.txt")
		require.NoError(t, err)
		require.EqualValues(t, "https://asciich.ch", url)
	})
}
