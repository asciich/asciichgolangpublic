package x509utils

import (
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func IsCertificateRootCa(cert *x509.Certificate) (isRootCa bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	return cert.Subject.String() == cert.Issuer.String(), nil
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
