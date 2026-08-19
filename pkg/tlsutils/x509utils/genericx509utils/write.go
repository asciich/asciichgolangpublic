package genericx509utils

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func WriteCertAsString(cert *x509.Certificate) (string, error) {
	if cert == nil {
		return "", tracederrors.TracedErrorNil("cert")
	}

	outputBytes, err := WriteCertAsBytes(cert)
	if err != nil {
		return "", err
	}

	return string(outputBytes), nil
}

func WriteCertAsBytes(cert *x509.Certificate) ([]byte, error) {
	if cert == nil {
		return nil, tracederrors.TracedErrorNil("cert")
	}

	if cert.Raw == nil {
		return nil, tracederrors.TracedError("cert.Raw is nil, cannot encode certificate")
	}

	var buf bytes.Buffer
	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}

	err := pem.Encode(io.Writer(&buf), pemBlock)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to PEM encode certificate: %w", err)
	}

	return buf.Bytes(), nil
}

func WriteCertsAsString(certs []*x509.Certificate) (string, error) {
	if certs == nil {
		return "", tracederrors.TracedErrorNil("certs")
	}

	if len(certs) == 0 {
		return "", tracederrors.TracedError("certs has no entries")
	}

	outputBytes, err := WriteCertsAsBytes(certs)
	if err != nil {
		return "", err
	}

	return string(outputBytes), nil
}

func WriteCertsAsBytes(certs []*x509.Certificate) ([]byte, error) {
	if certs == nil {
		return nil, tracederrors.TracedErrorNil("certs")
	}

	if len(certs) == 0 {
		return nil, tracederrors.TracedError("certs has no entries")
	}

	var buf bytes.Buffer

	for i, cert := range certs {
		if cert == nil {
			return nil, tracederrors.TracedErrorf("cert at index '%d' is nil", i)
		}

		if cert.Raw == nil {
			return nil, tracederrors.TracedErrorf("cert.Raw at index '%d' is nil, cannot encode certificate", i)
		}

		pemBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		}

		err := pem.Encode(io.Writer(&buf), pemBlock)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to PEM encode certificate at index '%d': %w", i, err)
		}
	}

	return buf.Bytes(), nil
}

func WriteCertificateAsPEMString(cert *x509.Certificate) (string, error) {
	return WriteCertAsString(cert)
}

func WriteCertificateAsPEMBytes(cert *x509.Certificate) ([]byte, error) {
	return WriteCertAsBytes(cert)
}

func WriteCertificatesAsPEMString(certs []*x509.Certificate) (string, error) {
	return WriteCertsAsString(certs)
}

func WriteCertificatesAsPEMBytes(certs []*x509.Certificate) ([]byte, error) {
	return WriteCertsAsBytes(certs)
}
