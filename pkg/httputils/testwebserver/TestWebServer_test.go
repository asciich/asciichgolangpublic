package testwebserver_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_TestWebServer_SetAndGetCertificate(t *testing.T) {
	ctx := getCtx()
	const port int = 9123

	testServer := testwebserver.NewTestWebServer()
	err := testServer.SetPort(port)
	require.NoError(t, err)

	certAndKey, err := testwebserver.GenerateCertAndKeyForTestWebserver(ctx)
	require.NoError(t, err)

	err = testServer.SetTlsCertAndKey(ctx, certAndKey)
	require.NoError(t, err)

	cert2, err := testServer.GetTlsCert()
	require.NoError(t, err)

	require.True(t, certAndKey.Cert.Equal(cert2))
}

func Test_TestWebsServer_GetUrl(t *testing.T) {
	testWebServer, err := testwebserver.GetTestWebServer(1234)
	require.NoError(t, err)

	url, err := testWebServer.GetUrl()
	require.NoError(t, err)
	require.EqualValues(t, "http://localhost:1234", url)
}

func Test_TestWebServer_ArchiveEndpoints(t *testing.T) {
	ctx := getCtx()
	const port int = 9130

	testWebserver, err := testwebserver.GetTestWebServer(port)
	require.NoError(t, err)
	err = testWebserver.StartInBackground(ctx)
	require.NoError(t, err)
	defer testWebserver.Stop(ctx)

	t.Run("hello_world.tar", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9130/hello_world.tar")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "application/x-tar", resp.Header.Get("Content-Type"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// Verify it's a valid tar archive containing hello_world.txt
		content, err := tarutils.ReadFileFromTarArchiveBytesAsBytes(body, "hello_world.txt")
		require.NoError(t, err)
		require.Equal(t, "hello world\n", string(content))
	})

	t.Run("hello_world.tar.gz", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9130/hello_world.tar.gz")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "application/gzip", resp.Header.Get("Content-Type"))

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Greater(t, len(body), 0)

		// Write to temp file to test extraction
		tmpFile := t.TempDir() + "/test.tar.gz"
		err = writeFile(tmpFile, body)
		require.NoError(t, err)

		// Verify it's a valid tar.gz archive containing hello_world.txt
		content, err := tarutils.ReadFileFromTarArchiveAsBytes(ctx, tmpFile, "hello_world.txt")
		require.NoError(t, err)
		require.Equal(t, "hello world\n", string(content))
	})
}

func writeFile(path string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	return err
}
