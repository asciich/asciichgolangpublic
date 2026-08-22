package x509utils_test

import (
	"context"
	"crypto"
	"crypto/x509"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/cryptoutils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/commandexecutorx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/nativex509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
)

// --- Implementation struct and registry ---

// x509Implementation holds a named set of function implementations
// so all implementations can be tested uniformly in a loop.
type x509Implementation struct {
	Name                                string
	ReadCertificateFromFile             func(ctx context.Context, pathToRead string) (*x509.Certificate, error)
	CreateRootCaCertificate             func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error)
	CreateIntermediateCertificate       func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error)
	CreateSelfSignedCertificate         func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error)
	CreateSignedIntermediateCertificate func(ctx context.Context, options *x509options.X509CreateCertificateOptions, rootCaCertAndKey *genericx509utils.X509CertKeyPair) (*genericx509utils.X509CertKeyPair, error)
	CreateSignedEndEntityCertificate    func(ctx context.Context, options *x509options.X509CreateCertificateOptions, caCertAndKey *genericx509utils.X509CertKeyPair) (*genericx509utils.X509CertKeyPair, error)
	GeneratePrivateKey                  func(ctx context.Context) (crypto.PrivateKey, error)
}

// getX509Implementations returns all implementations to test.
// This includes the convenience functions from the x509utils package itself.
func getX509Implementations() []x509Implementation {
	return []x509Implementation{
		{
			Name:                                "x509utils",
			ReadCertificateFromFile:             x509utils.ReadCertificateFromFile,
			CreateRootCaCertificate:             x509utils.CreateRootCaCertificate,
			CreateIntermediateCertificate:       x509utils.CreateIntermediateCertificate,
			CreateSelfSignedCertificate:         x509utils.CreateSelfSignedCertificate,
			CreateSignedIntermediateCertificate: x509utils.CreateSignedIntermediateCertificate,
			CreateSignedEndEntityCertificate:    x509utils.CreateSignedEndEntityCertificate,
			GeneratePrivateKey:                  x509utils.GeneratePrivateKey,
		},
		{
			Name:                                "nativex509utils",
			ReadCertificateFromFile:             nativex509utils.ReadCertificateFromFile,
			CreateRootCaCertificate:             nativex509utils.CreateRootCaCertificate,
			CreateIntermediateCertificate:       nativex509utils.CreateIntermediateCertificate,
			CreateSelfSignedCertificate:         nativex509utils.CreateSelfSignedCertificate,
			CreateSignedIntermediateCertificate: nativex509utils.CreateSignedIntermediateCertificate,
			CreateSignedEndEntityCertificate:    nativex509utils.CreateSignedEndEntityCertificate,
			GeneratePrivateKey:                  nativex509utils.GeneratePrivateKey,
		},
		{
			Name: "commandexecutorx509utils_exec",
			ReadCertificateFromFile: func(ctx context.Context, pathToRead string) (*x509.Certificate, error) {
				return commandexecutorx509utils.ReadCertificateFromFile(ctx, commandexecutorexecoo.Exec(), pathToRead)
			},
			CreateRootCaCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateRootCaCertificate(ctx, commandexecutorexecoo.Exec(), options)
			},
			CreateIntermediateCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateIntermediateCertificate(ctx, commandexecutorexecoo.Exec(), options)
			},
			CreateSelfSignedCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateSelfSignedCertificate(ctx, commandexecutorexecoo.Exec(), options)
			},
			CreateSignedIntermediateCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions, rootCaCertAndKey *genericx509utils.X509CertKeyPair) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateSignedIntermediateCertificate(ctx, commandexecutorexecoo.Exec(), options, rootCaCertAndKey)
			},
			CreateSignedEndEntityCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions, caCertAndKey *genericx509utils.X509CertKeyPair) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateSignedEndEntityCertificate(ctx, commandexecutorexecoo.Exec(), options, caCertAndKey)
			},
			GeneratePrivateKey: func(ctx context.Context) (crypto.PrivateKey, error) {
				return commandexecutorx509utils.GeneratePrivateKey(ctx, commandexecutorexecoo.Exec())
			},
		},
		{
			Name: "commandexecutorx509utils_bash",
			ReadCertificateFromFile: func(ctx context.Context, pathToRead string) (*x509.Certificate, error) {
				return commandexecutorx509utils.ReadCertificateFromFile(ctx, commandexecutorbashoo.Bash(), pathToRead)
			},
			CreateRootCaCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateRootCaCertificate(ctx, commandexecutorbashoo.Bash(), options)
			},
			CreateIntermediateCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateIntermediateCertificate(ctx, commandexecutorbashoo.Bash(), options)
			},
			CreateSelfSignedCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateSelfSignedCertificate(ctx, commandexecutorbashoo.Bash(), options)
			},
			CreateSignedIntermediateCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions, rootCaCertAndKey *genericx509utils.X509CertKeyPair) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateSignedIntermediateCertificate(ctx, commandexecutorbashoo.Bash(), options, rootCaCertAndKey)
			},
			CreateSignedEndEntityCertificate: func(ctx context.Context, options *x509options.X509CreateCertificateOptions, caCertAndKey *genericx509utils.X509CertKeyPair) (*genericx509utils.X509CertKeyPair, error) {
				return commandexecutorx509utils.CreateSignedEndEntityCertificate(ctx, commandexecutorbashoo.Bash(), options, caCertAndKey)
			},
			GeneratePrivateKey: func(ctx context.Context) (crypto.PrivateKey, error) {
				return commandexecutorx509utils.GeneratePrivateKey(ctx, commandexecutorbashoo.Bash())
			},
		},
	}
}

// --- Test helpers ---

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

func getDefaultRootCaOptions() *x509options.X509CreateCertificateOptions {
	return &x509options.X509CreateCertificateOptions{
		CountryName:    "CH",
		Locality:       "Zurich",
		Organization:   "TestRootOrg",
		CommonName:     "TestRootCA",
		PrivateKeySize: 2048,
	}
}

func getDefaultIntermediateOptions() *x509options.X509CreateCertificateOptions {
	return &x509options.X509CreateCertificateOptions{
		CountryName:    "CH",
		Locality:       "Bern",
		Organization:   "TestIntOrg",
		CommonName:     "TestIntermediateCA",
		PrivateKeySize: 2048,
	}
}

func getDefaultEndEntityOptions() *x509options.X509CreateCertificateOptions {
	return &x509options.X509CreateCertificateOptions{
		CountryName:    "CH",
		Locality:       "Basel",
		Organization:   "TestEEOrg",
		CommonName:     "server.example.com",
		PrivateKeySize: 2048,
	}
}

func getDefaultSelfSignedOptions() *x509options.X509CreateCertificateOptions {
	return &x509options.X509CreateCertificateOptions{
		CountryName:    "DE",
		Locality:       "Berlin",
		Organization:   "SelfSignedOrg",
		CommonName:     "selfsigned.example.com",
		PrivateKeySize: 2048,
	}
}

// --- ReadCertificateFromFile ---

// Test_ReadCertificateFromFile validates that all implementations of
// ReadCertificateFromFile behave identically.
func Test_ReadCertificateFromFile(t *testing.T) {
	implementations := getX509Implementations()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_empty path returns error", func(t *testing.T) {
			cert, err := impl.ReadCertificateFromFile(contextutils.ContextVerbose(), "")
			require.Error(t, err)
			require.Nil(t, cert)
		})

		t.Run(impl.Name+"_non-existent file returns error", func(t *testing.T) {
			cert, err := impl.ReadCertificateFromFile(contextutils.ContextVerbose(), "/non/existent/path.pem")
			require.Error(t, err)
			require.Nil(t, cert)
		})

		t.Run(impl.Name+"_valid certificate file", func(t *testing.T) {
			tmpDir := t.TempDir()
			certPath := generateSelfSignedCert(t, tmpDir)

			cert, err := impl.ReadCertificateFromFile(contextutils.ContextVerbose(), certPath)
			require.NoError(t, err)
			require.NotNil(t, cert)
			require.Equal(t, "TestOrg", cert.Subject.Organization[0])
			require.Equal(t, "CH", cert.Subject.Country[0])
			require.Equal(t, "Zurich", cert.Subject.Locality[0])
			require.Equal(t, "TestCert", cert.Subject.CommonName)
		})

		t.Run(impl.Name+"_invalid certificate file returns error", func(t *testing.T) {
			tmpDir := t.TempDir()
			invalidPath := filepath.Join(tmpDir, "invalid.pem")
			err := os.WriteFile(invalidPath, []byte("not a valid certificate"), 0644)
			require.NoError(t, err)

			cert, err := impl.ReadCertificateFromFile(contextutils.ContextVerbose(), invalidPath)
			require.Error(t, err)
			require.Nil(t, cert)
		})

		t.Run(impl.Name+"_reading same file twice yields equal certs", func(t *testing.T) {
			tmpDir := t.TempDir()
			certPath := generateSelfSignedCert(t, tmpDir)

			ctx := contextutils.ContextVerbose()

			cert1, err := impl.ReadCertificateFromFile(ctx, certPath)
			require.NoError(t, err)

			cert2, err := impl.ReadCertificateFromFile(ctx, certPath)
			require.NoError(t, err)

			require.True(t, cert1.Equal(cert2), "Reading the same certificate file twice should yield equal certificates")
		})
	}
}

// Test_ReadCertificateFromFile_implementations_return_same_result validates that
// all implementations return the exact same certificate for the same input file.
func Test_ReadCertificateFromFile_implementations_return_same_result(t *testing.T) {
	tmpDir := t.TempDir()
	certPath := generateSelfSignedCert(t, tmpDir)

	ctx := contextutils.ContextVerbose()
	implementations := getX509Implementations()

	var referenceCert *x509.Certificate
	var referenceName string

	for _, impl := range implementations {
		cert, err := impl.ReadCertificateFromFile(ctx, certPath)
		require.NoError(t, err)
		require.NotNil(t, cert)

		if referenceCert == nil {
			referenceCert = cert
			referenceName = impl.Name
		} else {
			require.True(
				t,
				referenceCert.Equal(cert),
				"Implementation '%s' returned a different certificate than '%s'",
				impl.Name,
				referenceName,
			)
		}
	}
}

// --- GeneratePrivateKey ---

func Test_GeneratePrivateKey(t *testing.T) {
	implementations := getX509Implementations()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_returns non nil private key", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			key, err := impl.GeneratePrivateKey(ctx)
			require.NoError(t, err)
			require.NotNil(t, key)
		})

		t.Run(impl.Name+"_generated keys are unique", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			key1, err := impl.GeneratePrivateKey(ctx)
			require.NoError(t, err)

			key2, err := impl.GeneratePrivateKey(ctx)
			require.NoError(t, err)

			isEqual, err := cryptoutils.IsPrivateKeyEqual(key1, key2)
			require.NoError(t, err)
			require.False(t, isEqual, "Two generated private keys should not be equal")
		})

		t.Run(impl.Name+"_generated key can extract public key", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			key, err := impl.GeneratePrivateKey(ctx)
			require.NoError(t, err)

			pubKey, err := cryptoutils.GetPublicKeyFromPrivateKey(key)
			require.NoError(t, err)
			require.NotNil(t, pubKey)
		})

		t.Run(impl.Name+"_generated key can be PEM encoded", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			key, err := impl.GeneratePrivateKey(ctx)
			require.NoError(t, err)

			pemStr, err := cryptoutils.EncodePrivateKeyAsPEMString(key)
			require.NoError(t, err)
			require.NotEmpty(t, pemStr)
			require.Contains(t, pemStr, "PRIVATE KEY")
		})
	}
}

// --- CreateRootCaCertificate ---

func Test_CreateRootCaCertificate(t *testing.T) {
	implementations := getX509Implementations()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_nil options returns error", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			pair, err := impl.CreateRootCaCertificate(ctx, nil)
			require.Error(t, err)
			require.Nil(t, pair)
		})

		t.Run(impl.Name+"_returns valid root CA cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultRootCaOptions()

			pair, err := impl.CreateRootCaCertificate(ctx, options)
			require.NoError(t, err)
			require.NotNil(t, pair)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)
			require.NotNil(t, cert)
		})

		t.Run(impl.Name+"_root CA has correct subject fields", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultRootCaOptions()

			pair, err := impl.CreateRootCaCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, "CH", cert.Subject.Country[0])
			require.Equal(t, "Zurich", cert.Subject.Locality[0])
			require.Equal(t, "TestRootOrg", cert.Subject.Organization[0])
			require.Equal(t, "TestRootCA", cert.Subject.CommonName)
		})

		t.Run(impl.Name+"_root CA is a CA certificate", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultRootCaOptions()

			pair, err := impl.CreateRootCaCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			require.True(t, cert.IsCA)
		})

		t.Run(impl.Name+"_root CA is self signed", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultRootCaOptions()

			pair, err := impl.CreateRootCaCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, cert.Subject.String(), cert.Issuer.String())
		})

		t.Run(impl.Name+"_root CA is detected as root CA by genericx509utils", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultRootCaOptions()

			pair, err := impl.CreateRootCaCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
			require.NoError(t, err)
			require.True(t, isRootCa)
		})

		t.Run(impl.Name+"_root CA key matches cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultRootCaOptions()

			pair, err := impl.CreateRootCaCertificate(ctx, options)
			require.NoError(t, err)

			err = pair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
		})

		t.Run(impl.Name+"_root CA cert can verify its own signature", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultRootCaOptions()

			pair, err := impl.CreateRootCaCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			isSigned, err := genericx509utils.IsSignedBy(ctx, cert, cert)
			require.NoError(t, err)
			require.True(t, isSigned)
		})
	}
}

// --- CreateIntermediateCertificate ---

func Test_CreateIntermediateCertificate(t *testing.T) {
	implementations := getX509Implementations()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_nil options returns error", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			pair, err := impl.CreateIntermediateCertificate(ctx, nil)
			require.Error(t, err)
			require.Nil(t, pair)
		})

		t.Run(impl.Name+"_returns valid intermediate cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultIntermediateOptions()

			pair, err := impl.CreateIntermediateCertificate(ctx, options)
			require.NoError(t, err)
			require.NotNil(t, pair)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)
			require.NotNil(t, cert)
		})

		t.Run(impl.Name+"_intermediate has correct subject fields", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultIntermediateOptions()

			pair, err := impl.CreateIntermediateCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, "CH", cert.Subject.Country[0])
			require.Equal(t, "Bern", cert.Subject.Locality[0])
			require.Equal(t, "TestIntOrg", cert.Subject.Organization[0])
			require.Equal(t, "TestIntermediateCA", cert.Subject.CommonName)
		})

		t.Run(impl.Name+"_intermediate is a CA certificate", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultIntermediateOptions()

			pair, err := impl.CreateIntermediateCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			require.True(t, cert.IsCA)
		})

		t.Run(impl.Name+"_intermediate key matches cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultIntermediateOptions()

			pair, err := impl.CreateIntermediateCertificate(ctx, options)
			require.NoError(t, err)

			err = pair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
		})
	}
}

// --- CreateSelfSignedCertificate ---

func Test_CreateSelfSignedCertificate(t *testing.T) {
	implementations := getX509Implementations()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_nil options returns error", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			pair, err := impl.CreateSelfSignedCertificate(ctx, nil)
			require.Error(t, err)
			require.Nil(t, pair)
		})

		t.Run(impl.Name+"_returns valid self signed cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultSelfSignedOptions()

			pair, err := impl.CreateSelfSignedCertificate(ctx, options)
			require.NoError(t, err)
			require.NotNil(t, pair)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)
			require.NotNil(t, cert)
		})

		t.Run(impl.Name+"_self signed has correct subject fields", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultSelfSignedOptions()

			pair, err := impl.CreateSelfSignedCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, "DE", cert.Subject.Country[0])
			require.Equal(t, "Berlin", cert.Subject.Locality[0])
			require.Equal(t, "SelfSignedOrg", cert.Subject.Organization[0])
			require.Equal(t, "selfsigned.example.com", cert.Subject.CommonName)
		})

		t.Run(impl.Name+"_self signed is self signed", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultSelfSignedOptions()

			pair, err := impl.CreateSelfSignedCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, cert.Subject.String(), cert.Issuer.String())
		})

		t.Run(impl.Name+"_self signed key matches cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultSelfSignedOptions()

			pair, err := impl.CreateSelfSignedCertificate(ctx, options)
			require.NoError(t, err)

			err = pair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
		})

		t.Run(impl.Name+"_self signed is not a CA", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()
			options := getDefaultSelfSignedOptions()

			pair, err := impl.CreateSelfSignedCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			require.False(t, cert.IsCA)
		})
	}
}

// --- CreateSignedIntermediateCertificate ---

func Test_CreateSignedIntermediateCertificate(t *testing.T) {
	implementations := getX509Implementations()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_nil options returns error", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			pair, err := impl.CreateSignedIntermediateCertificate(ctx, nil, rootPair)
			require.Error(t, err)
			require.Nil(t, pair)
		})

		t.Run(impl.Name+"_nil root CA returns error", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			pair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), nil)
			require.Error(t, err)
			require.Nil(t, pair)
		})

		t.Run(impl.Name+"_returns valid signed intermediate cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)
			require.NotNil(t, intPair)

			cert, err := intPair.GetX509Certificate()
			require.NoError(t, err)
			require.NotNil(t, cert)
		})

		t.Run(impl.Name+"_signed intermediate has correct subject fields", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			cert, err := intPair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, "CH", cert.Subject.Country[0])
			require.Equal(t, "Bern", cert.Subject.Locality[0])
			require.Equal(t, "TestIntOrg", cert.Subject.Organization[0])
			require.Equal(t, "TestIntermediateCA", cert.Subject.CommonName)
		})

		t.Run(impl.Name+"_signed intermediate is a CA certificate", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			cert, err := intPair.GetX509Certificate()
			require.NoError(t, err)

			require.True(t, cert.IsCA)
		})

		t.Run(impl.Name+"_signed intermediate is not self signed", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			cert, err := intPair.GetX509Certificate()
			require.NoError(t, err)

			require.NotEqual(t, cert.Subject.String(), cert.Issuer.String())
		})

		t.Run(impl.Name+"_signed intermediate is detected as intermediate by genericx509utils", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			cert, err := intPair.GetX509Certificate()
			require.NoError(t, err)

			isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, cert)
			require.NoError(t, err)
			require.True(t, isIntermediate)
		})

		t.Run(impl.Name+"_signed intermediate is signed by root CA", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			rootCert, err := rootPair.GetX509Certificate()
			require.NoError(t, err)

			intCert, err := intPair.GetX509Certificate()
			require.NoError(t, err)

			isSigned, err := genericx509utils.IsSignedBy(ctx, intCert, rootCert)
			require.NoError(t, err)
			require.True(t, isSigned)
		})

		t.Run(impl.Name+"_signed intermediate key matches cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			err = intPair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
		})

		t.Run(impl.Name+"_signed intermediate issuer matches root CA subject", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			rootCert, err := rootPair.GetX509Certificate()
			require.NoError(t, err)

			intCert, err := intPair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, rootCert.Subject.String(), intCert.Issuer.String())
		})
	}
}

// --- CreateSignedEndEntityCertificate ---

func Test_CreateSignedEndEntityCertificate(t *testing.T) {
	implementations := getX509Implementations()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_nil options returns error", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			pair, err := impl.CreateSignedEndEntityCertificate(ctx, nil, rootPair)
			require.Error(t, err)
			require.Nil(t, pair)
		})

		t.Run(impl.Name+"_nil CA returns error", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			pair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), nil)
			require.Error(t, err)
			require.Nil(t, pair)
		})

		t.Run(impl.Name+"_returns valid signed end entity cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
			require.NoError(t, err)
			require.NotNil(t, eePair)

			cert, err := eePair.GetX509Certificate()
			require.NoError(t, err)
			require.NotNil(t, cert)
		})

		t.Run(impl.Name+"_end entity has correct subject fields", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
			require.NoError(t, err)

			cert, err := eePair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, "CH", cert.Subject.Country[0])
			require.Equal(t, "Basel", cert.Subject.Locality[0])
			require.Equal(t, "TestEEOrg", cert.Subject.Organization[0])
			require.Equal(t, "server.example.com", cert.Subject.CommonName)
		})

		t.Run(impl.Name+"_end entity is not a CA certificate", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
			require.NoError(t, err)

			cert, err := eePair.GetX509Certificate()
			require.NoError(t, err)

			require.False(t, cert.IsCA)
		})

		t.Run(impl.Name+"_end entity is detected as end entity by genericx509utils", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
			require.NoError(t, err)

			cert, err := eePair.GetX509Certificate()
			require.NoError(t, err)

			isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, cert)
			require.NoError(t, err)
			require.True(t, isEndEntity)
		})

		t.Run(impl.Name+"_end entity is signed by CA", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
			require.NoError(t, err)

			rootCert, err := rootPair.GetX509Certificate()
			require.NoError(t, err)

			eeCert, err := eePair.GetX509Certificate()
			require.NoError(t, err)

			isSigned, err := genericx509utils.IsSignedBy(ctx, eeCert, rootCert)
			require.NoError(t, err)
			require.True(t, isSigned)
		})

		t.Run(impl.Name+"_end entity key matches cert", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
			require.NoError(t, err)

			err = eePair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
		})

		t.Run(impl.Name+"_end entity issuer matches CA subject", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
			require.NoError(t, err)

			rootCert, err := rootPair.GetX509Certificate()
			require.NoError(t, err)

			eeCert, err := eePair.GetX509Certificate()
			require.NoError(t, err)

			require.Equal(t, rootCert.Subject.String(), eeCert.Issuer.String())
		})

		t.Run(impl.Name+"_end entity signed by intermediate forms valid chain", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), intPair)
			require.NoError(t, err)

			rootCert, err := rootPair.GetX509Certificate()
			require.NoError(t, err)
			intCert, err := intPair.GetX509Certificate()
			require.NoError(t, err)
			eeCert, err := eePair.GetX509Certificate()
			require.NoError(t, err)

			isValidChain, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{eeCert, intCert, rootCert})
			require.NoError(t, err)
			require.True(t, isValidChain)
		})

		t.Run(impl.Name+"_end entity is not signed by unrelated CA", func(t *testing.T) {
			ctx := contextutils.ContextVerbose()

			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), rootPair)
			require.NoError(t, err)

			unrelatedRootOptions := &x509options.X509CreateCertificateOptions{
				CountryName:    "FR",
				Locality:       "Paris",
				Organization:   "UnrelatedOrg",
				CommonName:     "UnrelatedRoot",
				PrivateKeySize: 2048,
			}
			unrelatedPair, err := impl.CreateRootCaCertificate(ctx, unrelatedRootOptions)
			require.NoError(t, err)

			unrelatedCert, err := unrelatedPair.GetX509Certificate()
			require.NoError(t, err)

			eeCert, err := eePair.GetX509Certificate()
			require.NoError(t, err)

			isSigned, err := genericx509utils.IsSignedBy(ctx, eeCert, unrelatedCert)
			require.NoError(t, err)
			require.False(t, isSigned)
		})
	}
}

// --- Cross-implementation consistency ---

// Test_CreateRootCaCertificate_implementations_produce_consistent_results validates
// that all implementations produce root CA certificates with identical structural properties.
func Test_CreateRootCaCertificate_implementations_produce_consistent_results(t *testing.T) {
	ctx := contextutils.ContextVerbose()
	implementations := getX509Implementations()
	options := getDefaultRootCaOptions()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_produced cert has expected properties", func(t *testing.T) {
			pair, err := impl.CreateRootCaCertificate(ctx, options)
			require.NoError(t, err)

			cert, err := pair.GetX509Certificate()
			require.NoError(t, err)

			// All implementations should produce a cert with these properties:
			require.True(t, cert.IsCA)
			require.Equal(t, cert.Subject.String(), cert.Issuer.String())
			require.Equal(t, "CH", cert.Subject.Country[0])
			require.Equal(t, "Zurich", cert.Subject.Locality[0])
			require.Equal(t, "TestRootOrg", cert.Subject.Organization[0])
			require.Equal(t, "TestRootCA", cert.Subject.CommonName)

			isRootCa, err := genericx509utils.IsRootCaCert(ctx, cert)
			require.NoError(t, err)
			require.True(t, isRootCa)

			err = pair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
		})
	}
}

// --- Full chain integration test ---

// Test_FullChain_across_implementations validates that a full certificate chain
// (root -> intermediate -> end entity) works correctly for each implementation.
func Test_FullChain_across_implementations(t *testing.T) {
	ctx := contextutils.ContextVerbose()
	implementations := getX509Implementations()

	for _, impl := range implementations {
		impl := impl

		t.Run(impl.Name+"_full chain root-intermediate-end entity", func(t *testing.T) {
			// Create Root CA
			rootPair, err := impl.CreateRootCaCertificate(ctx, getDefaultRootCaOptions())
			require.NoError(t, err)

			// Create Intermediate signed by Root
			intPair, err := impl.CreateSignedIntermediateCertificate(ctx, getDefaultIntermediateOptions(), rootPair)
			require.NoError(t, err)

			// Create End Entity signed by Intermediate
			eePair, err := impl.CreateSignedEndEntityCertificate(ctx, getDefaultEndEntityOptions(), intPair)
			require.NoError(t, err)

			// Extract certs
			rootCert, err := rootPair.GetX509Certificate()
			require.NoError(t, err)
			intCert, err := intPair.GetX509Certificate()
			require.NoError(t, err)
			eeCert, err := eePair.GetX509Certificate()
			require.NoError(t, err)

			// Validate types
			isRootCa, err := genericx509utils.IsRootCaCert(ctx, rootCert)
			require.NoError(t, err)
			require.True(t, isRootCa)

			isIntermediate, err := genericx509utils.IsIntermediateCert(ctx, intCert)
			require.NoError(t, err)
			require.True(t, isIntermediate)

			isEndEntity, err := genericx509utils.IsEndEntityCert(ctx, eeCert)
			require.NoError(t, err)
			require.True(t, isEndEntity)

			// Validate signing relationships
			isSigned, err := genericx509utils.IsSignedBy(ctx, intCert, rootCert)
			require.NoError(t, err)
			require.True(t, isSigned)

			isSigned, err = genericx509utils.IsSignedBy(ctx, eeCert, intCert)
			require.NoError(t, err)
			require.True(t, isSigned)

			// Validate full chain
			isValidChain, err := genericx509utils.IsRootCaToEndEntityChain(ctx, []*x509.Certificate{eeCert, intCert, rootCert})
			require.NoError(t, err)
			require.True(t, isValidChain)

			// Validate signing chain
			isChain, err := genericx509utils.IsCertChain(ctx, []*x509.Certificate{eeCert, intCert, rootCert})
			require.NoError(t, err)
			require.True(t, isChain)

			// Validate all keys match their certs
			err = rootPair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
			err = intPair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
			err = eePair.CheckKeyMatchingCertificate()
			require.NoError(t, err)
		})
	}
}
