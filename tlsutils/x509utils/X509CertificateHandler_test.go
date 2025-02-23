package x509utils

import (
	"testing"

	"github.com/stretchr/testify/require"
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
				require := require.New(t)

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

				require.True(mustutils.Must(IsCertificateRootCa(caCert)))
				require.True(mustutils.Must(IsSelfSignedCertificate(caCert)))
				require.False(mustutils.Must(IsIntermediateCertificate(caCert)))
				require.False(mustutils.Must(IsEndEndityCertificate(caCert)))
				require.True(mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				require.True(mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				require.True(mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))
				require.True(mustutils.Must(IsSerialNumber(caCert, "12345")))

				require.EqualValues([]string{"CH"}, caCert.Issuer.Country)
				require.EqualValues([]string{"Zurich"}, caCert.Issuer.Locality)

				require.True(mustutils.Must(IsCertificateMatchingPrivateKey(caCert, privateKey)))
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
				require := require.New(t)

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

				require.True(mustutils.Must(IsCertificateMatchingPrivateKey(caCert, caPrivateKey)))
				require.True(mustutils.Must(IsCertificateRootCa(caCert)))
				require.False(mustutils.Must(IsIntermediateCertificate(caCert)))
				require.False(mustutils.Must(IsEndEndityCertificate(caCert)))
				require.True(mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				require.True(mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				require.True(mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))

				require.True(mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCert, intermediateKey)))
				require.False(mustutils.Must(IsCertificateRootCa(intermediateCert)))
				require.True(mustutils.Must(IsIntermediateCertificate(intermediateCert)))
				require.False(mustutils.Must(IsEndEndityCertificate(intermediateCert)))
				require.True(mustutils.Must(IsSubjectCountryName(intermediateCert, "CH")))
				require.True(mustutils.Must(IsSubjectLocalityName(intermediateCert, "Zurich")))
				require.True(mustutils.Must(IsSubjectOrganizationName(intermediateCert, "myOrg intermediate")))
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
				require := require.New(t)

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

				require.True(mustutils.Must(IsCertificateMatchingPrivateKey(caCert, caPrivateKey)))
				require.True(mustutils.Must(IsCertificateRootCa(caCert)))
				require.False(mustutils.Must(IsIntermediateCertificate(caCert)))
				require.False(mustutils.Must(IsEndEndityCertificate(caCert)))
				require.True(mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				require.True(mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				require.True(mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))

				require.True(mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCert, intermediateKey)))
				require.False(mustutils.Must(IsCertificateRootCa(intermediateCert)))
				require.True(mustutils.Must(IsIntermediateCertificate(intermediateCert)))
				require.False(mustutils.Must(IsEndEndityCertificate(intermediateCert)))
				require.True(mustutils.Must(IsSubjectCountryName(intermediateCert, "CH")))
				require.True(mustutils.Must(IsSubjectLocalityName(intermediateCert, "Zurich")))
				require.True(mustutils.Must(IsSubjectOrganizationName(intermediateCert, "myOrg intermediate")))

				require.True(mustutils.Must(IsCertificateMatchingPrivateKey(endEndityCertificate, endEndityKey)))
				require.False(mustutils.Must(IsCertificateRootCa(endEndityCertificate)))
				require.False(mustutils.Must(IsIntermediateCertificate(endEndityCertificate)))
				require.True(mustutils.Must(IsEndEndityCertificate(endEndityCertificate)))
				require.True(mustutils.Must(IsSubjectCountryName(endEndityCertificate, "CH")))
				require.True(mustutils.Must(IsSubjectLocalityName(endEndityCertificate, "Zurich")))
				require.True(mustutils.Must(IsSubjectOrganizationName(endEndityCertificate, "myOrg endEndity")))
			},
		)
	}
}

func TestX509Handler_CreateSelfSignedCertificate(t *testing.T) {
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

				caCert, privateKey := mustutils.Must2(handler.CreateSelfSignedCertificate(
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",
						SerialNumber: "12345",

						Verbose: true,
					},
				))

				require.False(mustutils.Must(IsCertificateRootCa(caCert)))
				require.True(mustutils.Must(IsSelfSignedCertificate(caCert)))
				require.False(mustutils.Must(IsIntermediateCertificate(caCert)))
				require.True(mustutils.Must(IsEndEndityCertificate(caCert)))
				require.True(mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				require.True(mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				require.True(mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))
				require.True(mustutils.Must(IsSerialNumber(caCert, "12345")))

				require.EqualValues([]string{"CH"}, caCert.Issuer.Country)
				require.EqualValues([]string{"Zurich"}, caCert.Issuer.Locality)

				require.True(mustutils.Must(IsCertificateMatchingPrivateKey(caCert, privateKey)))
			},
		)
	}
}
