package x509utils

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/logging"
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

func IsSelfSignedCertificate(cert *x509.Certificate) (isSelfSigend bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	return cert.Subject.String() == cert.Issuer.String(), nil
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

func LoadCertificateFromPEMString(pemEncoded string) (cert *x509.Certificate, err error) {
	if pemEncoded == "" {
		return nil, tracederrors.TracedErrorEmptyString("pemEncoded")
	}

	block, _ := pem.Decode([]byte(pemEncoded))
	var blockBytes []byte = nil
	if block == nil {
		return nil, tracederrors.TracedError("Failed to parse certificate PEM")
	} else {
		if block.Bytes == nil {
			return nil, tracederrors.TracedError("Decode returned block.Bytes as nil")
		} else {
			blockBytes = block.Bytes
		}
	}

	return LoadCertificateFromDerBytes(blockBytes)
}

func LoadPrivateKeyFromPEMString(pemEncoded string) (privateKey crypto.PrivateKey, err error) {
	if pemEncoded == "" {
		return nil, tracederrors.TracedErrorEmptyString("pemEncoded")
	}

	block, _ := pem.Decode([]byte(pemEncoded))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, tracederrors.TracedError("invalid private key.")
	}

	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		privateKey = key
	} else if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		privateKey = key
	} else {
		return nil, tracederrors.TracedErrorf("Unable to parse as key: %w", err)
	}

	return privateKey, nil
}

func EncodePrivateKeyAsPEMString(privateKey crypto.PrivateKey) (pemEncoded string, err error) {
	if privateKey == nil {
		return "", tracederrors.TracedErrorNil("privateKey")
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to marshal private key: '%w'", err)
	}

	var buf bytes.Buffer
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	err = pem.Encode(io.Writer(&buf), pemBlock)
	if err != nil {
		return "", err
	}

	pemEncoded = buf.String()
	pemEncoded = stringsutils.EnsureEndsWithExactlyOneLineBreak(pemEncoded)

	const minLen = 50
	if len(pemEncoded) < minLen {
		return "", tracederrors.TracedErrorf("pemBytes has less than '%v' bytes which is not enough for a pem certificate", minLen)
	}

	return pemEncoded, nil
}

func EncodeCertificateAsPEMString(cert *x509.Certificate) (pemEncoded string, err error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("derEncodecCertificate")
	}

	derBytes, err := EncodeCertificateAsDerBytes(cert)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	var pemBlock = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}
	err = pem.Encode(io.Writer(&buf), pemBlock)
	if err != nil {
		return "", err
	}

	pemEncoded = buf.String()
	pemEncoded = stringsutils.EnsureEndsWithExactlyOneLineBreak(pemEncoded)

	const minLen = 50
	if len(pemEncoded) < minLen {
		return "", tracederrors.TracedErrorf("pemBytes has less than '%v' bytes which is not enough for a pem certificate", minLen)
	}

	return pemEncoded, nil
}

func EncodeCertificateAsDerBytes(cert *x509.Certificate) (derEncodecCertificate []byte, err error) {
	if cert == nil {
		return
	}

	return cert.Raw, nil
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

func GetSerialNumberAsString(cert *x509.Certificate) (serialNumber string, err error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	serial := cert.SerialNumber
	if serial == nil {
		return "", tracederrors.TracedError("unable to get serial number from x509 certificate. SerialNumber is nil")
	}

	serialNumber = serial.String()
	if serialNumber == "" {
		return "", tracederrors.TracedError("Serial number is empty string after evaluation")
	}

	return serialNumber, nil
}

func IsSerialNumber(cert *x509.Certificate, expectedSerialNumber string) (isSerialNumber bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	serialNumber, err := GetSerialNumberAsString(cert)
	if err != nil {
		return false, err
	}

	return serialNumber == expectedSerialNumber, nil
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

func IsPrivateKeyEqual(key1 crypto.PrivateKey, key2 crypto.PrivateKey) (isEqual bool, err error) {
	if key1 == nil {
		return false, tracederrors.TracedErrorNil("key1")
	}

	if key2 == nil {
		return false, tracederrors.TracedErrorNil("key2")
	}

	withEqual, ok := key1.(interface {
		Equal(other crypto.PrivateKey) bool
	})
	if !ok {
		return false, tracederrors.TracedErrorf("key 1 does not implement Equal function to other private keys.")
	}

	return withEqual.Equal(key2), nil
}

func CreateSelfSignedCertificate(options *X509CreateCertificateOptions) (selfSignedCert *x509.Certificate, selfSignedCertKey crypto.PrivateKey, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	return GetNativeX509CertificateHandler().CreateSelfSignedCertificate(options)
}

func MustTlsCertToX509Cert(tlsCert *tls.Certificate) (cert *x509.Certificate) {
	cert, err := TlsCertToX509Cert(tlsCert)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cert
}

func TlsCertToX509Cert(tlsCert *tls.Certificate) (cert *x509.Certificate, err error) {
	if len(tlsCert.Certificate) == 0 {
		return nil, fmt.Errorf("tls.Certificate has no certificate data")
	}

	x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse x509 certificate: %w", err)
	}

	return x509Cert, nil
}
