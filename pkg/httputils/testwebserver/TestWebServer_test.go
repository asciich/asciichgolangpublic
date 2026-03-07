package testwebserver_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
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

	certAndKey, err := testwebserver.GenerateCertAndKeyForTestWebserver(getCtx())
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
