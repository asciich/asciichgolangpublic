package httputils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_TestWebServer_SetAndGetCertificate(t *testing.T) {
	ctx := getCtx()
	const port int = 9123

	testServer := httputils.NewTestWebServer()
	err := testServer.SetPort(port)
	require.NoError(t, err)

	certAndKey := mustutils.Must(httputils.GenerateCertAndKeyForTestWebserver(getCtx()))

	mustutils.Must0(testServer.SetTlsCertAndKey(ctx, certAndKey))

	cert2 := mustutils.Must(testServer.GetTlsCert())

	require.True(t, certAndKey.Cert.Equal(cert2))
}
