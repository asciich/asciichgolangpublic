package httputils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/mustutils"
)

func Test_TestWebServer_SetAndGetCertificate(t *testing.T) {
	const port int = 9123

	testServer := NewTestWebServer()
	testServer.MustSetPort(port)

	cert, key := mustutils.Must2(generateCertAndKeyForTestWebserver())

	mustutils.Must0(testServer.SetTlsCertAndKey(cert, key))

	cert2 := mustutils.Must(testServer.GetTlsCert())

	require.True(t, cert.Equal(cert2))
}