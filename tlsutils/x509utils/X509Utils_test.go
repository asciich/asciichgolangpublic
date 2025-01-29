package x509utils

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestX509Utils_GetPublicKeyFromPrivateKey(t *testing.T) {
	generatedKey := mustutils.Must(rsa.GenerateKey(rand.Reader, 4096))
	generatedKey2 := mustutils.Must(rsa.GenerateKey(rand.Reader, 4096))

	publicKey := mustutils.Must(GetPublicKeyFromPrivateKey(generatedKey))
	publicKey2 := mustutils.Must(GetPublicKeyFromPrivateKey(generatedKey2))

	assert.True(t, generatedKey.PublicKey.Equal(publicKey))
	assert.True(t, generatedKey2.PublicKey.Equal(publicKey2))

	assert.False(t, generatedKey.PublicKey.Equal(publicKey2))
	assert.False(t, generatedKey2.PublicKey.Equal(publicKey))
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
				assert := assert.New(t)

				handler := getX509CertificateHandlerToTest(tt.implementationName)
				cert, realKey := mustutils.Must2(handler.CreateRootCaCertificate(
					&X509CreateCertificateOptions{
						CountryName: "CH",
						Locality:    "Zurich",
						Organization: "RootOrg",

						Verbose: true,
					},
				))

				anotherKey := mustutils.Must(rsa.GenerateKey(rand.Reader, 4096))

				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(cert, realKey)))
				assert.False(mustutils.Must(IsCertificateMatchingPrivateKey(cert, anotherKey)))
			},
		)
	}
}
