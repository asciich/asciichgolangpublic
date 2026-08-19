package genericx509utils_test

import (
	"crypto/x509"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
)

func generateChainCerts(t *testing.T) (rootPEM string, intPEM string, eePEM string) {
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

	// Intermediate CA
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

func generateUnrelatedSelfSignedCertPEM(t *testing.T, subject string) string {
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
		"-addext", "basicConstraints=critical,CA:TRUE",
		"-addext", "keyUsage=critical,keyCertSign,cRLSign",
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl command failed: %s", string(output))

	certPEM, err := os.ReadFile(certPath)
	require.NoError(t, err)

	return string(certPEM)
}

func Test_IsSignedBy(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil certToCheck", func(t *testing.T) {
		rootPEM := generateUnrelatedSelfSignedCertPEM(t, "/C=CH/O=TestOrg/CN=Root")
		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, nil, rootCert)
		require.Error(t, err)
		require.False(t, isSigned)
	})

	t.Run("nil signingCert", func(t *testing.T) {
		rootPEM := generateUnrelatedSelfSignedCertPEM(t, "/C=CH/O=TestOrg/CN=Root")
		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, rootCert, nil)
		require.Error(t, err)
		require.False(t, isSigned)
	})

	t.Run("self signed cert is signed by itself", func(t *testing.T) {
		rootPEM := generateUnrelatedSelfSignedCertPEM(t, "/C=CH/O=TestOrg/CN=SelfSigned")
		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, rootCert, rootCert)
		require.NoError(t, err)
		require.True(t, isSigned)
	})

	t.Run("intermediate is signed by root", func(t *testing.T) {
		rootPEM, intPEM, _ := generateChainCerts(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, intCert, rootCert)
		require.NoError(t, err)
		require.True(t, isSigned)
	})

	t.Run("end entity is signed by intermediate", func(t *testing.T) {
		_, intPEM, eePEM := generateChainCerts(t)

		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, eeCert, intCert)
		require.NoError(t, err)
		require.True(t, isSigned)
	})

	t.Run("end entity is not signed by root directly", func(t *testing.T) {
		rootPEM, _, eePEM := generateChainCerts(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, eeCert, rootCert)
		require.NoError(t, err)
		require.False(t, isSigned)
	})

	t.Run("unrelated cert is not signed by another", func(t *testing.T) {
		rootPEM, _, _ := generateChainCerts(t)
		unrelatedPEM := generateUnrelatedSelfSignedCertPEM(t, "/C=DE/O=OtherOrg/CN=OtherRoot")

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		unrelatedCert, err := genericx509utils.ReadCertFromString(unrelatedPEM)
		require.NoError(t, err)

		isSigned, err := genericx509utils.IsSignedBy(ctx, unrelatedCert, rootCert)
		require.NoError(t, err)
		require.False(t, isSigned)
	})
}

func Test_IsCertChain(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil chain", func(t *testing.T) {
		isChain, err := genericx509utils.IsCertChain(ctx, nil)
		require.Error(t, err)
		require.False(t, isChain)
	})

	t.Run("empty chain", func(t *testing.T) {
		isChain, err := genericx509utils.IsCertChain(ctx, []*x509.Certificate{})
		require.NoError(t, err)
		require.False(t, isChain)
	})

	t.Run("single cert is not a chain", func(t *testing.T) {
		rootPEM := generateUnrelatedSelfSignedCertPEM(t, "/C=CH/O=TestOrg/CN=Root")
		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isChain, err := genericx509utils.IsCertChain(ctx, []*x509.Certificate{rootCert})
		require.NoError(t, err)
		require.False(t, isChain)
	})

	t.Run("valid full chain ee-int-root", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isChain, err := genericx509utils.IsCertChain(ctx, []*x509.Certificate{eeCert, intCert, rootCert})
		require.NoError(t, err)
		require.True(t, isChain)
	})

	t.Run("valid partial chain int-root", func(t *testing.T) {
		rootPEM, intPEM, _ := generateChainCerts(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)

		isChain, err := genericx509utils.IsCertChain(ctx, []*x509.Certificate{intCert, rootCert})
		require.NoError(t, err)
		require.True(t, isChain)
	})

	t.Run("wrong order root-int-ee is not valid chain", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isChain, err := genericx509utils.IsCertChain(ctx, []*x509.Certificate{rootCert, intCert, eeCert})
		require.NoError(t, err)
		require.False(t, isChain)
	})

	t.Run("unrelated certs are not a chain", func(t *testing.T) {
		cert1PEM := generateUnrelatedSelfSignedCertPEM(t, "/C=CH/O=Org1/CN=Cert1")
		cert2PEM := generateUnrelatedSelfSignedCertPEM(t, "/C=DE/O=Org2/CN=Cert2")

		cert1, err := genericx509utils.ReadCertFromString(cert1PEM)
		require.NoError(t, err)
		cert2, err := genericx509utils.ReadCertFromString(cert2PEM)
		require.NoError(t, err)

		isChain, err := genericx509utils.IsCertChain(ctx, []*x509.Certificate{cert1, cert2})
		require.NoError(t, err)
		require.False(t, isChain)
	})

	t.Run("nil cert in chain returns error", func(t *testing.T) {
		rootPEM := generateUnrelatedSelfSignedCertPEM(t, "/C=CH/O=TestOrg/CN=Root")
		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isChain, err := genericx509utils.IsCertChain(ctx, []*x509.Certificate{nil, rootCert})
		require.Error(t, err)
		require.False(t, isChain)
	})
}

func Test_IsCertChainFromString(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("empty string", func(t *testing.T) {
		isChain, err := genericx509utils.IsCertChainFromString(ctx, "")
		require.Error(t, err)
		require.False(t, isChain)
	})

	t.Run("invalid PEM", func(t *testing.T) {
		isChain, err := genericx509utils.IsCertChainFromString(ctx, "not a PEM")
		require.Error(t, err)
		require.False(t, isChain)
	})

	t.Run("valid chain from string", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)
		chainStr := eePEM + intPEM + rootPEM

		isChain, err := genericx509utils.IsCertChainFromString(ctx, chainStr)
		require.NoError(t, err)
		require.True(t, isChain)
	})

	t.Run("wrong order from string is not valid", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)
		chainStr := rootPEM + intPEM + eePEM

		isChain, err := genericx509utils.IsCertChainFromString(ctx, chainStr)
		require.NoError(t, err)
		require.False(t, isChain)
	})
}

func Test_IsCertChainFromBytes(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil input", func(t *testing.T) {
		isChain, err := genericx509utils.IsCertChainFromBytes(ctx, nil)
		require.Error(t, err)
		require.False(t, isChain)
	})

	t.Run("valid chain from bytes", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)
		chainBytes := []byte(eePEM + intPEM + rootPEM)

		isChain, err := genericx509utils.IsCertChainFromBytes(ctx, chainBytes)
		require.NoError(t, err)
		require.True(t, isChain)
	})

	t.Run("consistent with IsCertChainFromString", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)
		chainStr := eePEM + intPEM + rootPEM

		isChainFromString, err := genericx509utils.IsCertChainFromString(ctx, chainStr)
		require.NoError(t, err)

		isChainFromBytes, err := genericx509utils.IsCertChainFromBytes(ctx, []byte(chainStr))
		require.NoError(t, err)

		require.Equal(t, isChainFromString, isChainFromBytes)
	})
}

func Test_IsRootCaToEndEntityChain(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil chain", func(t *testing.T) {
		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, nil)
		require.Error(t, err)
		require.False(t, isValid)
	})

	t.Run("empty chain", func(t *testing.T) {
		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{})
		require.NoError(t, err)
		require.False(t, isValid)
	})

	t.Run("single cert is not valid", func(t *testing.T) {
		rootPEM := generateUnrelatedSelfSignedCertPEM(t, "/C=CH/O=TestOrg/CN=Root")
		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{rootCert})
		require.NoError(t, err)
		require.False(t, isValid)
	})

	t.Run("valid full chain ee-int-root", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{eeCert, intCert, rootCert})
		require.NoError(t, err)
		require.True(t, isValid)
	})

	t.Run("first cert is not end entity returns false", func(t *testing.T) {
		rootPEM, intPEM, _ := generateChainCerts(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)

		// intermediate as first cert should fail
		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{intCert, rootCert})
		require.NoError(t, err)
		require.False(t, isValid)
	})

	t.Run("last cert is not root CA returns false", func(t *testing.T) {
		_, intPEM, eePEM := generateChainCerts(t)

		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		// intermediate as last cert should fail
		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{eeCert, intCert})
		require.NoError(t, err)
		require.False(t, isValid)
	})

	t.Run("wrong order root-int-ee returns false", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		intCert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{rootCert, intCert, eeCert})
		require.NoError(t, err)
		require.False(t, isValid)
	})

	t.Run("valid chain without intermediate ee-root", func(t *testing.T) {
		// Generate an end entity directly signed by root
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
			"-subj", "/C=CH/O=DirectRootOrg/CN=DirectRoot",
			"-addext", "basicConstraints=critical,CA:TRUE",
			"-addext", "keyUsage=critical,keyCertSign,cRLSign",
		)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "openssl root CA generation failed: %s", string(output))

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
			"-subj", "/C=CH/O=DirectEEOrg/CN=direct.example.com",
		)
		output, err = cmd.CombinedOutput()
		require.NoError(t, err, "openssl end entity CSR generation failed: %s", string(output))

		err = os.WriteFile(eeExtPath, []byte("basicConstraints=critical,CA:FALSE\nkeyUsage=critical,digitalSignature,keyEncipherment\n"), 0644)
		require.NoError(t, err)

		cmd = exec.Command(
			"openssl", "x509", "-req",
			"-in", eeCsrPath,
			"-CA", rootCertPath,
			"-CAkey", rootKeyPath,
			"-CAcreateserial",
			"-out", eeCertPath,
			"-days", "1",
			"-extfile", eeExtPath,
		)
		output, err = cmd.CombinedOutput()
		require.NoError(t, err, "openssl end entity cert signing failed: %s", string(output))

		rootPEMBytes, err := os.ReadFile(rootCertPath)
		require.NoError(t, err)
		eePEMBytes, err := os.ReadFile(eeCertPath)
		require.NoError(t, err)

		rootCert, err := genericx509utils.ReadCertFromString(string(rootPEMBytes))
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(string(eePEMBytes))
		require.NoError(t, err)

		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{eeCert, rootCert})
		require.NoError(t, err)
		require.True(t, isValid)
	})

	t.Run("middle cert is not intermediate returns false", func(t *testing.T) {
		rootPEM, _, eePEM := generateChainCerts(t)
		unrelatedPEM := generateUnrelatedSelfSignedCertPEM(t, "/C=FR/O=FakeOrg/CN=FakeRoot")

		rootCert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)
		eeCert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)
		unrelatedCert, err := genericx509utils.ReadCertFromString(unrelatedPEM)
		require.NoError(t, err)

		// unrelatedCert is a self-signed root CA, not an intermediate
		isValid, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{eeCert, unrelatedCert, rootCert})
		require.NoError(t, err)
		require.False(t, isValid)
	})
}

func Test_IsRootCaToEndEntityChainFromString(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("empty string", func(t *testing.T) {
		isValid, err := genericx509utils.IsRootCaToEndEntityChainFromString(ctx, "")
		require.Error(t, err)
		require.False(t, isValid)
	})

	t.Run("invalid PEM", func(t *testing.T) {
		isValid, err := genericx509utils.IsRootCaToEndEntityChainFromString(ctx, "not a PEM")
		require.Error(t, err)
		require.False(t, isValid)
	})

	t.Run("valid chain from string", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)
		chainStr := eePEM + intPEM + rootPEM

		isValid, err := genericx509utils.IsRootCaToEndEntityChainFromString(ctx, chainStr)
		require.NoError(t, err)
		require.True(t, isValid)
	})

	t.Run("wrong order from string returns false", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)
		chainStr := rootPEM + intPEM + eePEM

		isValid, err := genericx509utils.IsRootCaToEndEntityChainFromString(ctx, chainStr)
		require.NoError(t, err)
		require.False(t, isValid)
	})
}

func Test_IsRootCaToEndEntityChainFromBytes(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil input", func(t *testing.T) {
		isValid, err := genericx509utils.IsRootCaToEndEntityChainFromBytes(ctx, nil)
		require.Error(t, err)
		require.False(t, isValid)
	})

	t.Run("valid chain from bytes", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)
		chainBytes := []byte(eePEM + intPEM + rootPEM)

		isValid, err := genericx509utils.IsRootCaToEndEntityChainFromBytes(ctx, chainBytes)
		require.NoError(t, err)
		require.True(t, isValid)
	})

	t.Run("consistent with FromString", func(t *testing.T) {
		rootPEM, intPEM, eePEM := generateChainCerts(t)
		chainStr := eePEM + intPEM + rootPEM

		isValidFromString, err := genericx509utils.IsRootCaToEndEntityChainFromString(ctx, chainStr)
		require.NoError(t, err)

		isValidFromBytes, err := genericx509utils.IsRootCaToEndEntityChainFromBytes(ctx, []byte(chainStr))
		require.NoError(t, err)

		require.Equal(t, isValidFromString, isValidFromBytes)
	})
}

func Test_CheckCertificateChainString(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("empty string returns error", func(t *testing.T) {
		err := genericx509utils.CheckCertificateChainString(ctx, "")
		require.Error(t, err)
	})

	t.Run("invalid PEM returns error", func(t *testing.T) {
		err := genericx509utils.CheckCertificateChainString(ctx, "not a PEM")
		require.Error(t, err)
	})

	t.Run("single cert is not a valid chain", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		rootPEM, err := genericx509utils.WriteCertAsString(rootPair.Cert)
		require.NoError(t, err)

		err = genericx509utils.CheckCertificateChainString(ctx, rootPEM)
		require.Error(t, err)
	})

	t.Run("valid full chain ee-int-root", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), intPair)
		require.NoError(t, err)

		eePEM, err := genericx509utils.WriteCertAsString(eePair.Cert)
		require.NoError(t, err)
		intPEM, err := genericx509utils.WriteCertAsString(intPair.Cert)
		require.NoError(t, err)
		rootPEM, err := genericx509utils.WriteCertAsString(rootPair.Cert)
		require.NoError(t, err)

		chainStr := eePEM + intPEM + rootPEM

		err = genericx509utils.CheckCertificateChainString(ctx, chainStr)
		require.NoError(t, err)
	})

	t.Run("wrong order root-int-ee returns error", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		intPair, err := genericx509utils.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), intPair)
		require.NoError(t, err)

		eePEM, err := genericx509utils.WriteCertAsString(eePair.Cert)
		require.NoError(t, err)
		intPEM, err := genericx509utils.WriteCertAsString(intPair.Cert)
		require.NoError(t, err)
		rootPEM, err := genericx509utils.WriteCertAsString(rootPair.Cert)
		require.NoError(t, err)

		chainStr := rootPEM + intPEM + eePEM

		err = genericx509utils.CheckCertificateChainString(ctx, chainStr)
		require.Error(t, err)
	})

	t.Run("unrelated certs are not a valid chain", func(t *testing.T) {
		rootPair1, err := genericx509utils.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		rootPair2, err := genericx509utils.CreateRootCaCertificate(ctx, &x509options.X509CreateCertificateOptions{
			CountryName:    "FR",
			Organization:   "OtherOrg",
			CommonName:     "OtherRoot",
			PrivateKeySize: 2048,
		})
		require.NoError(t, err)

		pem1, err := genericx509utils.WriteCertAsString(rootPair1.Cert)
		require.NoError(t, err)
		pem2, err := genericx509utils.WriteCertAsString(rootPair2.Cert)
		require.NoError(t, err)

		chainStr := pem1 + pem2

		err = genericx509utils.CheckCertificateChainString(ctx, chainStr)
		require.Error(t, err)
	})

	t.Run("only end entity cert without CA is not valid", func(t *testing.T) {
		rootPair, err := genericx509utils.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
		require.NoError(t, err)

		eePair, err := genericx509utils.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
		require.NoError(t, err)

		eePEM, err := genericx509utils.WriteCertAsString(eePair.Cert)
		require.NoError(t, err)

		err = genericx509utils.CheckCertificateChainString(ctx, eePEM)
		require.Error(t, err)
	})
}
