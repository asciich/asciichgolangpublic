package x509utils

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"slices"
	"time"

	"github.com/asciich/asciichgolangpublic/datatypes/bigintutils"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/cryptoutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var ErrNoValidCertificateChain = errors.New("no valid certificate chain")

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

	publicKey, err := cryptoutils.GetPublicKeyFromPrivateKey(privateKey)
	if err != nil {
		return false, err
	}

	return certPublicKeyWithEqual.Equal(publicKey), nil
}

func LoadCertificatesFromPEMString(pemEncoded string) ([]*x509.Certificate, error) {
	if pemEncoded == "" {
		return nil, tracederrors.TracedErrorEmptyString("pemEncoded")
	}

	rest := []byte(pemEncoded)

	ret := []*x509.Certificate{}
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			return nil, tracederrors.TracedError("Failed to parse certificate PEM")
		} else {
			if block.Bytes == nil {
				return nil, tracederrors.TracedError("Decode returned block.Bytes as nil")
			} else {
				cert, err := LoadCertificateFromDerBytes(block.Bytes)
				if err != nil {
					return nil, err
				}

				ret = append(ret, cert)
			}
		}

		bytes.TrimSpace(rest)
		if len(rest) <= 0 {
			break
		}
	}

	return ret, nil
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

func CheckExpired(ctx context.Context, cert *x509.Certificate) error {
	if cert == nil {
		return tracederrors.TracedErrorNil("cert")
	}

	name := FormatForLogging(cert)

	if cert.NotAfter.Before(time.Now()) {
		return tracederrors.TracedErrorf("%s is not valid anymore. Not after is %s", name, cert.NotAfter.String())
	}

	logging.LogInfoByCtxf(ctx, "%s is still valid until %s.", name, cert.NotAfter.String())
	return nil
}

func CheckCertificateChainString(ctx context.Context, chain string) error {
	if chain == "" {
		return tracederrors.TracedErrorEmptyString(chain)
	}

	logging.LogInfoByCtx(ctx, "Check x509 certificate chain in string started.")

	certs, err := LoadCertificatesFromPEMString(chain)
	if err != nil {
		return err
	}

	if len(certs) != 3 {
		return tracederrors.TracedErrorf("Expected a root, intermediate and endendity certificate but got '%d' certs.", len(certs))
	}

	endEndityCert := certs[0]
	intermediateCert := certs[1]
	rootCaCert := certs[2]

	err = CheckExpired(ctx, endEndityCert)
	if err != nil {
		return err
	}

	err = CheckExpired(ctx, intermediateCert)
	if err != nil {
		return err
	}

	err = CheckExpired(ctx, rootCaCert)
	if err != nil {
		return err
	}

	chains, err := ValidateCertificateChain(ctx, certs[0], []*x509.Certificate{rootCaCert}, []*x509.Certificate{intermediateCert})
	if err != nil {
		return err
	}

	if len(chains) != 1 {
		return tracederrors.TracedErrorf("No valid certificate chain found in string to validate.")
	}

	logging.LogInfoByCtxf(ctx, "Found valid certificate chain in string to validate.")

	logging.LogInfoByCtx(ctx, "Check x509 certificate chain in string finished.")

	return nil
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

func GetCommonName(cert *x509.Certificate) (commonName string, err error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	return cert.Subject.CommonName, nil
}

func GetSans(cert *x509.Certificate) (sans []string, err error) {
	if cert == nil {
		return nil, tracederrors.TracedErrorNil("cert")
	}

	sans = cert.DNSNames
	if sans == nil {
		sans = []string{}
	}

	return sans, nil
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

func IsCommonName(cert *x509.Certificate, expectedCommonName string) (isMatchingExpectedCommonName bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	commonName, err := GetCommonName(cert)
	if err != nil {
		return false, err
	}

	return commonName == expectedCommonName, nil
}

func IsAdditionalSANs(cert *x509.Certificate, expectedSANs []string) (isMatchingexpectedSANs bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	sans, err := GetSans(cert)
	if err != nil {
		return false, err
	}

	return slices.Equal(sans, expectedSANs), nil
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

func CreateSelfSignedCertificate(ctx context.Context, options *X509CreateCertificateOptions) (selfSignedCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	return GetNativeX509CertificateHandler().CreateSelfSignedCertificate(ctx, options)
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

// Returns the duration = notAfter - notBefore.
func GetValidityDuration(cert *x509.Certificate) (validityDuration *time.Duration, err error) {
	if cert == nil {
		return nil, tracederrors.TracedErrorNil("cert")
	}

	diff := cert.NotAfter.Sub(cert.NotBefore)

	return &diff, nil
}

func IsCertSignedBy(ctx context.Context, cert *x509.Certificate, issuerCert *x509.Certificate) (isSigned bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	if issuerCert == nil {
		return false, tracederrors.TracedErrorNil("issuerCert")
	}

	roots := x509.NewCertPool()
	roots.AddCert(issuerCert)

	_, err = cert.Verify(
		x509.VerifyOptions{
			Roots: roots,
		},
	)
	if err != nil {
		if errors.As(err, &x509.UnknownAuthorityError{}) {
			logging.LogInfoByCtxf(ctx, "Certificate '%s' is not signed by '%s'.", cert.Subject, issuerCert.Subject)
			return false, nil
		}

		return false, tracederrors.TracedErrorf("Cert '%s' signed by '%s' failed: %w", cert.Subject, issuerCert.Subject, err)
	}

	logging.LogInfoByCtxf(ctx, "Certificate '%s' is signed by '%s'.", cert.Subject, issuerCert.Subject)
	return true, nil
}

func GetX509CertificateDeepCopy(in *x509.Certificate) (out *x509.Certificate) {
	if in == nil {
		return nil
	}

	out = new(x509.Certificate)
	*out = *in

	return out
}

func FormatForLogging(cert *x509.Certificate) string {
	if cert == nil {
		return "cert is nil"
	}

	out, err := GetSubjectAndSerialString(cert)
	if err != nil {
		return "ERROR: " + err.Error()
	}

	var certType = "<Unknown cert type>"
	isRootCa, err := IsCertificateRootCa(cert)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	if isRootCa {
		certType = "RootCa"
	}

	isIntermediateCertificate, err := IsIntermediateCertificate(cert)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	if isIntermediateCertificate {
		certType = "Intermediate"
	}

	isEndEndityCert, err := IsEndEndityCertificate(cert)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	if isEndEndityCert {
		certType = "end endity"
	}

	out = "x509 " + certType + " certificate " + out

	return out
}

func GetSubjectAndSerialString(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("in")
	}

	serial, err := GetSerialNumberAsHexColonSeparated(cert)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s serial: %s",
		cert.Subject,
		serial,
	), nil
}

func GetSerialNumberAsHexColonSeparated(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("in")
	}

	return bigintutils.ToHexStringColonSeparated(cert.SerialNumber)
}

func GenerateCertificateSerialNumber(ctx context.Context) (serialNumber *big.Int, err error) {
	logging.LogInfoByCtx(ctx, "Generate certificate serial number started.")

	minNumber := big.NewInt(256 * 256 * 256)
	maxNumber := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)

	serialNumber, err = bigintutils.GetRandomBigIntByInts(minNumber, maxNumber)
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	hexRepresentation, err := bigintutils.ToHexStringColonSeparated(serialNumber)
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w", err)
	}

	logging.LogInfoByCtxf(ctx, "Generate certificate serial number finished. Generated serial number is '%s'.", hexRepresentation)

	return serialNumber, nil
}

func GenerateCertificateSerialNumberAsString(ctx context.Context) (string, error) {
	serial, err := GenerateCertificateSerialNumber(ctx)
	if err != nil {
		return "", err
	}

	ret, err := bigintutils.ToDecimalString(serial)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func ValidateCertificateChain(ctx context.Context, certToValidate *x509.Certificate, trustedList []*x509.Certificate, intermediatesList []*x509.Certificate) (chains [][]*x509.Certificate, err error) {
	if certToValidate == nil {
		return nil, tracederrors.TracedErrorNil("certToValidate")
	}

	if trustedList == nil {
		return nil, tracederrors.TracedErrorNil("trustedList")
	}

	if len(trustedList) == 0 {
		return nil, tracederrors.TracedError("trustedList has no entries.")
	}

	if intermediatesList == nil {
		return nil, tracederrors.TracedErrorNil("intermediatesList")
	}

	rootCAPool := x509.NewCertPool()
	for _, rootCert := range trustedList {
		if rootCert == nil {
			return nil, tracederrors.TracedErrorf("Nil pointer found in trustedList '%v'", trustedList)
		}
		rootCAPool.AddCert(rootCert)
	}

	intermediateCAPool := x509.NewCertPool()
	for _, intCert := range intermediatesList {
		if intCert == nil {
			return nil, tracederrors.TracedErrorf("Nil pointer found in intermediatesList '%v'", trustedList)
		}
		intermediateCAPool.AddCert(intCert)
	}

	verifyOptions := x509.VerifyOptions{
		Roots:         rootCAPool,
		Intermediates: intermediateCAPool,
		CurrentTime:   time.Now(),
	}

	chains, err = certToValidate.Verify(verifyOptions)
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w: %w", ErrNoValidCertificateChain, err)
	}

	logging.LogInfoByCtxf(ctx, "Found %d certificate chains to '%s'.", len(chains), FormatForLogging(certToValidate))

	return chains, nil
}
