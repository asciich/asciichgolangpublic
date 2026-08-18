package genericx509utils

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ReadCertFromString(input string) (cert *x509.Certificate, err error) {
	if input == "" {
		return nil, tracederrors.TracedErrorEmptyString("input")
	}

	return ReadCertFromBytes([]byte(input))
}

func ReadCertFromBytes(input []byte) (cert *x509.Certificate, err error) {
	if input == nil {
		return nil, tracederrors.TracedErrorNil("input")
	}

	block, _ := pem.Decode(input)
	if block == nil {
		return nil, tracederrors.TracedError("Failed to decode PEM block from input")
	}

	if block.Bytes == nil {
		return nil, tracederrors.TracedError("PEM decode returned block.Bytes as nil")
	}

	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to parse x509 certificate from DER bytes: %w", err)
	}

	return cert, nil
}

func ReadCertsFromString(input string) (certs []*x509.Certificate, err error) {
	if input == "" {
		return nil, tracederrors.TracedErrorEmptyString("input")
	}

	return ReadCertsFromBytes([]byte(input))
}

func ReadCertsFromBytes(input []byte) (certs []*x509.Certificate, err error) {
	if input == nil {
		return nil, tracederrors.TracedErrorNil("input")
	}

	rest := input
	certs = []*x509.Certificate{}

	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			if len(certs) == 0 {
				return nil, tracederrors.TracedError("Failed to decode any PEM block from input")
			}
			break
		}

		if block.Bytes == nil {
			return nil, tracederrors.TracedError("PEM decode returned block.Bytes as nil")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to parse x509 certificate from DER bytes: %w", err)
		}

		certs = append(certs, cert)

		rest = bytes.TrimSpace(rest)
		if len(rest) == 0 {
			break
		}
	}

	return certs, nil
}
