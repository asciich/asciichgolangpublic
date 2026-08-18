package genericx509utils

import (
	"context"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// --- Root CA ---

func IsRootCaCert(ctx context.Context, cert *x509.Certificate) (bool, error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	if cert.Version == 1 {
		return false, nil
	}

	if cert.Subject.String() != cert.Issuer.String() {
		return false, nil
	}

	return cert.IsCA, nil
}

func IsStringRootCaCert(ctx context.Context, input string) (bool, error) {
	if input == "" {
		return false, tracederrors.TracedErrorEmptyString("input")
	}

	cert, err := ReadCertFromString(input)
	if err != nil {
		return false, err
	}

	return IsRootCaCert(ctx, cert)
}

func IsBytesRootCaCert(ctx context.Context, input []byte) (bool, error) {
	if input == nil {
		return false, tracederrors.TracedErrorNil("input")
	}

	cert, err := ReadCertFromBytes(input)
	if err != nil {
		return false, err
	}

	return IsRootCaCert(ctx, cert)
}

// --- Intermediate ---

func IsIntermediateCert(ctx context.Context, cert *x509.Certificate) (bool, error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	if cert.Version == 1 {
		return false, nil
	}

	if cert.Subject.String() == cert.Issuer.String() {
		return false, nil
	}

	return cert.IsCA, nil
}

func IsStringIntermediateCert(ctx context.Context, input string) (bool, error) {
	if input == "" {
		return false, tracederrors.TracedErrorEmptyString("input")
	}

	cert, err := ReadCertFromString(input)
	if err != nil {
		return false, err
	}

	return IsIntermediateCert(ctx, cert)
}

func IsBytesIntermediateCert(ctx context.Context, input []byte) (bool, error) {
	if input == nil {
		return false, tracederrors.TracedErrorNil("input")
	}

	cert, err := ReadCertFromBytes(input)
	if err != nil {
		return false, err
	}

	return IsIntermediateCert(ctx, cert)
}

// --- End Entity (leaf certificate used by servers/services) ---

func IsEndEntityCert(ctx context.Context, cert *x509.Certificate) (bool, error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
	}

	return !cert.IsCA, nil
}

func IsStringEndEntityCert(ctx context.Context, input string) (bool, error) {
	if input == "" {
		return false, tracederrors.TracedErrorEmptyString("input")
	}

	cert, err := ReadCertFromString(input)
	if err != nil {
		return false, err
	}

	return IsEndEntityCert(ctx, cert)
}

func IsBytesEndEntityCert(ctx context.Context, input []byte) (bool, error) {
	if input == nil {
		return false, tracederrors.TracedErrorNil("input")
	}

	cert, err := ReadCertFromBytes(input)
	if err != nil {
		return false, err
	}

	return IsEndEntityCert(ctx, cert)
}
