package genericx509utils_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func generateSelfSignedCertPEM(t *testing.T, subject string) string {
	t.Helper()

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "key.pem")
	certPath := filepath.Join(tmpDir, "cert.pem")

	cmd := exec.Command(
		"openssl", "req", "-x509",
		"-newkey", "rsa:2048",
		"-keyout", keyPath,
		"-out", certPath,
		"-days", "1",
		"-nodes",
		"-subj", subject,
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl command failed: %s", string(output))

	certPEM, err := os.ReadFile(certPath)
	require.NoError(t, err)

	return string(certPEM)
}

func generateCertChainPEM(t *testing.T) (chainPEM string, certCount int) {
	t.Helper()

	tmpDir := t.TempDir()

	// Generate Root CA
	rootKeyPath := filepath.Join(tmpDir, "root_key.pem")
	rootCertPath := filepath.Join(tmpDir, "root_cert.pem")

	cmd := exec.Command(
		"openssl", "req", "-x509",
		"-newkey", "rsa:2048",
		"-keyout", rootKeyPath,
		"-out", rootCertPath,
		"-days", "1",
		"-nodes",
		"-subj", "/C=CH/O=RootOrg/CN=RootCA",
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl root CA generation failed: %s", string(output))

	// Generate Intermediate CA CSR and sign with Root
	intKeyPath := filepath.Join(tmpDir, "int_key.pem")
	intCsrPath := filepath.Join(tmpDir, "int_csr.pem")
	intCertPath := filepath.Join(tmpDir, "int_cert.pem")

	cmd = exec.Command(
		"openssl", "req", "-new",
		"-newkey", "rsa:2048",
		"-keyout", intKeyPath,
		"-out", intCsrPath,
		"-nodes",
		"-subj", "/C=CH/O=IntOrg/CN=IntermediateCA",
	)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "openssl intermediate CSR generation failed: %s", string(output))

	cmd = exec.Command(
		"openssl", "x509", "-req",
		"-in", intCsrPath,
		"-CA", rootCertPath,
		"-CAkey", rootKeyPath,
		"-CAcreateserial",
		"-out", intCertPath,
		"-days", "1",
	)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "openssl intermediate cert signing failed: %s", string(output))

	// Generate End Entity CSR and sign with Intermediate
	eeKeyPath := filepath.Join(tmpDir, "ee_key.pem")
	eeCsrPath := filepath.Join(tmpDir, "ee_csr.pem")
	eeCertPath := filepath.Join(tmpDir, "ee_cert.pem")

	cmd = exec.Command(
		"openssl", "req", "-new",
		"-newkey", "rsa:2048",
		"-keyout", eeKeyPath,
		"-out", eeCsrPath,
		"-nodes",
		"-subj", "/C=CH/O=EEOrg/CN=end-entity.example.com",
	)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "openssl end entity CSR generation failed: %s", string(output))

	cmd = exec.Command(
		"openssl", "x509", "-req",
		"-in", eeCsrPath,
		"-CA", intCertPath,
		"-CAkey", intKeyPath,
		"-CAcreateserial",
		"-out", eeCertPath,
		"-days", "1",
	)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "openssl end entity cert signing failed: %s", string(output))

	// Concatenate: end entity + intermediate + root
	eePEM, err := os.ReadFile(eeCertPath)
	require.NoError(t, err)
	intPEM, err := os.ReadFile(intCertPath)
	require.NoError(t, err)
	rootPEM, err := os.ReadFile(rootCertPath)
	require.NoError(t, err)

	chain := string(eePEM) + string(intPEM) + string(rootPEM)

	return chain, 3
}

func Test_ReadCertFromString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		cert, err := genericx509utils.ReadCertFromString("")
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("invalid PEM", func(t *testing.T) {
		cert, err := genericx509utils.ReadCertFromString("not a valid PEM")
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("valid certificate", func(t *testing.T) {
		certPEM := generateSelfSignedCertPEM(t, "/C=CH/L=Zurich/O=TestOrg/CN=TestCert")

		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)
		require.NotNil(t, cert)
		require.Equal(t, "TestOrg", cert.Subject.Organization[0])
		require.Equal(t, "CH", cert.Subject.Country[0])
		require.Equal(t, "Zurich", cert.Subject.Locality[0])
		require.Equal(t, "TestCert", cert.Subject.CommonName)
	})
}

func Test_ReadCertFromBytes(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		cert, err := genericx509utils.ReadCertFromBytes(nil)
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("invalid PEM bytes", func(t *testing.T) {
		cert, err := genericx509utils.ReadCertFromBytes([]byte("garbage data"))
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("valid certificate", func(t *testing.T) {
		certPEM := generateSelfSignedCertPEM(t, "/C=DE/O=ByteOrg/CN=ByteCert")

		cert, err := genericx509utils.ReadCertFromBytes([]byte(certPEM))
		require.NoError(t, err)
		require.NotNil(t, cert)
		require.Equal(t, "ByteOrg", cert.Subject.Organization[0])
		require.Equal(t, "DE", cert.Subject.Country[0])
		require.Equal(t, "ByteCert", cert.Subject.CommonName)
	})
}

func Test_ReadCertsFromString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		certs, err := genericx509utils.ReadCertsFromString("")
		require.Error(t, err)
		require.Nil(t, certs)
	})

	t.Run("invalid PEM", func(t *testing.T) {
		certs, err := genericx509utils.ReadCertsFromString("not a valid PEM")
		require.Error(t, err)
		require.Nil(t, certs)
	})

	t.Run("single certificate", func(t *testing.T) {
		certPEM := generateSelfSignedCertPEM(t, "/C=CH/O=SingleOrg/CN=SingleCert")

		certs, err := genericx509utils.ReadCertsFromString(certPEM)
		require.NoError(t, err)
		require.Len(t, certs, 1)
		require.Equal(t, "SingleOrg", certs[0].Subject.Organization[0])
		require.Equal(t, "SingleCert", certs[0].Subject.CommonName)
	})

	t.Run("certificate chain", func(t *testing.T) {
		chainPEM, expectedCount := generateCertChainPEM(t)

		certs, err := genericx509utils.ReadCertsFromString(chainPEM)
		require.NoError(t, err)
		require.Len(t, certs, expectedCount)
		require.Equal(t, "end-entity.example.com", certs[0].Subject.CommonName)
		require.Equal(t, "IntermediateCA", certs[1].Subject.CommonName)
		require.Equal(t, "RootCA", certs[2].Subject.CommonName)
	})
}

func Test_ReadCertsFromBytes(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		certs, err := genericx509utils.ReadCertsFromBytes(nil)
		require.Error(t, err)
		require.Nil(t, certs)
	})

	t.Run("invalid PEM bytes", func(t *testing.T) {
		certs, err := genericx509utils.ReadCertsFromBytes([]byte("garbage data"))
		require.Error(t, err)
		require.Nil(t, certs)
	})

	t.Run("certificate chain", func(t *testing.T) {
		chainPEM, expectedCount := generateCertChainPEM(t)

		certs, err := genericx509utils.ReadCertsFromBytes([]byte(chainPEM))
		require.NoError(t, err)
		require.Len(t, certs, expectedCount)
		require.Equal(t, "EEOrg", certs[0].Subject.Organization[0])
		require.Equal(t, "IntOrg", certs[1].Subject.Organization[0])
		require.Equal(t, "RootOrg", certs[2].Subject.Organization[0])
	})
}

func Test_ReadCertFromString_and_ReadCertFromBytes_consistency(t *testing.T) {
	t.Run("string and bytes return same result", func(t *testing.T) {
		certPEM := generateSelfSignedCertPEM(t, "/C=CH/O=ConsistencyOrg/CN=ConsistencyCert")

		certFromString, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		certFromBytes, err := genericx509utils.ReadCertFromBytes([]byte(certPEM))
		require.NoError(t, err)

		require.True(t, certFromString.Equal(certFromBytes))
	})
}

func Test_ReadCertsFromString_and_ReadCertsFromBytes_consistency(t *testing.T) {
	t.Run("string and bytes return same results", func(t *testing.T) {
		chainPEM, _ := generateCertChainPEM(t)

		certsFromString, err := genericx509utils.ReadCertsFromString(chainPEM)
		require.NoError(t, err)

		certsFromBytes, err := genericx509utils.ReadCertsFromBytes([]byte(chainPEM))
		require.NoError(t, err)

		require.Len(t, certsFromString, len(certsFromBytes))
		for i := range certsFromString {
			require.True(t, certsFromString[i].Equal(certsFromBytes[i]))
		}
	})
}
