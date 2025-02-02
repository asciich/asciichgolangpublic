package x509utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getX509CertificateHandlerToTest(implementationName string) (handler X509CertificateHandler) {
	if implementationName == "NativeX509CertificateHandler" {
		return GetNativeX509CertificateHandler()
	}

	logging.LogFatalWithTracef("Unknown implementationName '%s'", implementationName)

	return
}

func TestX509Handler_CreateRootCaCertificate(t *testing.T) {
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

				cert, privateKey := mustutils.Must2(handler.CreateRootCaCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",

						Verbose: true,
					},
				))

				assert.True(mustutils.Must(IsCertificateRootCa(cert)))
				assert.False(mustutils.Must(IsIntermediateCertificate(cert)))
				assert.False(mustutils.Must(IsEndEndityCertificate(cert)))

				assert.EqualValues([]string{"CH"}, cert.Issuer.Country)
				assert.EqualValues([]string{"Zurich"}, cert.Issuer.Locality)

				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(cert, privateKey)))
			},
		)
	}
}

func TestX509Handler_CreateIntermediateCertificate(t *testing.T) {
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

				caCert, caPrivateKey := mustutils.Must2(handler.CreateRootCaCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",

						Verbose: true,
					},
				))

				intermediateCert, intermediateKey := mustutils.Must2(handler.CreateSignedIntermediateCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg intermediate",

						Verbose: true,
					},
					caCert,
					caPrivateKey,
					true,
				))

				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(caCert, caPrivateKey)))
				assert.True(mustutils.Must(IsCertificateRootCa(caCert)))
				assert.False(mustutils.Must(IsIntermediateCertificate(caCert)))
				assert.False(mustutils.Must(IsEndEndityCertificate(caCert)))

				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCert, intermediateKey)))
				assert.False(mustutils.Must(IsCertificateRootCa(intermediateCert)))
				assert.True(mustutils.Must(IsIntermediateCertificate(intermediateCert)))
				assert.False(mustutils.Must(IsEndEndityCertificate(intermediateCert)))
			},
		)
	}
}
