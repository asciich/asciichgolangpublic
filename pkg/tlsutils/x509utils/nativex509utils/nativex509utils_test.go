package nativex509utils_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	native "github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/nativex509utils"
)

func generateSelfSignedCert(t *testing.T, dir string) string {
	t.Helper()

	keyPath := filepath.Join(dir, "key.pem")
	certPath := filepath.Join(dir, "cert.pem")

	cmd := exec.Command(
		"openssl", "req", "-x509",
		"-newkey", "rsa:2048",
		"-keyout", keyPath,
		"-out", certPath,
		"-days", "1",
		"-nodes",
		"-subj", "/C=CH/L=Zurich/O=TestOrg/CN=TestCert",
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl command failed: %s", string(output))

	return certPath
}

func Test_ReadCertificateFromFile(t *testing.T) {
	t.Run("nil path", func(t *testing.T) {
		cert, err := native.ReadCertificateFromFile(contextutils.ContextVerbose(), "")
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("non-existent file", func(t *testing.T) {
		cert, err := native.ReadCertificateFromFile(contextutils.ContextVerbose(), "/non/existent/path.pem")
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("valid certificate file", func(t *testing.T) {
		tmpDir := t.TempDir()
		certPath := generateSelfSignedCert(t, tmpDir)

		cert, err := native.ReadCertificateFromFile(contextutils.ContextVerbose(), certPath)
		require.NoError(t, err)
		require.NotNil(t, cert)
		require.Equal(t, "TestOrg", cert.Subject.Organization[0])
		require.Equal(t, "CH", cert.Subject.Country[0])
		require.Equal(t, "Zurich", cert.Subject.Locality[0])
		require.Equal(t, "TestCert", cert.Subject.CommonName)
	})

	t.Run("invalid certificate file", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidPath := filepath.Join(tmpDir, "invalid.pem")
		err := os.WriteFile(invalidPath, []byte("not a valid certificate"), 0644)
		require.NoError(t, err)

		cert, err := native.ReadCertificateFromFile(contextutils.ContextVerbose(), invalidPath)
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("certificate is parseable and consistent", func(t *testing.T) {
		tmpDir := t.TempDir()
		certPath := generateSelfSignedCert(t, tmpDir)

		ctx := contextutils.ContextVerbose()

		cert1, err := native.ReadCertificateFromFile(ctx, certPath)
		require.NoError(t, err)

		cert2, err := native.ReadCertificateFromFile(ctx, certPath)
		require.NoError(t, err)

		require.True(t, cert1.Equal(cert2), "Reading the same certificate file twice should yield equal certificates")
	})
}
