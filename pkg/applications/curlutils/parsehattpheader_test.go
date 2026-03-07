package curlutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/applications/curlutils"
)

func Test_ParseHttpHeader(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		parsed, err := curlutils.ParseHttpHeader([]byte(""))
		require.Error(t, err)
		require.Nil(t, parsed)
	})

	t.Run("200 OK only newline", func(t *testing.T) {
		content := []byte(`HTTP/1.1 200 OK
Date: Sat, 07 Mar 2026 18:43:44 GMT
Content-Length: 722
Content-Type: text/html; charset=utf-8
`)
		parsed, err := curlutils.ParseHttpHeader(content)
		require.NoError(t, err)
		require.EqualValues(t, 200, parsed.StatusCode)
	})

	t.Run("200 OK", func(t *testing.T) {
		content := []byte("HTTP/1.1 200 OK\r\nDate: Sat, 07 Mar 2026 18:43:44 GMT\r\nContent-Length: 722\r\nContent-Type: text/html; charset=utf-8\r\n")
		parsed, err := curlutils.ParseHttpHeader(content)
		require.NoError(t, err)
		require.EqualValues(t, 200, parsed.StatusCode)
	})

	t.Run("400 Not Found", func(t *testing.T) {
		content := []byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/plain; charset=utf-8\r\nX-Content-Type-Options: nosniff\r\nDate: Sat, 07 Mar 2026 19:02:51 GMT\r\nContent-Length: 19\r\n\r\n")

		parsed, err := curlutils.ParseHttpHeader(content)
		require.NoError(t, err)
		require.EqualValues(t, 404, parsed.StatusCode)
	})
}
