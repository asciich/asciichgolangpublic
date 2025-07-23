package x509utils_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/tlsutils/x509utils"
)

func Test_GenerateSelfSignedCertificateAndEncodeAsPem(t *testing.T) {
	// Get the default context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Generate a self signed certificate and key:
	certKeyPair, err := x509utils.CreateSelfSignedCertificate(ctx, &x509utils.X509CreateCertificateOptions{
		CommonName: "example.example.net",
		Organization: "Example org",
		Locality: "Zurich",
		CountryName: "CH",
		PrivateKeySize: 1024,
	})
	require.NoError(t, err)

	// Get the certificate
	cert, err := certKeyPair.GetX509Certificate()
	require.NoError(t, err)

	// Encode generated certificate as PEM:
	pem, err := x509utils.EncodeCertificateAsPEMString(cert)
	require.NoError(t, err)
	
	// pem does now contain the generated certificate as PEM string:
	require.True(t, strings.HasPrefix(pem, "-----BEGIN CERTIFICATE-----"))
}
