package httputils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/contextutils"
	"github.com/asciich/asciichgolangpublic/mustutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_TestWebServer_SetAndGetCertificate(t *testing.T) {
	ctx := getCtx()
	const port int = 9123

	testServer := NewTestWebServer()
	testServer.MustSetPort(port)

	certAndKey := mustutils.Must(generateCertAndKeyForTestWebserver(getCtx()))

	mustutils.Must0(testServer.SetTlsCertAndKey(ctx, certAndKey))

	cert2 := mustutils.Must(testServer.GetTlsCert())

	require.True(t, certAndKey.Cert.Equal(cert2))
}
