package x509utils

import (
	"crypto"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func IsCertificateRootCa(cert *x509.Certificate) (isRootCa bool, err error) {
	if cert == nil {
		return false, tracederrors.TracedErrorNil("cert")
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
