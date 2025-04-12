package x509utils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/datatypes/bigintutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/cryptoutils"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

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
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",
					},
				)
				require.NoError(t, err)

				anotherKey, err := handler.GeneratePrivateKey(getCtx())
				require.NoError(t, err)

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(rootCaCertAndKey.Cert, rootCaCertAndKey.Key)))
				require.False(t, mustutils.Must(IsCertificateMatchingPrivateKey(rootCaCertAndKey.Cert, anotherKey)))
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
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",
					},
				)
				require.NoError(t, err)

				intermediateCertAndKey, err := handler.CreateSignedIntermediateCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "IntermediateOrg",
					},
					rootCaCertAndKey,
				)
				require.NoError(t, err)

				selfSignedCertAndKey, err := handler.CreateSelfSignedCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "SelfSignedOrg",
					},
				)
				require.NoError(t, err)

				ctx := contextutils.ContextVerbose()

				require.True(t, mustutils.Must(IsCertSignedBy(ctx, intermediateCertAndKey.Cert, rootCaCertAndKey.Cert)))
				require.False(t, mustutils.Must(IsCertSignedBy(ctx, intermediateCertAndKey.Cert, selfSignedCertAndKey.Cert)))
				require.False(t, mustutils.Must(IsCertSignedBy(ctx, rootCaCertAndKey.Cert, intermediateCertAndKey.Cert)))
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
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",
					},
				)
				require.NoError(t, err)

				derEncoded := mustutils.Must(EncodeCertificateAsDerBytes(rootCaCertAndKey.Cert))
				cert2 := mustutils.Must(LoadCertificateFromDerBytes(derEncoded))

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
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",
					},
				)
				require.NoError(t, err)

				derEncoded := mustutils.Must(EncodeCertificateAsPEMString(rootCaCertAndKey.Cert))
				require.True(t, strings.HasPrefix(derEncoded, "-----BEGIN CERTIFICATE-----\n"))
				require.True(t, strings.HasSuffix(derEncoded, "\n-----END CERTIFICATE-----\n"))

				cert2 := mustutils.Must(LoadCertificateFromPEMString(derEncoded))

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
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",
					},
				)
				require.NoError(t, err)

				derEncoded := mustutils.Must(cryptoutils.EncodePrivateKeyAsPEMString(rootCaCertAndKey.Key))
				require.True(t, strings.HasPrefix(derEncoded, "-----BEGIN PRIVATE KEY-----\n"))
				require.True(t, strings.HasSuffix(derEncoded, "\n-----END PRIVATE KEY-----\n"))

				key2 := mustutils.Must(cryptoutils.LoadPrivateKeyFromPEMString(derEncoded))

				require.True(t, mustutils.Must(IsPrivateKeyEqual(rootCaCertAndKey.Key, key2)))
			},
		)
	}
}

func Test_GetValidityDuration(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		vd, err := GetValidityDuration(nil)
		require.Error(t, err)
		require.Nil(t, vd)
	})

	t.Run("1day", func(t *testing.T) {
		start := time.Now()
		cert := &x509.Certificate{
			NotBefore: start,
			NotAfter:  start.Add(time.Hour * 24),
		}
		vd, err := GetValidityDuration(cert)
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
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",
					},
				)
				require.NoError(t, err)

				out, err := GetSubjectAndSerialString(rootCaCertAndKey.Cert)
				require.NoError(t, err)
				require.EqualValues(t, "O=RootOrg,L=Zurich,ST=,C=CH serial: 01", out)
			},
		)
	}
}

func Test_GenerateSerialNumber(t *testing.T) {
	generated, err := GenerateCertificateSerialNumber(getCtx())
	require.NoError(t, err)

	generatedStr, err := bigintutils.ToHexStringColonSeparated(generated)
	require.NoError(t, err)
	require.True(t, len(generatedStr) > 4)
}
