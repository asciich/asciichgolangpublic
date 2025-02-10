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

				caCert, privateKey := mustutils.Must2(handler.CreateRootCaCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",
						SerialNumber: "12345",

						Verbose: true,
					},
				))

				assert.True(mustutils.Must(IsCertificateRootCa(caCert)))
				assert.False(mustutils.Must(IsIntermediateCertificate(caCert)))
				assert.False(mustutils.Must(IsEndEndityCertificate(caCert)))
				assert.True(mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				assert.True(mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				assert.True(mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))
				assert.True(mustutils.Must(IsSerialNumber(caCert, "12345")))

				assert.EqualValues([]string{"CH"}, caCert.Issuer.Country)
				assert.EqualValues([]string{"Zurich"}, caCert.Issuer.Locality)

				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(caCert, privateKey)))
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
				assert.True(mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				assert.True(mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				assert.True(mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))

				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCert, intermediateKey)))
				assert.False(mustutils.Must(IsCertificateRootCa(intermediateCert)))
				assert.True(mustutils.Must(IsIntermediateCertificate(intermediateCert)))
				assert.False(mustutils.Must(IsEndEndityCertificate(intermediateCert)))
				assert.True(mustutils.Must(IsSubjectCountryName(intermediateCert, "CH")))
				assert.True(mustutils.Must(IsSubjectLocalityName(intermediateCert, "Zurich")))
				assert.True(mustutils.Must(IsSubjectOrganizationName(intermediateCert, "myOrg intermediate")))
			},
		)
	}
}

func TestX509Handler_CreateEndEndityCertificate(t *testing.T) {
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

				endEndityCertificate, endEndityKey := mustutils.Must2(handler.CreateSignedEndEndityCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg endEndity",

						Verbose: true,
					},
					intermediateCert,
					intermediateKey,
					true,
				))

				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(caCert, caPrivateKey)))
				assert.True(mustutils.Must(IsCertificateRootCa(caCert)))
				assert.False(mustutils.Must(IsIntermediateCertificate(caCert)))
				assert.False(mustutils.Must(IsEndEndityCertificate(caCert)))
				assert.True(mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				assert.True(mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				assert.True(mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))

				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCert, intermediateKey)))
				assert.False(mustutils.Must(IsCertificateRootCa(intermediateCert)))
				assert.True(mustutils.Must(IsIntermediateCertificate(intermediateCert)))
				assert.False(mustutils.Must(IsEndEndityCertificate(intermediateCert)))
				assert.True(mustutils.Must(IsSubjectCountryName(intermediateCert, "CH")))
				assert.True(mustutils.Must(IsSubjectLocalityName(intermediateCert, "Zurich")))
				assert.True(mustutils.Must(IsSubjectOrganizationName(intermediateCert, "myOrg intermediate")))


				assert.True(mustutils.Must(IsCertificateMatchingPrivateKey(endEndityCertificate, endEndityKey)))
				assert.False(mustutils.Must(IsCertificateRootCa(endEndityCertificate)))
				assert.False(mustutils.Must(IsIntermediateCertificate(endEndityCertificate)))
				assert.True(mustutils.Must(IsEndEndityCertificate(endEndityCertificate)))
				assert.True(mustutils.Must(IsSubjectCountryName(endEndityCertificate, "CH")))
				assert.True(mustutils.Must(IsSubjectLocalityName(endEndityCertificate, "Zurich")))
				assert.True(mustutils.Must(IsSubjectOrganizationName(endEndityCertificate, "myOrg endEndity")))
			},
		)
	}
}
