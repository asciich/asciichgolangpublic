package nativex509utils_test

import (
	"crypto/x509"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/nativex509utils"
)

func generateTestCertPEM(t *testing.T, subject string) string {
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

func generateTestCertChainPEMs(t *testing.T) (rootPEM string, intPEM string, eePEM string) {
	t.Helper()

	tmpDir := t.TempDir()

	// Root CA
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
		"-addext", "basicConstraints=critical,CA:TRUE",
		"-addext", "keyUsage=critical,keyCertSign,cRLSign",
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl root CA generation failed: %s", string(output))

	// Intermediate
	intKeyPath := filepath.Join(tmpDir, "int_key.pem")
	intCsrPath := filepath.Join(tmpDir, "int_csr.pem")
	intCertPath := filepath.Join(tmpDir, "int_cert.pem")
	intExtPath := filepath.Join(tmpDir, "int_ext.cnf")

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

	err = os.WriteFile(intExtPath, []byte("basicConstraints=critical,CA:TRUE,pathlen:0\nkeyUsage=critical,keyCertSign,cRLSign\n"), 0644)
	require.NoError(t, err)

	cmd = exec.Command(
		"openssl", "x509", "-req",
		"-in", intCsrPath,
		"-CA", rootCertPath,
		"-CAkey", rootKeyPath,
		"-CAcreateserial",
		"-out", intCertPath,
		"-days", "1",
		"-extfile", intExtPath,
	)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "openssl intermediate cert signing failed: %s", string(output))

	// End Entity
	eeKeyPath := filepath.Join(tmpDir, "ee_key.pem")
	eeCsrPath := filepath.Join(tmpDir, "ee_csr.pem")
	eeCertPath := filepath.Join(tmpDir, "ee_cert.pem")
	eeExtPath := filepath.Join(tmpDir, "ee_ext.cnf")

	cmd = exec.Command(
		"openssl", "req", "-new",
		"-newkey", "rsa:2048",
		"-keyout", eeKeyPath,
		"-out", eeCsrPath,
		"-nodes",
		"-subj", "/C=CH/O=EEOrg/CN=server.example.com",
	)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "openssl end entity CSR generation failed: %s", string(output))

	err = os.WriteFile(eeExtPath, []byte("basicConstraints=critical,CA:FALSE\nkeyUsage=critical,digitalSignature,keyEncipherment\n"), 0644)
	require.NoError(t, err)

	cmd = exec.Command(
		"openssl", "x509", "-req",
		"-in", eeCsrPath,
		"-CA", intCertPath,
		"-CAkey", intKeyPath,
		"-CAcreateserial",
		"-out", eeCertPath,
		"-days", "1",
		"-extfile", eeExtPath,
	)
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "openssl end entity cert signing failed: %s", string(output))

	rootPEMBytes, err := os.ReadFile(rootCertPath)
	require.NoError(t, err)
	intPEMBytes, err := os.ReadFile(intCertPath)
	require.NoError(t, err)
	eePEMBytes, err := os.ReadFile(eeCertPath)
	require.NoError(t, err)

	return string(rootPEMBytes), string(intPEMBytes), string(eePEMBytes)
}

func Test_WriteCertAsString(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		output, err := nativex509utils.WriteCertAsString(nil)
		require.Error(t, err)
		require.Empty(t, output)
	})

	t.Run("valid cert roundtrip", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=WriteOrg/CN=WriteCert")

		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		output, err := nativex509utils.WriteCertAsString(cert)
		require.NoError(t, err)
		require.NotEmpty(t, output)

		certReadBack, err := genericx509utils.ReadCertFromString(output)
		require.NoError(t, err)
		require.True(t, cert.Equal(certReadBack))
	})
}

func Test_WriteCertAsBytes(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		output, err := nativex509utils.WriteCertAsBytes(nil)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("valid cert roundtrip", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=DE/O=ByteWriteOrg/CN=ByteWriteCert")

		cert, err := genericx509utils.ReadCertFromBytes([]byte(certPEM))
		require.NoError(t, err)

		output, err := nativex509utils.WriteCertAsBytes(cert)
		require.NoError(t, err)
		require.NotNil(t, output)

		certReadBack, err := genericx509utils.ReadCertFromBytes(output)
		require.NoError(t, err)
		require.True(t, cert.Equal(certReadBack))
	})
}

func Test_WriteCertToFile(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		err := nativex509utils.WriteCertToFile(nil, "/tmp/test.pem")
		require.Error(t, err)
	})

	t.Run("empty file path", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=FileOrg/CN=FileCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		err = nativex509utils.WriteCertToFile(cert, "")
		require.Error(t, err)
	})

	t.Run("valid cert write and read back", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=FileOrg/CN=FileCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "output_cert.pem")

		err = nativex509utils.WriteCertToFile(cert, filePath)
		require.NoError(t, err)

		fileContent, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.NotEmpty(t, fileContent)

		certReadBack, err := genericx509utils.ReadCertFromBytes(fileContent)
		require.NoError(t, err)
		require.True(t, cert.Equal(certReadBack))
	})

	t.Run("invalid directory", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=TestOrg/CN=TestCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		err = nativex509utils.WriteCertToFile(cert, "/non/existent/dir/cert.pem")
		require.Error(t, err)
	})
}

func Test_WriteCertsAsString(t *testing.T) {
	t.Run("nil certs", func(t *testing.T) {
		output, err := nativex509utils.WriteCertsAsString(nil)
		require.Error(t, err)
		require.Empty(t, output)
	})

	t.Run("empty slice", func(t *testing.T) {
		output, err := nativex509utils.WriteCertsAsString([]*x509.Certificate{})
		require.Error(t, err)
		require.Empty(t, output)
	})

	t.Run("nil cert in slice", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=TestOrg/CN=TestCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		output, err := nativex509utils.WriteCertsAsString([]*x509.Certificate{cert, nil})
		require.Error(t, err)
		require.Empty(t, output)
	})

	t.Run("single cert roundtrip", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=SingleOrg/CN=SingleCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		output, err := nativex509utils.WriteCertsAsString([]*x509.Certificate{cert})
		require.NoError(t, err)
		require.NotEmpty(t, output)

		certsReadBack, err := genericx509utils.ReadCertsFromString(output)
		require.NoError(t, err)
		require.Len(t, certsReadBack, 1)
		require.True(t, cert.Equal(certsReadBack[0]))
	})

	t.Run("certificate chain roundtrip", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateTestCertChainPEMs(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		certs := []*x509.Certificate{eeCert, intCert, rootCert}

		output, err := nativex509utils.WriteCertsAsString(certs)
		require.NoError(t, err)
		require.NotEmpty(t, output)

		certsReadBack, err := genericx509utils.ReadCertsFromString(output)
		require.NoError(t, err)
		require.Len(t, certsReadBack, 3)
		require.True(t, eeCert.Equal(certsReadBack[0]))
		require.True(t, intCert.Equal(certsReadBack[1]))
		require.True(t, rootCert.Equal(certsReadBack[2]))
	})
}

func Test_WriteCertsAsBytes(t *testing.T) {
	t.Run("nil certs", func(t *testing.T) {
		output, err := nativex509utils.WriteCertsAsBytes(nil)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("empty slice", func(t *testing.T) {
		output, err := nativex509utils.WriteCertsAsBytes([]*x509.Certificate{})
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("certificate chain roundtrip", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateTestCertChainPEMs(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		certs := []*x509.Certificate{eeCert, intCert, rootCert}

		output, err := nativex509utils.WriteCertsAsBytes(certs)
		require.NoError(t, err)
		require.NotNil(t, output)

		certsReadBack, err := genericx509utils.ReadCertsFromBytes(output)
		require.NoError(t, err)
		require.Len(t, certsReadBack, 3)

		for i := range certs {
			require.True(t, certs[i].Equal(certsReadBack[i]))
		}
	})
}

func Test_WriteCertsToFile(t *testing.T) {
	t.Run("nil certs", func(t *testing.T) {
		err := nativex509utils.WriteCertsToFile(nil, "/tmp/test.pem")
		require.Error(t, err)
	})

	t.Run("empty slice", func(t *testing.T) {
		err := nativex509utils.WriteCertsToFile([]*x509.Certificate{}, "/tmp/test.pem")
		require.Error(t, err)
	})

	t.Run("empty file path", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=TestOrg/CN=TestCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		err = nativex509utils.WriteCertsToFile([]*x509.Certificate{cert}, "")
		require.Error(t, err)
	})

	t.Run("certificate chain write and read back", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateTestCertChainPEMs(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		certs := []*x509.Certificate{eeCert, intCert, rootCert}

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "chain.pem")

		err = nativex509utils.WriteCertsToFile(certs, filePath)
		require.NoError(t, err)

		fileContent, err := os.ReadFile(filePath)
		require.NoError(t, err)
		require.NotEmpty(t, fileContent)

		certsReadBack, err := genericx509utils.ReadCertsFromBytes(fileContent)
		require.NoError(t, err)
		require.Len(t, certsReadBack, 3)

		for i := range certs {
			require.True(t, certs[i].Equal(certsReadBack[i]))
		}
	})

	t.Run("invalid directory", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=TestOrg/CN=TestCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		err = nativex509utils.WriteCertsToFile([]*x509.Certificate{cert}, "/non/existent/dir/chain.pem")
		require.Error(t, err)
	})
}

func Test_WriteCertAsString_and_WriteCertAsBytes_consistency(t *testing.T) {
	t.Run("string and bytes produce same output", func(t *testing.T) {
		certPEM := generateTestCertPEM(t, "/C=CH/O=ConsistencyOrg/CN=ConsistencyCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		outputString, err := nativex509utils.WriteCertAsString(cert)
		require.NoError(t, err)

		outputBytes, err := nativex509utils.WriteCertAsBytes(cert)
		require.NoError(t, err)

		require.Equal(t, outputString, string(outputBytes))
	})
}
