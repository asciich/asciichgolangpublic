package x509utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
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

				caCertKeyPair, err := handler.CreateRootCaCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",
						SerialNumber: "12345",
					},
				)
				require.NoError(t, err)

				require.True(t, mustutils.Must(IsCertificateRootCa(caCertKeyPair.Cert)))
				require.True(t, mustutils.Must(IsSelfSignedCertificate(caCertKeyPair.Cert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(caCertKeyPair.Cert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(caCertKeyPair.Cert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(caCertKeyPair.Cert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(caCertKeyPair.Cert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(caCertKeyPair.Cert, "myOrg root")))
				require.True(t, mustutils.Must(IsSerialNumber(caCertKeyPair.Cert, "12345")))

				require.EqualValues(t, []string{"CH"}, caCertKeyPair.Cert.Issuer.Country)
				require.EqualValues(t, []string{"Zurich"}, caCertKeyPair.Cert.Issuer.Locality)
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(caCertKeyPair.Cert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(caCertKeyPair.Cert, caCertKeyPair.Key)))
			},
		)
	}
}

func TestX509Handler_CreateIntermediateCertificate(t *testing.T) {
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

				caCertKeyPair, err := handler.CreateRootCaCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",
					},
				)
				require.NoError(t, err)

				intermediateCertAndKey, err := handler.CreateSignedIntermediateCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg intermediate",
					},
					caCertKeyPair,
				)
				require.NoError(t, err)

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(caCertKeyPair.Cert, caCertKeyPair.Key)))
				require.True(t, mustutils.Must(IsCertificateRootCa(caCertKeyPair.Cert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(caCertKeyPair.Cert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(caCertKeyPair.Cert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(caCertKeyPair.Cert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(caCertKeyPair.Cert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(caCertKeyPair.Cert, "myOrg root")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(caCertKeyPair.Cert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCertAndKey.Cert, intermediateCertAndKey.Key)))
				require.False(t, mustutils.Must(IsCertificateRootCa(intermediateCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsIntermediateCertificate(intermediateCertAndKey.Cert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(intermediateCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(intermediateCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(intermediateCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(intermediateCertAndKey.Cert, "myOrg intermediate")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(intermediateCertAndKey.Cert)))
			},
		)
	}
}

func TestX509Handler_CreateEndEndityCertificate(t *testing.T) {
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

				rootCeCertAndKey, err := handler.CreateRootCaCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",
					},
				)
				require.NoError(t, err)

				intermediateCertAndKey, err := handler.CreateSignedIntermediateCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg intermediate",
					},
					rootCeCertAndKey,
				)
				require.NoError(t, err)

				endEndityCertAndKey, err := handler.CreateSignedEndEndityCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CommonName:     "mytestcn.example.net",
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "myOrg endEndity",
						AdditionalSans: []string{"mytestsan1.example.net", "mytestsan2.example.net"},
					},
					intermediateCertAndKey,
				)
				require.NoError(t, err)

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(rootCeCertAndKey.Cert, rootCeCertAndKey.Key)))
				require.True(t, mustutils.Must(IsCertificateRootCa(rootCeCertAndKey.Cert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(rootCeCertAndKey.Cert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(rootCeCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(rootCeCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(rootCeCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(rootCeCertAndKey.Cert, "myOrg root")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(rootCeCertAndKey.Cert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(intermediateCertAndKey.Cert, intermediateCertAndKey.Key)))
				require.False(t, mustutils.Must(IsCertificateRootCa(intermediateCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsIntermediateCertificate(intermediateCertAndKey.Cert)))
				require.False(t, mustutils.Must(IsEndEndityCertificate(intermediateCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(intermediateCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(intermediateCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(intermediateCertAndKey.Cert, "myOrg intermediate")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(intermediateCertAndKey.Cert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(endEndityCertAndKey.Cert, endEndityCertAndKey.Key)))
				require.False(t, mustutils.Must(IsCertificateRootCa(endEndityCertAndKey.Cert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(endEndityCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsEndEndityCertificate(endEndityCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(endEndityCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(endEndityCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(endEndityCertAndKey.Cert, "myOrg endEndity")))
				require.True(t, mustutils.Must(IsCommonName(endEndityCertAndKey.Cert, "mytestcn.example.net")))
				require.True(t, mustutils.Must(IsAdditionalSANs(endEndityCertAndKey.Cert, []string{"mytestsan1.example.net", "mytestsan2.example.net"})))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(endEndityCertAndKey.Cert)))
			},
		)
	}
}

func TestX509Handler_CreateSelfSignedCertificate(t *testing.T) {
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

				rootCaCertAndKey, err := handler.CreateSelfSignedCertificate(
					ctx,
					&X509CreateCertificateOptions{
						CountryName:  "CH",
						Locality:     "Zurich",
						Organization: "myOrg root",
						SerialNumber: "12345",
					},
				)
				require.NoError(t, err)

				require.False(t, mustutils.Must(IsCertificateRootCa(rootCaCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsSelfSignedCertificate(rootCaCertAndKey.Cert)))
				require.False(t, mustutils.Must(IsIntermediateCertificate(rootCaCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsEndEndityCertificate(rootCaCertAndKey.Cert)))
				require.True(t, mustutils.Must(IsSubjectCountryName(rootCaCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(IsSubjectLocalityName(rootCaCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(IsSubjectOrganizationName(rootCaCertAndKey.Cert, "myOrg root")))
				require.True(t, mustutils.Must(IsSerialNumber(rootCaCertAndKey.Cert, "12345")))

				require.EqualValues(t, []string{"CH"}, rootCaCertAndKey.Cert.Issuer.Country)
				require.EqualValues(t, []string{"Zurich"}, rootCaCertAndKey.Cert.Issuer.Locality)
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(GetValidityDuration(rootCaCertAndKey.Cert)))

				require.True(t, mustutils.Must(IsCertificateMatchingPrivateKey(rootCaCertAndKey.Cert, rootCaCertAndKey.Key)))
			},
		)
	}
}
