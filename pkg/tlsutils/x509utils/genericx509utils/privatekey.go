package genericx509utils

import (
	"crypto"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/cryptoutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

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
