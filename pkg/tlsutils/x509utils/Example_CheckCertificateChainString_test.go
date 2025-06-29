package x509utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils"
)

func Test_CheckCertificateChainString(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// ---
	// Preparation start
	//
	// Generate the rootCA
	rootCaCertAndKey, err := x509utils.CreateRootCa(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg root",
			PrivateKeySize: 1024,
		},
	)
	require.NoError(t, err)
	rootPem, err := x509utils.EncodeCertificateAsPEMString(rootCaCertAndKey.Cert)
	require.NoError(t, err)

	// Generate the intermediate CA
	intermediateCertAndKey, err := x509utils.CreateSignedIntermediateCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg intermediate",
			PrivateKeySize: 1024,
		},
		rootCaCertAndKey,
	)
	require.NoError(t, err)
	intermediatePem, err := x509utils.EncodeCertificateAsPEMString(intermediateCertAndKey.Cert)
	require.NoError(t, err)

	endEndityCertAndKey, err := x509utils.CreateSignedEndEndityCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CommonName:     "mytestcn.example.net",
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg endEndity",
			AdditionalSans: []string{"mytestsan1.example.net", "mytestsan2.example.net"},
			PrivateKeySize: 1024,
		},
		intermediateCertAndKey,
	)
	require.NoError(t, err)
	certPem, err := x509utils.EncodeCertificateAsPEMString(endEndityCertAndKey.Cert)
	require.NoError(t, err)
	// Preparation end
	// ---

	// An empty string is not valid:
	err = x509utils.CheckCertificateChainString(ctx, "")
	require.Error(t, err)

	// only the end endity certificate is not valid
	err = x509utils.CheckCertificateChainString(ctx, certPem)
	require.Error(t, err)

	// a missing root certificate is also not valid
	err = x509utils.CheckCertificateChainString(ctx, certPem+"\n"+intermediatePem)
	require.Error(t, err)

	// a complete chain is valud
	err = x509utils.CheckCertificateChainString(ctx, certPem+"\n"+intermediatePem+"\n"+rootPem)
	require.NoError(t, err)
}
