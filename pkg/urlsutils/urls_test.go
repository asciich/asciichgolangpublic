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


func TestGetPath(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"https://example.com/foo/bar?query=1", "/foo/bar"},
		{"http://example.com/", "/"},
		{"http://example.com", ""},               // no trailing slash -> empty path
		{"ftp://host/some/path/file.txt", "/some/path/file.txt"},
		{"//example.com/relative", "/relative"},  // scheme-less absolute URL (Parse handles it)
		{"/just/a/path", "/just/a/path"},         // relative path only
	}

	for _, tc := range tests {
		got, err := urlsutils.GetPath(tc.in)
		require.NoError(t, err)
		require.Equal(t, tc.want, got, "input: %q", tc.in)
	}
}