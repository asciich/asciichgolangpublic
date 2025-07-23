package x509utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/mustutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/tlsutils/x509utils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
)

func getX509CertificateHandlerToTest(implementationName string) (handler x509utils.X509CertificateHandler) {
	if implementationName == "NativeX509CertificateHandler" {
		return x509utils.GetNativeX509CertificateHandler()
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
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "myOrg root",
						SerialNumber:   "12345",
						PrivateKeySize: 1024,
					},
				)
				require.NoError(t, err)

				require.True(t, mustutils.Must(x509utils.IsCertificateRootCa(caCertKeyPair.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSelfSignedCertificate(caCertKeyPair.Cert)))
				require.False(t, mustutils.Must(x509utils.IsIntermediateCertificate(caCertKeyPair.Cert)))
				require.False(t, mustutils.Must(x509utils.IsEndEndityCertificate(caCertKeyPair.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSubjectCountryName(caCertKeyPair.Cert, "CH")))
				require.True(t, mustutils.Must(x509utils.IsSubjectLocalityName(caCertKeyPair.Cert, "Zurich")))
				require.True(t, mustutils.Must(x509utils.IsSubjectOrganizationName(caCertKeyPair.Cert, "myOrg root")))
				require.True(t, mustutils.Must(x509utils.IsSerialNumber(caCertKeyPair.Cert, "12345")))

				require.EqualValues(t, []string{"CH"}, caCertKeyPair.Cert.Issuer.Country)
				require.EqualValues(t, []string{"Zurich"}, caCertKeyPair.Cert.Issuer.Locality)
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(x509utils.GetValidityDuration(caCertKeyPair.Cert)))

				require.True(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(caCertKeyPair.Cert, caCertKeyPair.Key)))
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
					caCertKeyPair,
				)
				require.NoError(t, err)

				require.True(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(caCertKeyPair.Cert, caCertKeyPair.Key)))
				require.True(t, mustutils.Must(x509utils.IsCertificateRootCa(caCertKeyPair.Cert)))
				require.False(t, mustutils.Must(x509utils.IsIntermediateCertificate(caCertKeyPair.Cert)))
				require.False(t, mustutils.Must(x509utils.IsEndEndityCertificate(caCertKeyPair.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSubjectCountryName(caCertKeyPair.Cert, "CH")))
				require.True(t, mustutils.Must(x509utils.IsSubjectLocalityName(caCertKeyPair.Cert, "Zurich")))
				require.True(t, mustutils.Must(x509utils.IsSubjectOrganizationName(caCertKeyPair.Cert, "myOrg root")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(x509utils.GetValidityDuration(caCertKeyPair.Cert)))

				require.True(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(intermediateCertAndKey.Cert, intermediateCertAndKey.Key)))
				require.False(t, mustutils.Must(x509utils.IsCertificateRootCa(intermediateCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsIntermediateCertificate(intermediateCertAndKey.Cert)))
				require.False(t, mustutils.Must(x509utils.IsEndEndityCertificate(intermediateCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSubjectCountryName(intermediateCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(x509utils.IsSubjectLocalityName(intermediateCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(x509utils.IsSubjectOrganizationName(intermediateCertAndKey.Cert, "myOrg intermediate")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(x509utils.GetValidityDuration(intermediateCertAndKey.Cert)))
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

				require.True(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(rootCeCertAndKey.Cert, rootCeCertAndKey.Key)))
				require.True(t, mustutils.Must(x509utils.IsCertificateRootCa(rootCeCertAndKey.Cert)))
				require.False(t, mustutils.Must(x509utils.IsIntermediateCertificate(rootCeCertAndKey.Cert)))
				require.False(t, mustutils.Must(x509utils.IsEndEndityCertificate(rootCeCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSubjectCountryName(rootCeCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(x509utils.IsSubjectLocalityName(rootCeCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(x509utils.IsSubjectOrganizationName(rootCeCertAndKey.Cert, "myOrg root")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(x509utils.GetValidityDuration(rootCeCertAndKey.Cert)))

				require.True(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(intermediateCertAndKey.Cert, intermediateCertAndKey.Key)))
				require.False(t, mustutils.Must(x509utils.IsCertificateRootCa(intermediateCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsIntermediateCertificate(intermediateCertAndKey.Cert)))
				require.False(t, mustutils.Must(x509utils.IsEndEndityCertificate(intermediateCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSubjectCountryName(intermediateCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(x509utils.IsSubjectLocalityName(intermediateCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(x509utils.IsSubjectOrganizationName(intermediateCertAndKey.Cert, "myOrg intermediate")))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(x509utils.GetValidityDuration(intermediateCertAndKey.Cert)))

				require.True(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(endEndityCertAndKey.Cert, endEndityCertAndKey.Key)))
				require.False(t, mustutils.Must(x509utils.IsCertificateRootCa(endEndityCertAndKey.Cert)))
				require.False(t, mustutils.Must(x509utils.IsIntermediateCertificate(endEndityCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsEndEndityCertificate(endEndityCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSubjectCountryName(endEndityCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(x509utils.IsSubjectLocalityName(endEndityCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(x509utils.IsSubjectOrganizationName(endEndityCertAndKey.Cert, "myOrg endEndity")))
				require.True(t, mustutils.Must(x509utils.IsCommonName(endEndityCertAndKey.Cert, "mytestcn.example.net")))
				require.True(t, mustutils.Must(x509utils.IsAdditionalSANs(endEndityCertAndKey.Cert, []string{"mytestcn.example.net", "mytestsan1.example.net", "mytestsan2.example.net"})))
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(x509utils.GetValidityDuration(endEndityCertAndKey.Cert)))
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
					&x509utils.X509CreateCertificateOptions{
						CountryName:    "CH",
						Locality:       "Zurich",
						Organization:   "myOrg root",
						SerialNumber:   "12345",
						PrivateKeySize: 1024,
					},
				)
				require.NoError(t, err)

				require.False(t, mustutils.Must(x509utils.IsCertificateRootCa(rootCaCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSelfSignedCertificate(rootCaCertAndKey.Cert)))
				require.False(t, mustutils.Must(x509utils.IsIntermediateCertificate(rootCaCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsEndEndityCertificate(rootCaCertAndKey.Cert)))
				require.True(t, mustutils.Must(x509utils.IsSubjectCountryName(rootCaCertAndKey.Cert, "CH")))
				require.True(t, mustutils.Must(x509utils.IsSubjectLocalityName(rootCaCertAndKey.Cert, "Zurich")))
				require.True(t, mustutils.Must(x509utils.IsSubjectOrganizationName(rootCaCertAndKey.Cert, "myOrg root")))
				require.True(t, mustutils.Must(x509utils.IsSerialNumber(rootCaCertAndKey.Cert, "12345")))

				require.EqualValues(t, []string{"CH"}, rootCaCertAndKey.Cert.Issuer.Country)
				require.EqualValues(t, []string{"Zurich"}, rootCaCertAndKey.Cert.Issuer.Locality)
				require.EqualValues(t, time.Hour*24*45, *mustutils.Must(x509utils.GetValidityDuration(rootCaCertAndKey.Cert)))

				require.True(t, mustutils.Must(x509utils.IsCertificateMatchingPrivateKey(rootCaCertAndKey.Cert, rootCaCertAndKey.Key)))
			},
		)
	}
}
