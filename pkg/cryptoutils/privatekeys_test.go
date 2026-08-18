package cryptoutils_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/cryptoutils"
)

func generateRSAPrivateKeyPEM(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "key.pem")

	// Generate RSA key in PKCS8 format (which uses "PRIVATE KEY" header)
	cmd := exec.Command(
		"openssl", "genpkey",
		"-algorithm", "RSA",
		"-pkeyopt", "rsa_keygen_bits:2048",
		"-out", keyPath,
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl genpkey failed: %s", string(output))

	keyPEM, err := os.ReadFile(keyPath)
	require.NoError(t, err)

	return string(keyPEM)
}

func generateECPrivateKeyPEM(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "key.pem")

	cmd := exec.Command(
		"openssl", "genpkey",
		"-algorithm", "EC",
		"-pkeyopt", "ec_paramgen_curve:P-256",
		"-out", keyPath,
	)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "openssl genpkey failed: %s", string(output))

	keyPEM, err := os.ReadFile(keyPath)
	require.NoError(t, err)

	return string(keyPEM)
}

func Test_LoadPrivateKeyFromPEMString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		key, err := cryptoutils.LoadPrivateKeyFromPEMString("")
		require.Error(t, err)
		require.Nil(t, key)
	})

	t.Run("invalid PEM", func(t *testing.T) {
		key, err := cryptoutils.LoadPrivateKeyFromPEMString("not a valid PEM")
		require.Error(t, err)
		require.Nil(t, key)
	})

	t.Run("valid RSA key", func(t *testing.T) {
		keyPEM := generateRSAPrivateKeyPEM(t)

		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)
		require.NotNil(t, key)
	})

	t.Run("valid EC key", func(t *testing.T) {
		keyPEM := generateECPrivateKeyPEM(t)

		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)
		require.NotNil(t, key)
	})

	t.Run("wrong PEM type rejected", func(t *testing.T) {
		// Generate a cert PEM (not a key)
		tmpDir := t.TempDir()
		certPath := filepath.Join(tmpDir, "cert.pem")
		keyPath := filepath.Join(tmpDir, "key.pem")

		cmd := exec.Command(
			"openssl", "req", "-x509",
			"-newkey", "rsa:2048",
			"-keyout", keyPath,
			"-out", certPath,
			"-days", "1",
			"-nodes",
			"-subj", "/CN=Test",
		)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "openssl command failed: %s", string(output))

		certPEM, err := os.ReadFile(certPath)
		require.NoError(t, err)

		key, err := cryptoutils.LoadPrivateKeyFromPEMString(string(certPEM))
		require.Error(t, err)
		require.Nil(t, key)
	})
}

func Test_EncodePrivateKeyAsPEMString(t *testing.T) {
	t.Run("nil key", func(t *testing.T) {
		pemStr, err := cryptoutils.EncodePrivateKeyAsPEMString(nil)
		require.Error(t, err)
		require.Empty(t, pemStr)
	})

	t.Run("RSA key roundtrip", func(t *testing.T) {
		keyPEM := generateRSAPrivateKeyPEM(t)

		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		encoded, err := cryptoutils.EncodePrivateKeyAsPEMString(key)
		require.NoError(t, err)
		require.NotEmpty(t, encoded)
		require.Contains(t, encoded, "BEGIN PRIVATE KEY")
		require.Contains(t, encoded, "END PRIVATE KEY")

		// Roundtrip: parse back
		keyParsedBack, err := cryptoutils.LoadPrivateKeyFromPEMString(encoded)
		require.NoError(t, err)
		require.NotNil(t, keyParsedBack)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(key, keyParsedBack)
		require.NoError(t, err)
		require.True(t, isEqual)
	})

	t.Run("EC key roundtrip", func(t *testing.T) {
		keyPEM := generateECPrivateKeyPEM(t)

		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		encoded, err := cryptoutils.EncodePrivateKeyAsPEMString(key)
		require.NoError(t, err)
		require.NotEmpty(t, encoded)
		require.Contains(t, encoded, "PRIVATE KEY")

		keyParsedBack, err := cryptoutils.LoadPrivateKeyFromPEMString(encoded)
		require.NoError(t, err)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(key, keyParsedBack)
		require.NoError(t, err)
		require.True(t, isEqual)
	})

	t.Run("encoded PEM ends with single newline", func(t *testing.T) {
		keyPEM := generateRSAPrivateKeyPEM(t)

		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		encoded, err := cryptoutils.EncodePrivateKeyAsPEMString(key)
		require.NoError(t, err)

		require.True(t, encoded[len(encoded)-1] == '\n')
		require.False(t, encoded[len(encoded)-2] == '\n')
	})
}

func Test_GetPublicKeyFromPrivateKey(t *testing.T) {
	t.Run("nil key", func(t *testing.T) {
		pubKey, err := cryptoutils.GetPublicKeyFromPrivateKey(nil)
		require.Error(t, err)
		require.Nil(t, pubKey)
	})

	t.Run("RSA key returns public key", func(t *testing.T) {
		keyPEM := generateRSAPrivateKeyPEM(t)

		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		pubKey, err := cryptoutils.GetPublicKeyFromPrivateKey(key)
		require.NoError(t, err)
		require.NotNil(t, pubKey)
	})

	t.Run("EC key returns public key", func(t *testing.T) {
		keyPEM := generateECPrivateKeyPEM(t)

		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		pubKey, err := cryptoutils.GetPublicKeyFromPrivateKey(key)
		require.NoError(t, err)
		require.NotNil(t, pubKey)
	})

	t.Run("two different keys have different public keys", func(t *testing.T) {
		key1PEM := generateRSAPrivateKeyPEM(t)
		key2PEM := generateRSAPrivateKeyPEM(t)

		key1, err := cryptoutils.LoadPrivateKeyFromPEMString(key1PEM)
		require.NoError(t, err)
		key2, err := cryptoutils.LoadPrivateKeyFromPEMString(key2PEM)
		require.NoError(t, err)

		pubKey1, err := cryptoutils.GetPublicKeyFromPrivateKey(key1)
		require.NoError(t, err)
		pubKey2, err := cryptoutils.GetPublicKeyFromPrivateKey(key2)
		require.NoError(t, err)

		// Use the Equal interface to compare
		pubKey1WithEqual, ok := pubKey1.(interface{ Equal(x interface{}) bool })
		if ok {
			require.False(t, pubKey1WithEqual.Equal(pubKey2))
		}
	})
}

func Test_IsPrivateKeyEqual(t *testing.T) {
	t.Run("nil key1", func(t *testing.T) {
		keyPEM := generateRSAPrivateKeyPEM(t)
		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(nil, key)
		require.Error(t, err)
		require.False(t, isEqual)
	})

	t.Run("nil key2", func(t *testing.T) {
		keyPEM := generateRSAPrivateKeyPEM(t)
		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(key, nil)
		require.Error(t, err)
		require.False(t, isEqual)
	})

	t.Run("same key is equal to itself", func(t *testing.T) {
		keyPEM := generateRSAPrivateKeyPEM(t)
		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(key, key)
		require.NoError(t, err)
		require.True(t, isEqual)
	})

	t.Run("key parsed twice from same PEM is equal", func(t *testing.T) {
		keyPEM := generateRSAPrivateKeyPEM(t)

		key1, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		key2, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(key1, key2)
		require.NoError(t, err)
		require.True(t, isEqual)
	})

	t.Run("two different RSA keys are not equal", func(t *testing.T) {
		key1PEM := generateRSAPrivateKeyPEM(t)
		key2PEM := generateRSAPrivateKeyPEM(t)

		key1, err := cryptoutils.LoadPrivateKeyFromPEMString(key1PEM)
		require.NoError(t, err)
		key2, err := cryptoutils.LoadPrivateKeyFromPEMString(key2PEM)
		require.NoError(t, err)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(key1, key2)
		require.NoError(t, err)
		require.False(t, isEqual)
	})

	t.Run("two different EC keys are not equal", func(t *testing.T) {
		key1PEM := generateECPrivateKeyPEM(t)
		key2PEM := generateECPrivateKeyPEM(t)

		key1, err := cryptoutils.LoadPrivateKeyFromPEMString(key1PEM)
		require.NoError(t, err)
		key2, err := cryptoutils.LoadPrivateKeyFromPEMString(key2PEM)
		require.NoError(t, err)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(key1, key2)
		require.NoError(t, err)
		require.False(t, isEqual)
	})

	t.Run("RSA and EC key are not equal", func(t *testing.T) {
		rsaKeyPEM := generateRSAPrivateKeyPEM(t)
		ecKeyPEM := generateECPrivateKeyPEM(t)

		rsaKey, err := cryptoutils.LoadPrivateKeyFromPEMString(rsaKeyPEM)
		require.NoError(t, err)
		ecKey, err := cryptoutils.LoadPrivateKeyFromPEMString(ecKeyPEM)
		require.NoError(t, err)

		isEqual, err := cryptoutils.IsPrivateKeyEqual(rsaKey, ecKey)
		require.NoError(t, err)
		require.False(t, isEqual)
	})
}
