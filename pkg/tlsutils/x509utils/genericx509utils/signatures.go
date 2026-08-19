package genericx509utils

import (
	"context"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func IsSignedBy(ctx context.Context, certToCheck *x509.Certificate, signingCert *x509.Certificate) (bool, error) {
	if certToCheck == nil {
		return false, tracederrors.TracedErrorNil("certToCheck")
	}

	if signingCert == nil {
		return false, tracederrors.TracedErrorNil("signingCert")
	}

	certToCheckStr, err := GetSubjectAndSerialString(certToCheck)
	if err != nil {
		return false, err
	}

	signingCertStr, err := GetSubjectAndSerialString(signingCert)
	if err != nil {
		return false, err
	}

	err = certToCheck.CheckSignatureFrom(signingCert)
	if err != nil {
		logging.LogInfoByCtxf(ctx, "Certificate '%s' is NOT signed by '%s': %v", certToCheckStr, signingCertStr, err)
		return false, nil
	}

	logging.LogInfoByCtxf(ctx, "Certificate '%s' IS signed by '%s'.", certToCheckStr, signingCertStr)
	return true, nil
}

// Check if at least two certs are in the chainToCheck.
// The order is validated as: The first entry is signed by the second which is signed by the next...
func IsCertChain(ctx context.Context, chainToCheck []*x509.Certificate) (bool, error) {
	if chainToCheck == nil {
		return false, tracederrors.TracedErrorNil("chainToCheck")
	}

	nCerts := len(chainToCheck)
	if nCerts < 2 {
		logging.LogInfoByCtxf(ctx, "Chain has '%d' certificates, at least 2 are required for a valid chain.", nCerts)
		return false, nil
	}

	for i := 0; i < nCerts-1; i++ {
		current := chainToCheck[i]
		next := chainToCheck[i+1]

		if current == nil {
			return false, tracederrors.TracedErrorf("cert at index '%d' is nil", i)
		}

		if next == nil {
			return false, tracederrors.TracedErrorf("cert at index '%d' is nil", i+1)
		}

		currentStr, err := GetSubjectAndSerialString(current)
		if err != nil {
			return false, err
		}

		nextStr, err := GetSubjectAndSerialString(next)
		if err != nil {
			return false, err
		}

		isSigned, err := IsSignedBy(ctx, current, next)
		if err != nil {
			return false, err
		}

		if !isSigned {
			logging.LogInfoByCtxf(ctx, "Chain is broken at index '%d': certificate '%s' is NOT signed by '%s'.", i, currentStr, nextStr)
			return false, nil
		}

		logging.LogInfoByCtxf(ctx, "Chain link '%d' valid: certificate '%s' is signed by '%s'.", i, currentStr, nextStr)
	}

	logging.LogInfoByCtxf(ctx, "Certificate chain with '%d' certificates is valid.", nCerts)
	return true, nil
}

func IsCertChainFromString(ctx context.Context, chainToCheck string) (bool, error) {
	if chainToCheck == "" {
		return false, tracederrors.TracedErrorEmptyString("chainToCheck")
	}

	certs, err := ReadCertsFromString(chainToCheck)
	if err != nil {
		return false, err
	}

	return IsCertChain(ctx, certs)
}

func IsCertChainFromBytes(ctx context.Context, chainToCheck []byte) (bool, error) {
	if chainToCheck == nil {
		return false, tracederrors.TracedErrorNil("chainToCheck")
	}

	certs, err := ReadCertsFromBytes(chainToCheck)
	if err != nil {
		return false, err
	}

	return IsCertChain(ctx, certs)
}

// Check if at least two certs are in the chainToCheck.
// The order is validated as: The first entry is signed by the second which is signed by the next...
// The first one must be an end entity certificate.
// The last one must be the RootCa.
func IsRootCaToEndEntityChain(ctx context.Context, chainToCheck []*x509.Certificate) (bool, error) {
	if chainToCheck == nil {
		return false, tracederrors.TracedErrorNil("chainToCheck")
	}

	nCerts := len(chainToCheck)
	if nCerts < 2 {
		logging.LogInfoByCtxf(ctx, "Chain has '%d' certificates, at least 2 are required for a valid RootCA-to-EndEntity chain.", nCerts)
		return false, nil
	}

	// Check the first cert is an end entity certificate.
	firstCert := chainToCheck[0]
	if firstCert == nil {
		return false, tracederrors.TracedError("first cert in chain is nil")
	}

	firstCertStr, err := GetSubjectAndSerialString(firstCert)
	if err != nil {
		return false, err
	}

	isEndEntity, err := IsEndEntityCert(ctx, firstCert)
	if err != nil {
		return false, err
	}

	if !isEndEntity {
		logging.LogInfoByCtxf(ctx, "First certificate in chain '%s' is NOT an end entity certificate.", firstCertStr)
		return false, nil
	}

	logging.LogInfoByCtxf(ctx, "First certificate in chain '%s' IS an end entity certificate.", firstCertStr)

	// Check the last cert is a Root CA certificate.
	lastCert := chainToCheck[nCerts-1]
	if lastCert == nil {
		return false, tracederrors.TracedError("last cert in chain is nil")
	}

	lastCertStr, err := GetSubjectAndSerialString(lastCert)
	if err != nil {
		return false, err
	}

	isRootCa, err := IsRootCaCert(ctx, lastCert)
	if err != nil {
		return false, err
	}

	if !isRootCa {
		logging.LogInfoByCtxf(ctx, "Last certificate in chain '%s' is NOT a Root CA certificate.", lastCertStr)
		return false, nil
	}

	logging.LogInfoByCtxf(ctx, "Last certificate in chain '%s' IS a Root CA certificate.", lastCertStr)

	// Check intermediate certificates if any.
	for i := 1; i < nCerts-1; i++ {
		intermediateCert := chainToCheck[i]
		if intermediateCert == nil {
			return false, tracederrors.TracedErrorf("cert at index '%d' is nil", i)
		}

		intermediateCertStr, err := GetSubjectAndSerialString(intermediateCert)
		if err != nil {
			return false, err
		}

		isIntermediate, err := IsIntermediateCert(ctx, intermediateCert)
		if err != nil {
			return false, err
		}

		if !isIntermediate {
			logging.LogInfoByCtxf(ctx, "Certificate at index '%d' ('%s') is NOT an intermediate certificate.", i, intermediateCertStr)
			return false, nil
		}

		logging.LogInfoByCtxf(ctx, "Certificate at index '%d' ('%s') IS an intermediate certificate.", i, intermediateCertStr)
	}

	// Validate the signing chain.
	isChain, err := IsCertChain(ctx, chainToCheck)
	if err != nil {
		return false, err
	}

	if !isChain {
		logging.LogInfoByCtxf(ctx, "Certificates have correct types but the signing chain is invalid.")
		return false, nil
	}

	logging.LogInfoByCtxf(ctx, "Valid RootCA-to-EndEntity chain with '%d' certificates verified.", nCerts)
	return true, nil
}

func IsRootCaToEndEntityChainFromString(ctx context.Context, chainToCheck string) (bool, error) {
	if chainToCheck == "" {
		return false, tracederrors.TracedErrorEmptyString("chainToCheck")
	}

	certs, err := ReadCertsFromString(chainToCheck)
	if err != nil {
		return false, err
	}

	return IsRootCaToEndEntityChain(ctx, certs)
}

func IsRootCaToEndEntityChainFromBytes(ctx context.Context, chainToCheck []byte) (bool, error) {
	if chainToCheck == nil {
		return false, tracederrors.TracedErrorNil("chainToCheck")
	}

	certs, err := ReadCertsFromBytes(chainToCheck)
	if err != nil {
		return false, err
	}

	return IsRootCaToEndEntityChain(ctx, certs)
}

func CheckCertificateChainString(ctx context.Context, chain string) error {
	if chain == "" {
		return tracederrors.TracedErrorEmptyString("chain")
	}

	certs, err := ReadCertsFromString(chain)
	if err != nil {
		return err
	}

	if len(certs) < 3 {
		return tracederrors.TracedErrorf("Expected at least a root, intermediate and end entity certificate but got '%d' certs.", len(certs))
	}

	isValid, err := IsRootCaToEndEntityChain(ctx, certs)
	if err != nil {
		return err
	}

	if !isValid {
		return tracederrors.TracedError("Certificate chain string is not a valid RootCA-to-EndEntity chain.")
	}

	logging.LogInfoByCtx(ctx, "Certificate chain string is valid.")

	return nil
}
