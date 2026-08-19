package httputils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/testwebserver"
)

// Example how to perform a GET request and log the server's certificate information:
// This demonstrates how to collect and inspect SSL/TLS certificates from the server.
func Test_Example_GetRequestAndLogCertificateInfos(t *testing.T) {
	// Preparation start...

	// Enable verbose output
	ctx := contextutils.WithVerbose(context.TODO())

	// Initialize the test web server with TLS:
	const port int = 9124
	testServer, err := testwebserver.GetTestWebServerWithTLS(port)
	require.NoError(t, err)
	defer testServer.Stop(ctx)
	err = testServer.StartInBackground(ctx)
	require.NoError(t, err)
	// ... preparation end.

	// To perform a GET request and collect certificates use:
	response, err := httputils.SendRequest(
		ctx,
		&httpoptions.RequestOptions{
			// Add the URL to request here:
			Url: "https://localhost:9124/hello_world.txt",

			// Enable certificate collection:
			CollectCertificates: true,

			// Skip TLS validation for test purposes (not recommended for production):
			SkipTLSvalidation: true,
		},
	)
	require.NoError(t, err)

	// Log the server's certificate information:
	err = response.LogCertInfo(ctx)
	require.NoError(t, err)

	// You can also access the certificate chain directly:
	certChain, err := response.GetServerCertificateChain(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, certChain)

	// Or get just the end entity (leaf) certificate:
	leafCert, err := response.GetServerEndEntitiyCertificate(ctx)
	require.NoError(t, err)
	require.NotNil(t, leafCert)

	// Now you can inspect certificate properties:
	t.Logf("Certificate Subject: %s", leafCert.Subject.CommonName)
	t.Logf("Certificate Issuer: %s", leafCert.Issuer.CommonName)
	t.Logf("Certificate Valid From: %s", leafCert.NotBefore)
	t.Logf("Certificate Valid Until: %s", leafCert.NotAfter)
}
