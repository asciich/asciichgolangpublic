package x509utils

import (
	"crypto"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func IsIntermediateCertificate(cert *x509.Certificate) (isIntermediateCertificate bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	if cert.Version == 1 {
		// version 1 is not supported for root CA.
		return false, nil
	}

	if cert.Subject.String() == cert.Issuer.String() {
		return false, err
	}

	return cert.IsCA, nil
}

// An End-Endity certificate is a cert used by the systems/ services.
// So it's neither an intermedate nor a rootCA certificate.
func IsEndEndityCertificate(cert *x509.Certificate) (isIntermediateCertificate bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	return !cert.IsCA, nil
}

func IsCertificateRootCa(cert *x509.Certificate) (isRootCa bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	if cert.Version == 1 {
		// version 1 is not supported for root CA.
		return false, nil
	}

	if cert.Subject.String() != cert.Issuer.String() {
		return false, err
	}

	return cert.IsCA, nil
}

func GetPublicKeyFromPrivateKey(privateKey crypto.PrivateKey) (publicKey crypto.PublicKey, err error) {
	if privateKey == "" {
		return nil, tracederrors.TracedErrorNil("privateKey")
	}

	withPublic, ok := privateKey.(interface{ Public() crypto.PublicKey })
	if !ok {
		return nil, tracederrors.TracedErrorf("Unable to get publicKey out of privateKey. privateKey does not implement Public as expected (and done by all stdlib private keys).")
	}

	return withPublic.Public(), nil
}

func IsCertificateVersion1(cert *x509.Certificate) (isV1 bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	return cert.Version == 1, nil
}

func IsCertificateMatchingPrivateKey(cert *x509.Certificate, privateKey crypto.PrivateKey) (isMatching bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	if privateKey == nil {
		return false, tracederrors.TracedErrorNil("privateKey")
	}

	certPublicKey := cert.PublicKey

	certPublicKeyWithEqual, ok := certPublicKey.(interface{ Equal(x crypto.PublicKey) bool })
	if !ok {
		return false, tracederrors.TracedError("certPublicKey does not implement Equal as expected (and done by all stdlib private keys).")
	}

	publicKey, err := GetPublicKeyFromPrivateKey(privateKey)
	if err != nil {
		return false, err
	}

	return certPublicKeyWithEqual.Equal(publicKey), nil
}

func LoadCertificateFromDerBytes(derEncodecCertificate []byte) (cert *x509.Certificate, err error) {
	if derEncodecCertificate == nil {
		return nil, tracederrors.TracedErrorNil("derEncodecCertificate")
	}

	cert, err = x509.ParseCertificate(derEncodecCertificate)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to decode derEncodecCertificate: %w", err)
	}

	return cert, nil
}

func GetSubjectCountryName(cert *x509.Certificate) (countryName string, err error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	country := cert.Subject.Country

	nCountries := len(country)
	if nCountries == 0 {
		return "", nil
	}

	if nCountries == 1 {
		return country[0], nil
	}

	return "", tracederrors.TracedErrorf(
		"Not implemented for nCountries != 1. Got '%d' countries: '%v'",
		nCountries,
		country,
	)
}


func GetSubjectLocalityName(cert *x509.Certificate) (locality string, err error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	country := cert.Subject.Locality

	nLocalities := len(country)
	if nLocalities == 0 {
		return "", nil
	}

	if nLocalities == 1 {
		return country[0], nil
	}

	return "", tracederrors.TracedErrorf(
		"Not implemented for nLocalities != 1. Got '%d' localities: '%v'",
		nLocalities,
		country,
	)
}

func GetSubjectOrganizationName(cert *x509.Certificate) (organizationName string, err error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	organization := cert.Subject.Organization

	nOrganizations := len(organization)
	if nOrganizations == 0 {
		return "", nil
	}

	if nOrganizations == 1 {
		return organization[0], nil
	}

	return "", tracederrors.TracedErrorf(
		"Not implemented for nLocalities != 1. Got '%d' localities: '%v'",
		nOrganizations,
		organization,
	)
}

func IsSubjectCountryName(cert *x509.Certificate, expectedCountryName string) (isMatchingExpectedCountryName bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	countryName, err := GetSubjectCountryName(cert)
	if err != nil {
		return false, err
	}

	return countryName == expectedCountryName, nil
}

func IsSubjectLocalityName(cert *x509.Certificate, expectedLocalityName string) (isMatchingExpectedLocalityName bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	localityName, err := GetSubjectLocalityName(cert)
	if err != nil {
		return false, err
	}

	return localityName == expectedLocalityName, nil
}

func IsSubjectOrganizationName(cert *x509.Certificate, expectedOrganizationName string) (isMatchingExpectedOrganizationName bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	organizationName, err := GetSubjectOrganizationName(cert)
	if err != nil {
		return false, err
	}

	return organizationName == expectedOrganizationName, nil
}
