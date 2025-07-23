package x509utils_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/bigintutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/contextutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/cryptoutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/tlsutils/x509utils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
)

func Test_GetPublicKeyFromPrivateKey(t *testing.T) {
	generatedKey := mustutils.Must(rsa.GenerateKey(rand.Reader, 4096))
	generatedKey2 := mustutils.Must(rsa.GenerateKey(rand.Reader, 4096))

	publicKey := mustutils.Must(cryptoutils.GetPublicKeyFromPrivateKey(generatedKey))
	publicKey2 := mustutils.Must(cryptoutils.GetPublicKeyFromPrivateKey(generatedKey2))

	require.True(t, generatedKey.PublicKey.Equal(publicKey))
	require.True(t, generatedKey2.PublicKey.Equal(publicKey2))

	require.False(t, generatedKey.PublicKey.Equal(publicKey2))
	require.False(t, generatedKey2.PublicKey.Equal(publicKey))
}

func Test_IsCertificateMatchingPrivateKey(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"NativeX509CertificateHandler"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				handler := getX509CertificateHandlerToTest(tt.implementationName)
				rootCaCertAndKey, err := handler.CreateRootCaCertificate(
					getCtx(),
					&x509utils.X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",
					},
				)
				require.NoError(t, err)

				anotherKey, err := handler.GeneratePrivateKey(getCtx())
				require.NoError(t, err)

				require.True(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(rootCaCertAndKey.Cert, rootCaCertAndKey.Key)))
				require.False(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(rootCaCertAndKey.Cert, anotherKey)))
			},
		)
	}
}

func Test_IsCertificateSignedBy(t *testing.T) {
	ctx := getCtx()

	tests := []struct {
		implementationName string
	}{
		{"NativeX509CertificateHandler"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				handler := getX509CertificateHandlerToTest(tt.implementationName)
				rootCaCertAndKey, err := handler.CreateRootCaCertificate(
					ctx,
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "RootOrg",
						PrivateKeySize: 1024,
					},
				)
				require.NoError(t, err)

				intermediateCertAndKey, err := handler.CreateSignedIntermediateCertificate(
					ctx,
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "IntermediateOrg",
						PrivateKeySize: 1024,
					},
					rootCaCertAndKey,
				)
				require.NoError(t, err)

				selfSignedCertAndKey, err := handler.CreateSelfSignedCertificate(
					ctx,
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "SelfSignedOrg",
						PrivateKeySize: 1024,
					},
				)
				require.NoError(t, err)

				ctx := contextutils.ContextVerbose()

				require.True(t, mustutils.Must(x509utils.IsCertSignedBy(ctx, intermediateCertAndKey.Cert, rootCaCertAndKey.Cert)))
				require.False(t, mustutils.Must(x509utils.IsCertSignedBy(ctx, intermediateCertAndKey.Cert, selfSignedCertAndKey.Cert)))
				require.False(t, mustutils.Must(x509utils.IsCertSignedBy(ctx, rootCaCertAndKey.Cert, intermediateCertAndKey.Cert)))
			},
		)
	}

}

func Test_EndcodeAndDecodeAsDER(t *testing.T) {
	ctx := getCtx()

	tests := []struct {
		implementationName string
	}{
		{"NativeX509CertificateHandler"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				handler := getX509CertificateHandlerToTest(tt.implementationName)
				rootCaCertAndKey, err := handler.CreateRootCaCertificate(
					ctx,
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "RootOrg",
						PrivateKeySize: 1024,
					},
				)
				require.NoError(t, err)

				derEncoded := mustutils.Must(x509utils.EncodeCertificateAsDerBytes(rootCaCertAndKey.Cert))
				cert2 := mustutils.Must(x509utils.LoadCertificateFromDerBytes(derEncoded))

				require.True(t, rootCaCertAndKey.Cert.Equal(cert2))
			},
		)
	}
}

func Test_EndcodeAndDecodeAsPEM(t *testing.T) {
	ctx := getCtx()

	tests := []struct {
		implementationName string
	}{
		{"NativeX509CertificateHandler"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				handler := getX509CertificateHandlerToTest(tt.implementationName)
				rootCaCertAndKey, err := handler.CreateRootCaCertificate(
					ctx,
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "RootOrg",
						PrivateKeySize: 1024,
					},
				)
				require.NoError(t, err)

				derEncoded := mustutils.Must(x509utils.EncodeCertificateAsPEMString(rootCaCertAndKey.Cert))
				require.True(t, strings.HasPrefix(derEncoded, "-----BEGIN CERTIFICATE-----\n"))
				require.True(t, strings.HasSuffix(derEncoded, "\n-----END CERTIFICATE-----\n"))

				cert2 := mustutils.Must(x509utils.LoadCertificateFromPEMString(derEncoded))

				require.True(t, rootCaCertAndKey.Cert.Equal(cert2))
			},
		)
	}
}

func Test_EndcodeAndDecodePrivateKeyAsPem(t *testing.T) {
	ctx := getCtx()

	tests := []struct {
		implementationName string
	}{
		{"NativeX509CertificateHandler"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				handler := getX509CertificateHandlerToTest(tt.implementationName)
				rootCaCertAndKey, err := handler.CreateRootCaCertificate(
					ctx,
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "RootOrg",
						PrivateKeySize: 1024,
					},
				)
				require.NoError(t, err)

				derEncoded := mustutils.Must(cryptoutils.EncodePrivateKeyAsPEMString(rootCaCertAndKey.Key))
				require.True(t, strings.HasPrefix(derEncoded, "-----BEGIN PRIVATE KEY-----\n"))
				require.True(t, strings.HasSuffix(derEncoded, "\n-----END PRIVATE KEY-----\n"))

				key2 := mustutils.Must(cryptoutils.LoadPrivateKeyFromPEMString(derEncoded))

				require.True(t, mustutils.Must(x509utils.IsPrivateKeyEqual(rootCaCertAndKey.Key, key2)))
			},
		)
	}
}

func Test_GetValidityDuration(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		vd, err := x509utils.GetValidityDuration(nil)
		require.Error(t, err)
		require.Nil(t, vd)
	})

	t.Run("1day", func(t *testing.T) {
		start := time.Now()
		cert := &x509.Certificate{
			NotBefore: start,
			NotAfter:  start.Add(time.Hour * 24),
		}
		vd, err := x509utils.GetValidityDuration(cert)
		require.NoError(t, err)
		require.EqualValues(t, time.Hour*24, *vd)
	})

}

func Test_GetSubjectAndSerialString(t *testing.T) {
	ctx := getCtx()

	tests := []struct {
		implementationName string
	}{
		{"NativeX509CertificateHandler"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				handler := getX509CertificateHandlerToTest(tt.implementationName)
				rootCaCertAndKey, err := handler.CreateRootCaCertificate(
					ctx,
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "RootOrg",
						SerialNumber:   "1",
						PrivateKeySize: 1024,
					},
				)
				require.NoError(t, err)

				out, err := x509utils.GetSubjectAndSerialString(rootCaCertAndKey.Cert)
				require.NoError(t, err)
				require.EqualValues(t, "O=RootOrg,L=Zurich,ST=,C=CH serial: 01", out)
			},
		)
	}
}

func Test_GenerateSerialNumber(t *testing.T) {
	generated, err := x509utils.GenerateCertificateSerialNumber(getCtx())
	require.NoError(t, err)

	generatedStr, err := bigintutils.ToHexStringColonSeparated(generated)
	require.NoError(t, err)
	require.True(t, len(generatedStr) > 4)
}

func Test_ValidateCertificateChain(t *testing.T) {
	ctx := getCtx()

	handler := x509utils.GetDefaultHandler()

	rootCeCertAndKey, err := handler.CreateRootCaCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg root",
			PrivateKeySize: 1024,
		},
	)
	require.NoError(t, err)

	intermediateCertAndKey, err := handler.CreateSignedIntermediateCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg intermediate",
			PrivateKeySize: 1024,
		},
		rootCeCertAndKey,
	)
	require.NoError(t, err)

	endEndityCertAndKey, err := handler.CreateSignedEndEndityCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CommonName:     "mytestcn.example.net",
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg endEndity",
			AdditionalSans: []string{"mytestsan1.example.net", "mytestsan2.example.net"},
			PrivateKeySize: 1024,
		},
		intermediateCertAndKey,
	)
	require.NoError(t, err)

	rootCeCertAndKey2, err := handler.CreateRootCaCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg root2",
			PrivateKeySize: 1024,
		},
	)
	require.NoError(t, err)

	intermediateCertAndKey2, err := handler.CreateSignedIntermediateCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg intermediate2",
			PrivateKeySize: 1024,
		},
		rootCeCertAndKey2,
	)
	require.NoError(t, err)

	endEndityCertAndKey2, err := handler.CreateSignedEndEndityCertificate(
		ctx,
		&x509utils.X509CreateCertificateOptions{
			CommonName:     "mytestcn.example.net",
			CountryName:    "CH",
			Locality:       "Zurich",
			Organization:   "myOrg endEndity2",
			AdditionalSans: []string{"mytestsan21.example.net", "mytestsan22.example.net"},
			PrivateKeySize: 1024,
		},
		intermediateCertAndKey2,
	)
	require.NoError(t, err)

	t.Run("all nil", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(ctx, nil, nil, nil)
		require.Error(t, err)
		require.Nil(t, chains)
	})

	t.Run("all nil but certToValidate given", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(ctx, endEndityCertAndKey2.Cert, nil, nil)
		require.Error(t, err)
		require.Nil(t, chains)
	})

	t.Run("all nil but certToValidate and cert and trusted given", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(ctx, endEndityCertAndKey2.Cert, []*x509.Certificate{rootCeCertAndKey.Cert}, nil)
		require.Error(t, err)
		require.Nil(t, chains)
	})

	t.Run("valid chain 1", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(
			ctx,
			endEndityCertAndKey.Cert,
			[]*x509.Certificate{rootCeCertAndKey.Cert},
			[]*x509.Certificate{intermediateCertAndKey.Cert})
		require.NoError(t, err)
		require.NotNil(t, chains)
		require.Len(t, chains, 1)
		require.EqualValues(t, chains[0][0], endEndityCertAndKey.Cert)
		require.EqualValues(t, chains[0][1], intermediateCertAndKey.Cert)
		require.EqualValues(t, chains[0][2], rootCeCertAndKey.Cert)
	})

	t.Run("valid chain 2", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(
			ctx,
			endEndityCertAndKey2.Cert,
			[]*x509.Certificate{rootCeCertAndKey2.Cert},
			[]*x509.Certificate{intermediateCertAndKey2.Cert})
		require.NoError(t, err)
		require.NotNil(t, chains)
		require.Len(t, chains, 1)
		require.EqualValues(t, chains[0][0], endEndityCertAndKey2.Cert)
		require.EqualValues(t, chains[0][1], intermediateCertAndKey2.Cert)
		require.EqualValues(t, chains[0][2], rootCeCertAndKey2.Cert)
	})

	t.Run("two root CAs", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(
			ctx,
			endEndityCertAndKey2.Cert,
			[]*x509.Certificate{rootCeCertAndKey2.Cert, rootCeCertAndKey.Cert},
			[]*x509.Certificate{intermediateCertAndKey2.Cert})
		require.NoError(t, err)
		require.NotNil(t, chains)
		require.Len(t, chains, 1)
		require.EqualValues(t, chains[0][0], endEndityCertAndKey2.Cert)
		require.EqualValues(t, chains[0][1], intermediateCertAndKey2.Cert)
		require.EqualValues(t, chains[0][2], rootCeCertAndKey2.Cert)
	})

	t.Run("two root CAs and intermediates", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(
			ctx,
			endEndityCertAndKey2.Cert,
			[]*x509.Certificate{rootCeCertAndKey2.Cert, rootCeCertAndKey.Cert},
			[]*x509.Certificate{intermediateCertAndKey.Cert, intermediateCertAndKey2.Cert})
		require.NoError(t, err)
		require.NotNil(t, chains)
		require.Len(t, chains, 1)
		require.EqualValues(t, chains[0][0], endEndityCertAndKey2.Cert)
		require.EqualValues(t, chains[0][1], intermediateCertAndKey2.Cert)
		require.EqualValues(t, chains[0][2], rootCeCertAndKey2.Cert)
	})

	t.Run("Invalid root", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(
			ctx,
			endEndityCertAndKey2.Cert,
			[]*x509.Certificate{rootCeCertAndKey.Cert},
			[]*x509.Certificate{intermediateCertAndKey2.Cert},
		)
		require.ErrorIs(t, err, x509utils.ErrNoValidCertificateChain)
		require.Len(t, chains, 0)
		require.Nil(t, chains)
	})

	t.Run("Invalid intermediate", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(
			ctx,
			endEndityCertAndKey2.Cert,
			[]*x509.Certificate{rootCeCertAndKey2.Cert},
			[]*x509.Certificate{intermediateCertAndKey.Cert},
		)
		require.ErrorIs(t, err, x509utils.ErrNoValidCertificateChain)
		require.Len(t, chains, 0)
		require.Nil(t, chains)
	})

	t.Run("Invalid endEndity", func(t *testing.T) {
		chains, err := x509utils.ValidateCertificateChain(
			ctx,
			endEndityCertAndKey.Cert,
			[]*x509.Certificate{rootCeCertAndKey2.Cert},
			[]*x509.Certificate{intermediateCertAndKey2.Cert},
		)
		require.ErrorIs(t, err, x509utils.ErrNoValidCertificateChain)
		require.Len(t, chains, 0)
		require.Nil(t, chains)
	})
}
