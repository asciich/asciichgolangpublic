package genericx509utils

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func TlsCertToX509Cert(tlsCert *tls.Certificate) (*x509.Certificate, error) {
	if tlsCert == nil {
		return nil, tracederrors.TracedErrorNil("tlsCert")
	}

	if len(tlsCert.Certificate) == 0 {
		return nil, tracederrors.TracedError("tls.Certificate has no certificate data")
	}

	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse x509 certificate from tls.Certificate: %w", err)
	}

	return cert, nil
}
