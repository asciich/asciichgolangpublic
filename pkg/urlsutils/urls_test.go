package urlsutils_test

import (
	"strconv"
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
		{"http://example.com", ""}, // no trailing slash -> empty path
		{"ftp://host/some/path/file.txt", "/some/path/file.txt"},
		{"//example.com/relative", "/relative"}, // scheme-less absolute URL (Parse handles it)
		{"/just/a/path", "/just/a/path"},        // relative path only
	}

	for _, tt := range tests {
		got, err := urlsutils.GetPath(tt.in)
		require.NoError(t, err)
		require.Equal(t, tt.want, got, "input: %q", tt.in)
	}
}

func TestSetPort(t *testing.T) {
	tests := []struct {
		url      string
		port     int
		expected string
	}{
		{"https://example.com", 123, "https://example.com:123"},
		{"https://example.com", 443, "https://example.com:443"},
		{"https://example.com/index.html", 123, "https://example.com:123/index.html"},
		{"https://example.com/index.html", 443, "https://example.com:443/index.html"},
		{"https://example.com/path/index.html", 123, "https://example.com:123/path/index.html"},
		{"https://example.com/path/index.html", 443, "https://example.com:443/path/index.html"},
		{"https://example.com:80", 123, "https://example.com:123"},
		{"https://example.com:80", 443, "https://example.com:443"},
		{"https://example.com:80/index.html", 123, "https://example.com:123/index.html"},
		{"https://example.com:80/index.html", 443, "https://example.com:443/index.html"},
		{"https://example.com:80/path/index.html", 123, "https://example.com:123/path/index.html"},
		{"https://example.com:80/path/index.html", 443, "https://example.com:443/path/index.html"},
		{"https://example.com:443", 123, "https://example.com:123"},
		{"https://example.com:443", 443, "https://example.com:443"},
		{"https://example.com:443/index.html", 123, "https://example.com:123/index.html"},
		{"https://example.com:443/index.html", 443, "https://example.com:443/index.html"},
		{"https://example.com:443/path/index.html", 123, "https://example.com:123/path/index.html"},
		{"https://example.com:443/path/index.html", 443, "https://example.com:443/path/index.html"},
		{"https://example.com:1234", 123, "https://example.com:123"},
		{"https://example.com:1234", 443, "https://example.com:443"},
		{"https://example.com:1234/index.html", 123, "https://example.com:123/index.html"},
		{"https://example.com:1234/index.html", 443, "https://example.com:443/index.html"},
		{"https://example.com:1234/path/index.html", 123, "https://example.com:123/path/index.html"},
		{"https://example.com:1234/path/index.html", 443, "https://example.com:443/path/index.html"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got, err := urlsutils.SetPort(tt.url, tt.port)
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, got)
		})
	}

	t.Run("empty URL", func(t *testing.T) {
		got, err := urlsutils.SetPort("", 123)
		require.Error(t, err)
		require.EqualValues(t, "", got)
	})

	for _, port := range []int{0, -1, -80, -443} {
		t.Run("invalid Port "+strconv.Itoa(port), func(t *testing.T) {
			got, err := urlsutils.SetPort("https://example.com", port)
			require.Error(t, err)
			require.EqualValues(t, "", got)
		})
	}
}



func TestSetPath(t *testing.T) {
	tests := []struct {
		url      string
		path     string
		expected string
	}{
		{"https://example.com", "", "https://example.com"},
		{"https://example.com/", "", "https://example.com"},
		{"https://example.com/this", "", "https://example.com"},
		{"https://example.com/this/", "", "https://example.com"},
		{"https://example.com/this/will", "", "https://example.com"},
		{"https://example.com/this/will/", "", "https://example.com"},
		{"https://example.com/this/will/be", "", "https://example.com"},
		{"https://example.com/this/will/be/", "", "https://example.com"},
		{"https://example.com/this/will/be/removed", "", "https://example.com"},
		{"https://example.com", "/", "https://example.com/"},
		{"https://example.com:443", "", "https://example.com:443"},
		{"https://example.com:443", "/", "https://example.com:443/"},
		{"https://example.com/this/will/be/overwritten", "by/this", "https://example.com/by/this"},
		{"https://example.com/this/will/be/overwritten", "/by/this", "https://example.com/by/this"},
		{"https://example.com:8443/this/will/be/overwritten", "by/this", "https://example.com:8443/by/this"},
		{"https://example.com:8443/this/will/be/overwritten", "/by/this", "https://example.com:8443/by/this"},
	}
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got, err := urlsutils.SetPath(tt.url, tt.path)
			require.NoError(t, err)
			require.EqualValues(t, tt.expected, got)
		})
	}

	t.Run("empty URL", func(t *testing.T) {
		got, err := urlsutils.SetPath("", "")
		require.Error(t, err)
		require.EqualValues(t, "", got)
	})
}
