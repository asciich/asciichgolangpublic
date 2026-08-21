package nativetruststoreoo_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils/nativetruststoreoo"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils/truststoreinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
)

func TestNativeTrustStore_InterfaceCompliance(t *testing.T) {
	var _ truststoreinterfaces.TrustStore = (*nativetruststoreoo.NativeTrustStore)(nil)
}

func TestNewNativeTrustStore(t *testing.T) {
	t.Run("create with temp dir", func(t *testing.T) {
		ctx := context.Background()
		tempDirPath, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer os.RemoveAll(tempDirPath)

		trustStore, err := nativetruststoreoo.NewNativeTrustStore(tempDirPath)
		require.NoError(t, err)
		require.NotNil(t, trustStore)
	})

	t.Run("create with default OS trust store (empty path)", func(t *testing.T) {
		trustStore, err := nativetruststoreoo.NewNativeTrustStore("")
		require.NoError(t, err)
		require.NotNil(t, trustStore)
	})
}

func TestNativeTrustStore_ListCaCertificates(t *testing.T) {
	t.Run("list certificates from custom truststore path", func(t *testing.T) {
		ctx := context.Background()
		tempDirPath, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer os.RemoveAll(tempDirPath)

		certName := "Test Native CA 1"
		options := &x509options.X509CreateCertificateOptions{
			CommonName:     certName,
			Organization:   "Test Org",
			CountryName:    "CH",
			Locality:       "Zurich",
			PrivateKeySize: 2048,
		}

		certKeyPair, err := x509utils.CreateRootCaCertificate(ctx, options)
		require.NoError(t, err)

		certPEMBytes, err := genericx509utils.WriteCertAsBytes(certKeyPair.Cert)
		require.NoError(t, err)

		certPath := filepath.Join(tempDirPath, "test-ca.crt")
		err = os.WriteFile(certPath, certPEMBytes, 0644)
		require.NoError(t, err)

		trustStore, err := nativetruststoreoo.NewNativeTrustStore(tempDirPath)
		require.NoError(t, err)

		certs, err := trustStore.ListCaCertificates(ctx)
		require.NoError(t, err)
		require.Len(t, certs, 1)

		loadedCert := certs[0]
		require.NotNil(t, loadedCert)
		require.True(t, certKeyPair.Cert.Equal(loadedCert))
		require.Equal(t, certName, loadedCert.Subject.CommonName)
		require.Equal(t, "Test Org", loadedCert.Subject.Organization[0])
		require.Equal(t, "CH", loadedCert.Subject.Country[0])
		require.Equal(t, "Zurich", loadedCert.Subject.Locality[0])
		require.True(t, loadedCert.IsCA)
	})

	t.Run("empty directory returns empty list without error", func(t *testing.T) {
		ctx := context.Background()
		tempDirPath, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer os.RemoveAll(tempDirPath)

		trustStore, err := nativetruststoreoo.NewNativeTrustStore(tempDirPath)
		require.NoError(t, err)

		certs, err := trustStore.ListCaCertificates(ctx)
		require.NoError(t, err)
		require.Empty(t, certs)
	})
}

func TestNativeTrustStore_ContextCancellation(t *testing.T) {
	t.Run("ListCaCertificates respects context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		tempDirPath, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer os.RemoveAll(tempDirPath)

		cancel()

		trustStore, err := nativetruststoreoo.NewNativeTrustStore(tempDirPath)
		require.NoError(t, err)

		certs, err := trustStore.ListCaCertificates(ctx)
		require.ErrorIs(t, err, context.Canceled)
		require.Nil(t, certs)
	})
}
