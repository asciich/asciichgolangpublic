package genericx509utils_test

import (
	"crypto/x509"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func generateWriteTestCertPEM(t *testing.T, subject string) string {
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

func generateWriteTestCertChainPEMs(t *testing.T) (rootPEM string, intPEM string, eePEM string) {
	t.Helper()

	tmpDir := t.TempDir()

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
		output, err := genericx509utils.WriteCertAsString(nil)
		require.Error(t, err)
		require.Empty(t, output)
	})

	t.Run("valid cert roundtrip", func(t *testing.T) {
		certPEM := generateWriteTestCertPEM(t, "/C=CH/O=WriteOrg/CN=WriteCert")

		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		output, err := genericx509utils.WriteCertAsString(cert)
		require.NoError(t, err)
		require.NotEmpty(t, output)
		require.Contains(t, output, "BEGIN CERTIFICATE")
		require.Contains(t, output, "END CERTIFICATE")

		certReadBack, err := genericx509utils.ReadCertFromString(output)
		require.NoError(t, err)
		require.True(t, cert.Equal(certReadBack))
	})

	t.Run("preserves subject fields", func(t *testing.T) {
		certPEM := generateWriteTestCertPEM(t, "/C=DE/L=Berlin/O=SubjectOrg/CN=SubjectCert")

		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		output, err := genericx509utils.WriteCertAsString(cert)
		require.NoError(t, err)

		certReadBack, err := genericx509utils.ReadCertFromString(output)
		require.NoError(t, err)
		require.Equal(t, "DE", certReadBack.Subject.Country[0])
		require.Equal(t, "Berlin", certReadBack.Subject.Locality[0])
		require.Equal(t, "SubjectOrg", certReadBack.Subject.Organization[0])
		require.Equal(t, "SubjectCert", certReadBack.Subject.CommonName)
	})
}

func Test_WriteCertAsBytes(t *testing.T) {
	t.Run("nil cert", func(t *testing.T) {
		output, err := genericx509utils.WriteCertAsBytes(nil)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("valid cert roundtrip", func(t *testing.T) {
		certPEM := generateWriteTestCertPEM(t, "/C=CH/O=ByteOrg/CN=ByteCert")

		cert, err := genericx509utils.ReadCertFromBytes([]byte(certPEM))
		require.NoError(t, err)

		output, err := genericx509utils.WriteCertAsBytes(cert)
		require.NoError(t, err)
		require.NotNil(t, output)

		certReadBack, err := genericx509utils.ReadCertFromBytes(output)
		require.NoError(t, err)
		require.True(t, cert.Equal(certReadBack))
	})
}

func Test_WriteCertsAsString(t *testing.T) {
	t.Run("nil certs", func(t *testing.T) {
		output, err := genericx509utils.WriteCertsAsString(nil)
		require.Error(t, err)
		require.Empty(t, output)
	})

	t.Run("empty slice", func(t *testing.T) {
		output, err := genericx509utils.WriteCertsAsString([]*x509.Certificate{})
		require.Error(t, err)
		require.Empty(t, output)
	})

	t.Run("nil cert in slice", func(t *testing.T) {
		certPEM := generateWriteTestCertPEM(t, "/C=CH/O=TestOrg/CN=TestCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		output, err := genericx509utils.WriteCertsAsString([]*x509.Certificate{cert, nil})
		require.Error(t, err)
		require.Empty(t, output)
	})

	t.Run("single cert roundtrip", func(t *testing.T) {
		certPEM := generateWriteTestCertPEM(t, "/C=CH/O=SingleOrg/CN=SingleCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		output, err := genericx509utils.WriteCertsAsString([]*x509.Certificate{cert})
		require.NoError(t, err)
		require.NotEmpty(t, output)

		certsReadBack, err := genericx509utils.ReadCertsFromString(output)
		require.NoError(t, err)
		require.Len(t, certsReadBack, 1)
		require.True(t, cert.Equal(certsReadBack[0]))
	})

	t.Run("certificate chain roundtrip", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateWriteTestCertChainPEMs(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		certs := []*x509.Certificate{eeCert, intCert, rootCert}

		output, err := genericx509utils.WriteCertsAsString(certs)
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
		output, err := genericx509utils.WriteCertsAsBytes(nil)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("empty slice", func(t *testing.T) {
		output, err := genericx509utils.WriteCertsAsBytes([]*x509.Certificate{})
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("certificate chain roundtrip", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateWriteTestCertChainPEMs(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		certs := []*x509.Certificate{eeCert, intCert, rootCert}

		output, err := genericx509utils.WriteCertsAsBytes(certs)
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

func Test_WriteCertAsString_and_WriteCertAsBytes_consistency(t *testing.T) {
	t.Run("string and bytes produce same output", func(t *testing.T) {
		certPEM := generateWriteTestCertPEM(t, "/C=CH/O=ConsistencyOrg/CN=ConsistencyCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		outputString, err := genericx509utils.WriteCertAsString(cert)
		require.NoError(t, err)

		outputBytes, err := genericx509utils.WriteCertAsBytes(cert)
		require.NoError(t, err)

		require.Equal(t, outputString, string(outputBytes))
	})
}

func Test_WriteCertsAsString_and_WriteCertsAsBytes_consistency(t *testing.T) {
	t.Run("string and bytes produce same output", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateWriteTestCertChainPEMs(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		certs := []*x509.Certificate{eeCert, intCert, rootCert}

		outputString, err := genericx509utils.WriteCertsAsString(certs)
		require.NoError(t, err)

		outputBytes, err := genericx509utils.WriteCertsAsBytes(certs)
		require.NoError(t, err)

		require.Equal(t, outputString, string(outputBytes))
	})
}
