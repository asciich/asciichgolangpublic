package x509utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestX509Utils_GetPublicKeyFromPrivateKey(t *testing.T) {
	generatedKey := mustutils.Must(rsa.GenerateKey(rand.Reader, 4096))
	generatedKey2 := mustutils.Must(rsa.GenerateKey(rand.Reader, 4096))

	publicKey := mustutils.Must(GetPublicKeyFromPrivateKey(generatedKey))
	publicKey2 := mustutils.Must(GetPublicKeyFromPrivateKey(generatedKey2))

	require.True(t, generatedKey.PublicKey.Equal(publicKey))
	require.True(t, generatedKey2.PublicKey.Equal(publicKey2))

	require.False(t, generatedKey.PublicKey.Equal(publicKey2))
	require.False(t, generatedKey2.PublicKey.Equal(publicKey))
}

func TestX509Utils_IsCertificateMatchingPrivateKey(t *testing.T) {
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
				require := require.New(t)

				handler := getX509CertificateHandlerToTest(tt.implementationName)
				cert, realKey := mustutils.Must2(handler.CreateRootCaCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",

						Verbose: true,
					},
				))

				anotherKey := mustutils.Must(handler.GeneratePrivateKey())

				require.True(mustutils.Must(IsCertificateMatchingPrivateKey(cert, realKey)))
				require.False(mustutils.Must(IsCertificateMatchingPrivateKey(cert, anotherKey)))
			},
		)
	}
}

func TestX509Utils_EndcodeAndDecodeAsDER(t *testing.T) {
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
				require := require.New(t)

				handler := getX509CertificateHandlerToTest(tt.implementationName)
				cert, _ := mustutils.Must2(handler.CreateRootCaCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",

						Verbose: true,
					},
				))

				derEncoded := mustutils.Must(EncodeCertificateAsDerBytes(cert))
				cert2 := mustutils.Must(LoadCertificateFromDerBytes(derEncoded))

				require.True(cert.Equal(cert2))
			},
		)
	}
}

func TestX509Utils_EndcodeAndDecodeAsPEM(t *testing.T) {
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
				require := require.New(t)

				handler := getX509CertificateHandlerToTest(tt.implementationName)
				cert, _ := mustutils.Must2(handler.CreateRootCaCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",

						Verbose: true,
					},
				))

				derEncoded := mustutils.Must(EncodeCertificateAsPEMString(cert))
				require.True(strings.HasPrefix(derEncoded, "-----BEGIN CERTIFICATE-----\n"))
				require.True(strings.HasSuffix(derEncoded, "\n-----END CERTIFICATE-----\n"))

				cert2 := mustutils.Must(LoadCertificateFromPEMString(derEncoded))

				require.True(cert.Equal(cert2))
			},
		)
	}
}

func TestX509Utils_EndcodeAndDecodePrivateKeyAsPem(t *testing.T) {
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
				require := require.New(t)

				handler := getX509CertificateHandlerToTest(tt.implementationName)
				_, key := mustutils.Must2(handler.CreateRootCaCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "RootOrg",

						Verbose: true,
					},
				))

				derEncoded := mustutils.Must(EncodePrivateKeyAsPEMString(key))
				require.True(strings.HasPrefix(derEncoded, "-----BEGIN PRIVATE KEY-----\n"))
				require.True(strings.HasSuffix(derEncoded, "\n-----END PRIVATE KEY-----\n"))

				key2 := mustutils.Must(LoadPrivateKeyFromPEMString(derEncoded))

				require.True(mustutils.Must(IsPrivateKeyEqual(key, key2)))
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
