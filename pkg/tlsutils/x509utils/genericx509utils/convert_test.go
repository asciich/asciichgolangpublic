package genericx509utils_test

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func Test_TlsCertToX509Cert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil tls cert returns error", func(t *testing.T) {
		cert, err := genericx509utils.TlsCertToX509Cert(nil)
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("empty certificate data returns error", func(t *testing.T) {
		tlsCert := &tls.Certificate{}

		cert, err := genericx509utils.TlsCertToX509Cert(tlsCert)
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("valid tls cert returns x509 cert", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		certPEM, err := genericx509utils.WriteCertAsBytes(rootPair.Cert)
		require.NoError(t, err)

		keyPEM, err := rootPair.GetPrivateKeyAsPEMString()
		require.NoError(t, err)

		tlsCert, err := tls.X509KeyPair(certPEM, []byte(keyPEM))
		require.NoError(t, err)

		x509Cert, err := genericx509utils.TlsCertToX509Cert(&tlsCert)
		require.NoError(t, err)
		require.NotNil(t, x509Cert)
		require.True(t, rootPair.Cert.Equal(x509Cert))
	})

	t.Run("returned cert has correct subject fields", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		certPEM, err := genericx509utils.WriteCertAsBytes(rootPair.Cert)
		require.NoError(t, err)

		keyPEM, err := rootPair.GetPrivateKeyAsPEMString()
		require.NoError(t, err)

		tlsCert, err := tls.X509KeyPair(certPEM, []byte(keyPEM))
		require.NoError(t, err)

		x509Cert, err := genericx509utils.TlsCertToX509Cert(&tlsCert)
		require.NoError(t, err)

		require.Equal(t, "CH", x509Cert.Subject.Country[0])
		require.Equal(t, "Zurich", x509Cert.Subject.Locality[0])
		require.Equal(t, "TestRootOrg", x509Cert.Subject.Organization[0])
		require.Equal(t, "TestRootCA", x509Cert.Subject.CommonName)
	})
}
