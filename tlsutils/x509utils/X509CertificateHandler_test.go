package x509utils

import (
	"testing"
	"time"

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

				require.True(t, mustutils.Must(IsCertificateRootCa(caCert)))
				require.True(t, mustutils.Must(IsSelfSignedCertificate(caCert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(caCert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(caCert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))
				require.True(t, mustutils.Must(IsSerialNumber(caCert, "12345")))

				require.EqualValues(t, []string{"CH"}, caCert.Issuer.Country)
				require.EqualValues(t, []string{"Zurich"}, caCert.Issuer.Locality)
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(caCert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(caCert, privateKey)))
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

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(caCert, caPrivateKey)))
				require.True(t, mustutils.Must(IsCertificateRootCa(caCert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(caCert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(caCert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(caCert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCert, intermediateKey)))
				require.False(t, mustutils.Must(IsCertificateRootCa(intermediateCert)))
				require.True(t, mustutils.Must(IsIntermediateCertificate(intermediateCert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(intermediateCert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(intermediateCert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(intermediateCert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(intermediateCert, "myOrg intermediate")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(intermediateCert)))
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
						CommonName:   "mytestcn.example.net",
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg endEndity",

						Verbose: true,
					},
					intermediateCert,
					intermediateKey,
					true,
				))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(caCert, caPrivateKey)))
				require.True(t, mustutils.Must(IsCertificateRootCa(caCert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(caCert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(caCert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(caCert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCert, intermediateKey)))
				require.False(t, mustutils.Must(IsCertificateRootCa(intermediateCert)))
				require.True(t, mustutils.Must(IsIntermediateCertificate(intermediateCert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(intermediateCert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(intermediateCert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(intermediateCert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(intermediateCert, "myOrg intermediate")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(intermediateCert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(endEndityCertificate, endEndityKey)))
				require.False(t, mustutils.Must(IsCertificateRootCa(endEndityCertificate)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(endEndityCertificate)))
				require.True(t, mustutils.Must(IsEndEndityCertificate(endEndityCertificate)))
				require.True(t, mustutils.Must(IsSubjectCountryName(endEndityCertificate, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(endEndityCertificate, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(endEndityCertificate, "myOrg endEndity")))
				require.True(t, mustutils.Must(IsCommonName(endEndityCertificate, "mytestcn.example.net")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(endEndityCertificate)))
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

				require.False(t, mustutils.Must(IsCertificateRootCa(caCert)))
				require.True(t, mustutils.Must(IsSelfSignedCertificate(caCert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(caCert)))
				require.True(t, mustutils.Must(IsEndEndityCertificate(caCert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(caCert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(caCert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(caCert, "myOrg root")))
				require.True(t, mustutils.Must(IsSerialNumber(caCert, "12345")))

				require.EqualValues(t, []string{"CH"}, caCert.Issuer.Country)
				require.EqualValues(t, []string{"Zurich"}, caCert.Issuer.Locality)
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(caCert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(caCert, privateKey)))
			},
		)
	}
}
