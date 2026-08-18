package genericx509utils_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func generateRootCaCertPEM(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "root_key.pem")
	certPath := filepath.Join(tmpDir, "root_cert.pem")

	cmd := exec.Command(
		"openssl", "req", "-x509",
		"-newkey", "rsa:2048",
		"-keyout", keyPath,
		"-out", certPath,
		"-days", "1",
		"-nodes",
		"-subj", "/C=CH/O=RootOrg/CN=RootCA",
		"-addext", "basicConstraints=critical,CA:TRUE",
		"-addext", "keyUsage=critical,keyCertSign,cRLSign",
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl root CA generation failed: %s", string(output))

	certPEM, err := os.ReadFile(certPath)
	require.NoError(t, err)

	return string(certPEM)
}

func generateIntermediateCertPEM(t *testing.T) string {
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

	// Intermediate CSR
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

	// Extension file for intermediate CA
	extContent := "basicConstraints=critical,CA:TRUE,pathlen:0\nkeyUsage=critical,keyCertSign,cRLSign\n"
	err = os.WriteFile(intExtPath, []byte(extContent), 0644)
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

	certPEM, err := os.ReadFile(intCertPath)
	require.NoError(t, err)

	return string(certPEM)
}

func generateEndEntityCertPEM(t *testing.T) string {
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
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl root CA generation failed: %s", string(output))

	// End Entity CSR
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

	// Extension file for end entity (no CA)
	extContent := "basicConstraints=critical,CA:FALSE\nkeyUsage=critical,digitalSignature,keyEncipherment\n"
	err = os.WriteFile(eeExtPath, []byte(extContent), 0644)
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

	certPEM, err := os.ReadFile(eeCertPath)
	require.NoError(t, err)

	return string(certPEM)
}

func Test_IsRootCaCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil cert", func(t *testing.T) {
		isRootCa, err := genericx509utils.IsRootCaCert(ctx, nil)
		require.Error(t, err)
		require.False(t, isRootCa)
	})

	t.Run("root CA cert returns true", func(t *testing.T) {
		rootPEM := generateRootCaCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isRootCa)
	})

	t.Run("intermediate cert returns false", func(t *testing.T) {
		intPEM := generateIntermediateCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)

		isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isRootCa)
	})

	t.Run("end entity cert returns false", func(t *testing.T) {
		eePEM := generateEndEntityCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isRootCa)
	})
}

func Test_IsStringRootCaCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("empty string", func(t *testing.T) {
		isRootCa, err := genericx509utils.IsStringRootCaCert(ctx, "")
		require.Error(t, err)
		require.False(t, isRootCa)
	})

	t.Run("invalid PEM", func(t *testing.T) {
		isRootCa, err := genericx509utils.IsStringRootCaCert(ctx, "not a PEM")
		require.Error(t, err)
		require.False(t, isRootCa)
	})

	t.Run("root CA returns true", func(t *testing.T) {
		rootPEM := generateRootCaCertPEM(t)

		isRootCa, err := genericx509utils.IsStringRootCaCert(ctx, rootPEM)
		require.NoError(t, err)
		require.True(t, isRootCa)
	})

	t.Run("end entity returns false", func(t *testing.T) {
		eePEM := generateEndEntityCertPEM(t)

		isRootCa, err := genericx509utils.IsStringRootCaCert(ctx, eePEM)
		require.NoError(t, err)
		require.False(t, isRootCa)
	})
}

func Test_IsBytesRootCaCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil input", func(t *testing.T) {
		isRootCa, err := genericx509utils.IsBytesRootCaCert(ctx, nil)
		require.Error(t, err)
		require.False(t, isRootCa)
	})

	t.Run("root CA returns true", func(t *testing.T) {
		rootPEM := generateRootCaCertPEM(t)

		isRootCa, err := genericx509utils.IsBytesRootCaCert(ctx, []byte(rootPEM))
		require.NoError(t, err)
		require.True(t, isRootCa)
	})
}

func Test_IsIntermediateCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil cert", func(t *testing.T) {
		isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, nil)
		require.Error(t, err)
		require.False(t, isIntermediate)
	})

	t.Run("intermediate cert returns true", func(t *testing.T) {
		intPEM := generateIntermediateCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)

		isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isIntermediate)
	})

	t.Run("root CA cert returns false", func(t *testing.T) {
		rootPEM := generateRootCaCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isIntermediate)
	})

	t.Run("end entity cert returns false", func(t *testing.T) {
		eePEM := generateEndEntityCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isIntermediate)
	})
}

func Test_IsStringIntermediateCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("empty string", func(t *testing.T) {
		isIntermediate, err := genericx509utils.IsStringIntermediateCert(ctx, "")
		require.Error(t, err)
		require.False(t, isIntermediate)
	})

	t.Run("invalid PEM", func(t *testing.T) {
		isIntermediate, err := genericx509utils.IsStringIntermediateCert(ctx, "garbage")
		require.Error(t, err)
		require.False(t, isIntermediate)
	})

	t.Run("intermediate returns true", func(t *testing.T) {
		intPEM := generateIntermediateCertPEM(t)

		isIntermediate, err := genericx509utils.IsStringIntermediateCert(ctx, intPEM)
		require.NoError(t, err)
		require.True(t, isIntermediate)
	})

	t.Run("root CA returns false", func(t *testing.T) {
		rootPEM := generateRootCaCertPEM(t)

		isIntermediate, err := genericx509utils.IsStringIntermediateCert(ctx, rootPEM)
		require.NoError(t, err)
		require.False(t, isIntermediate)
	})
}

func Test_IsBytesIntermediateCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil input", func(t *testing.T) {
		isIntermediate, err := genericx509utils.IsBytesIntermediateCert(ctx, nil)
		require.Error(t, err)
		require.False(t, isIntermediate)
	})

	t.Run("intermediate returns true", func(t *testing.T) {
		intPEM := generateIntermediateCertPEM(t)

		isIntermediate, err := genericx509utils.IsBytesIntermediateCert(ctx, []byte(intPEM))
		require.NoError(t, err)
		require.True(t, isIntermediate)
	})
}

func Test_IsEndEntityCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil cert", func(t *testing.T) {
		isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, nil)
		require.Error(t, err)
		require.False(t, isEndEntity)
	})

	t.Run("end entity cert returns true", func(t *testing.T) {
		eePEM := generateEndEntityCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isEndEntity)
	})

	t.Run("root CA cert returns false", func(t *testing.T) {
		rootPEM := generateRootCaCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isEndEntity)
	})

	t.Run("intermediate cert returns false", func(t *testing.T) {
		intPEM := generateIntermediateCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)

		isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isEndEntity)
	})
}

func Test_IsStringEndEntityCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("empty string", func(t *testing.T) {
		isEndEntity, err := genericx509utils.IsStringEndEntityCert(ctx, "")
		require.Error(t, err)
		require.False(t, isEndEntity)
	})

	t.Run("invalid PEM", func(t *testing.T) {
		isEndEntity, err := genericx509utils.IsStringEndEntityCert(ctx, "invalid")
		require.Error(t, err)
		require.False(t, isEndEntity)
	})

	t.Run("end entity returns true", func(t *testing.T) {
		eePEM := generateEndEntityCertPEM(t)

		isEndEntity, err := genericx509utils.IsStringEndEntityCert(ctx, eePEM)
		require.NoError(t, err)
		require.True(t, isEndEntity)
	})

	t.Run("root CA returns false", func(t *testing.T) {
		rootPEM := generateRootCaCertPEM(t)

		isEndEntity, err := genericx509utils.IsStringEndEntityCert(ctx, rootPEM)
		require.NoError(t, err)
		require.False(t, isEndEntity)
	})
}

func Test_IsBytesEndEntityCert(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("nil input", func(t *testing.T) {
		isEndEntity, err := genericx509utils.IsBytesEndEntityCert(ctx, nil)
		require.Error(t, err)
		require.False(t, isEndEntity)
	})

	t.Run("end entity returns true", func(t *testing.T) {
		eePEM := generateEndEntityCertPEM(t)

		isEndEntity, err := genericx509utils.IsBytesEndEntityCert(ctx, []byte(eePEM))
		require.NoError(t, err)
		require.True(t, isEndEntity)
	})
}

func Test_MutualExclusivity(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	t.Run("root CA is exclusively root CA", func(t *testing.T) {
		rootPEM := generateRootCaCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(rootPEM)
		require.NoError(t, err)

		isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isRootCa)

		isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isIntermediate)

		isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isEndEntity)
	})

	t.Run("intermediate is exclusively intermediate", func(t *testing.T) {
		intPEM := generateIntermediateCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(intPEM)
		require.NoError(t, err)

		isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isRootCa)

		isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isIntermediate)

		isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isEndEntity)
	})

	t.Run("end entity is exclusively end entity", func(t *testing.T) {
		eePEM := generateEndEntityCertPEM(t)
		cert, err := genericx509utils.ReadCertFromString(eePEM)
		require.NoError(t, err)

		isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isRootCa)

		isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, cert)
		require.NoError(t, err)
		require.False(t, isIntermediate)

		isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, cert)
		require.NoError(t, err)
		require.True(t, isEndEntity)
	})
}

// Ensure unused import is satisfied
var _ context.Context
