package genericx509utils_test

import (
	"context"
	"crypto"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/cryptoutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

func generateCertAndKeyPEM(t *testing.T, subject string) (certPEM string, keyPEM string) {
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

	certPEMBytes, err := os.ReadFile(certPath)
	require.NoError(t, err)

	keyPEMBytes, err := os.ReadFile(keyPath)
	require.NoError(t, err)

	return string(certPEMBytes), string(keyPEMBytes)
}

func loadCertAndKeyPair(t *testing.T, subject string) *genericx509utils.X509CertKeyPair {
	t.Helper()

	certPEM, keyPEM := generateCertAndKeyPEM(t, subject)

	cert, err := genericx509utils.ReadCertFromString(certPEM)
	require.NoError(t, err)

	key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
	require.NoError(t, err)

	return &genericx509utils.X509CertKeyPair{
		Cert: cert,
		Key:  key,
	}
}

func Test_X509CertKeyPair_GetX509Certificate(t *testing.T) {
	t.Run("cert not set", func(t *testing.T) {
		pair := &genericx509utils.X509CertKeyPair{}

		cert, err := pair.GetX509Certificate()
		require.Error(t, err)
		require.Nil(t, cert)
	})

	t.Run("returns deep copy of cert", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=TestOrg/CN=TestCert")

		cert1, err := pair.GetX509Certificate()
		require.NoError(t, err)
		require.NotNil(t, cert1)

		cert2, err := pair.GetX509Certificate()
		require.NoError(t, err)
		require.NotNil(t, cert2)

		require.True(t, cert1.Equal(cert2))
		require.True(t, cert1.Equal(pair.Cert))
	})

	t.Run("returned cert has correct subject", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/L=Zurich/O=DeepCopyOrg/CN=DeepCopyCert")

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)
		require.Equal(t, "DeepCopyCert", cert.Subject.CommonName)
		require.Equal(t, "CH", cert.Subject.Country[0])
		require.Equal(t, "Zurich", cert.Subject.Locality[0])
		require.Equal(t, "DeepCopyOrg", cert.Subject.Organization[0])
	})
}

func Test_X509CertKeyPair_GetPrivateKey(t *testing.T) {
	t.Run("key not set", func(t *testing.T) {
		pair := &genericx509utils.X509CertKeyPair{}

		key, err := pair.GetPrivateKey()
		require.Error(t, err)
		require.Nil(t, key)
	})

	t.Run("returns private key", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=TestOrg/CN=TestCert")

		key, err := pair.GetPrivateKey()
		require.NoError(t, err)
		require.NotNil(t, key)
	})
}

func Test_X509CertKeyPair_GetPrivateKeyAsPEMString(t *testing.T) {
	t.Run("key not set", func(t *testing.T) {
		pair := &genericx509utils.X509CertKeyPair{}

		pemStr, err := pair.GetPrivateKeyAsPEMString()
		require.Error(t, err)
		require.Empty(t, pemStr)
	})

	t.Run("returns non empty PEM string", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=TestOrg/CN=TestCert")

		pemStr, err := pair.GetPrivateKeyAsPEMString()
		require.NoError(t, err)
		require.NotEmpty(t, pemStr)
		require.Contains(t, pemStr, "PRIVATE KEY")
	})

	t.Run("PEM string can be parsed back", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=TestOrg/CN=TestCert")

		pemStr, err := pair.GetPrivateKeyAsPEMString()
		require.NoError(t, err)

		parsedKey, err := cryptoutils.LoadPrivateKeyFromPEMString(pemStr)
		require.NoError(t, err)
		require.NotNil(t, parsedKey)
	})
}

func Test_X509CertKeyPair_IsKeyMatchingCert(t *testing.T) {
	t.Run("cert not set", func(t *testing.T) {
		_, keyPEM := generateCertAndKeyPEM(t, "/C=CH/O=TestOrg/CN=TestCert")
		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		pair := &genericx509utils.X509CertKeyPair{Key: key}

		isMatching, err := pair.IsKeyMatchingCert()
		require.Error(t, err)
		require.False(t, isMatching)
	})

	t.Run("key not set", func(t *testing.T) {
		certPEM, _ := generateCertAndKeyPEM(t, "/C=CH/O=TestOrg/CN=TestCert")
		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)

		pair := &genericx509utils.X509CertKeyPair{Cert: cert}

		isMatching, err := pair.IsKeyMatchingCert()
		require.Error(t, err)
		require.False(t, isMatching)
	})

	t.Run("matching cert and key returns true", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=MatchOrg/CN=MatchCert")

		isMatching, err := pair.IsKeyMatchingCert()
		require.NoError(t, err)
		require.True(t, isMatching)
	})

	t.Run("mismatched cert and key returns false", func(t *testing.T) {
		certPEM, _ := generateCertAndKeyPEM(t, "/C=CH/O=Cert1Org/CN=Cert1")
		_, keyPEM := generateCertAndKeyPEM(t, "/C=DE/O=Cert2Org/CN=Cert2")

		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)
		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		pair := &genericx509utils.X509CertKeyPair{Cert: cert, Key: key}

		isMatching, err := pair.IsKeyMatchingCert()
		require.NoError(t, err)
		require.False(t, isMatching)
	})
}

func Test_X509CertKeyPair_CheckKeyMatchingCert(t *testing.T) {
	t.Run("matching cert and key returns no error", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=CheckOrg/CN=CheckCert")

		err := pair.CheckKeyMatchingCertificate()
		require.NoError(t, err)
	})

	t.Run("mismatched cert and key returns error", func(t *testing.T) {
		certPEM, _ := generateCertAndKeyPEM(t, "/C=CH/O=Cert1Org/CN=Cert1")
		_, keyPEM := generateCertAndKeyPEM(t, "/C=DE/O=Cert2Org/CN=Cert2")

		cert, err := genericx509utils.ReadCertFromString(certPEM)
		require.NoError(t, err)
		key, err := cryptoutils.LoadPrivateKeyFromPEMString(keyPEM)
		require.NoError(t, err)

		pair := &genericx509utils.X509CertKeyPair{Cert: cert, Key: key}

		err = pair.CheckKeyMatchingCertificate()
		require.Error(t, err)
	})

	t.Run("cert not set returns error", func(t *testing.T) {
		pair := &genericx509utils.X509CertKeyPair{}

		err := pair.CheckKeyMatchingCertificate()
		require.Error(t, err)
	})
}

func Test_X509CertKeyPair_GetPublicKey(t *testing.T) {
	t.Run("key not set", func(t *testing.T) {
		pair := &genericx509utils.X509CertKeyPair{}

		pubKey, err := pair.GetPublicKey()
		require.Error(t, err)
		require.Nil(t, pubKey)
	})

	t.Run("returns non nil public key", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=PubKeyOrg/CN=PubKeyCert")

		pubKey, err := pair.GetPublicKey()
		require.NoError(t, err)
		require.NotNil(t, pubKey)
	})

	t.Run("public key matches certificate public key", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=MatchPubOrg/CN=MatchPubCert")

		pubKey, err := pair.GetPublicKey()
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		certPubKey := cert.PublicKey

		certPubKeyWithEqual, ok := certPubKey.(interface{ Equal(x crypto.PublicKey) bool })
		require.True(t, ok)
		require.True(t, certPubKeyWithEqual.Equal(pubKey))
	})
}

func Test_X509CertKeyPair_WriteCertificatePemToFile(t *testing.T) {
	t.Run("nil toWrite returns error", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=WriteCertOrg/CN=WriteCert")

		err := pair.WriteCertificatePemToFile(context.Background(), nil)
		require.Error(t, err)
	})

	t.Run("writes certificate PEM to file", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=WriteCertFileOrg/CN=WriteCertFile")

		tmpDir := t.TempDir()
		certPath := filepath.Join(tmpDir, "cert.pem")
		certFile := files.MustNewLocalFileByPath(certPath)

		err := pair.WriteCertificatePemToFile(context.Background(), certFile)
		require.NoError(t, err)

		content, err := os.ReadFile(certPath)
		require.NoError(t, err)
		require.NotEmpty(t, content)
		require.Contains(t, string(content), "BEGIN CERTIFICATE")
		require.Contains(t, string(content), "END CERTIFICATE")

		certReadBack, err := genericx509utils.ReadCertFromString(string(content))
		require.NoError(t, err)

		originalCert, err := pair.GetX509Certificate()
		require.NoError(t, err)
		require.True(t, originalCert.Equal(certReadBack))
	})

	t.Run("written certificate has correct subject", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=DE/L=Munich/O=SubjectWriteOrg/CN=SubjectWriteCert")

		tmpDir := t.TempDir()
		certPath := filepath.Join(tmpDir, "cert.pem")
		certFile := files.MustNewLocalFileByPath(certPath)

		err := pair.WriteCertificatePemToFile(context.Background(), certFile)
		require.NoError(t, err)

		content, err := os.ReadFile(certPath)
		require.NoError(t, err)

		certReadBack, err := genericx509utils.ReadCertFromString(string(content))
		require.NoError(t, err)
		require.Equal(t, "DE", certReadBack.Subject.Country[0])
		require.Equal(t, "Munich", certReadBack.Subject.Locality[0])
		require.Equal(t, "SubjectWriteOrg", certReadBack.Subject.Organization[0])
		require.Equal(t, "SubjectWriteCert", certReadBack.Subject.CommonName)
	})
}

func Test_X509CertKeyPair_WritePrivateKeyToFile(t *testing.T) {
	t.Run("nil toWrite returns error", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=WriteKeyOrg/CN=WriteKey")

		err := pair.WritePrivateKeyToFile(context.Background(), nil)
		require.Error(t, err)
	})

	t.Run("writes private key PEM to file", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=WriteKeyFileOrg/CN=WriteKeyFile")

		tmpDir := t.TempDir()
		keyPath := filepath.Join(tmpDir, "key.pem")
		keyFile := files.MustNewLocalFileByPath(keyPath)

		err := pair.WritePrivateKeyToFile(context.Background(), keyFile)
		require.NoError(t, err)

		content, err := os.ReadFile(keyPath)
		require.NoError(t, err)
		require.NotEmpty(t, content)
		require.Contains(t, string(content), "BEGIN PRIVATE KEY")
		require.Contains(t, string(content), "END PRIVATE KEY")

		keyReadBack, err := cryptoutils.LoadPrivateKeyFromPEMString(string(content))
		require.NoError(t, err)
		require.NotNil(t, keyReadBack)

		originalKey, err := pair.GetPrivateKey()
		require.NoError(t, err)

		keyEqual, ok := originalKey.(interface {
			Equal(x crypto.PrivateKey) bool
		})
		require.True(t, ok)
		require.True(t, keyEqual.Equal(keyReadBack))
	})

	t.Run("written key can decrypt data encrypted by certificate public key", func(t *testing.T) {
		pair := loadCertAndKeyPair(t, "/C=CH/O=DecryptTestOrg/CN=DecryptTest")

		tmpDir := t.TempDir()
		keyPath := filepath.Join(tmpDir, "key.pem")
		keyFile := files.MustNewLocalFileByPath(keyPath)

		err := pair.WritePrivateKeyToFile(context.Background(), keyFile)
		require.NoError(t, err)

		content, err := os.ReadFile(keyPath)
		require.NoError(t, err)

		keyReadBack, err := cryptoutils.LoadPrivateKeyFromPEMString(string(content))
		require.NoError(t, err)

		cert, err := pair.GetX509Certificate()
		require.NoError(t, err)

		publicKey := cert.PublicKey
		require.NotNil(t, publicKey)
		require.Equal(t, keyReadBack.(interface{ Public() crypto.PublicKey }).Public(), publicKey)
	})
}
