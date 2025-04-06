package x509utils

import (
	"crypto"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type X509CertKeyPair struct {
	Cert *x509.Certificate
	Key  crypto.PrivateKey
}

func (x *X509CertKeyPair) GetX509Certificate() (*x509.Certificate, error) {
	if x.Cert == nil {
		return nil, tracederrors.TracedError("Cert not set")
	}

	return GetX509CertificateDeepCopy(x.Cert), nil
}

func (x *X509CertKeyPair) GetPrivateKey() (crypto.PrivateKey, error) {
	if x.Key == nil {
		return nil, tracederrors.TracedError("Cert not set")
	}

	return x.Key, nil
}

func (x *X509CertKeyPair) IsKeyMatchingCert() (bool, error) {
	cert, err := x.GetX509Certificate()
	if err != nil {
		return false, err
	}

	key, err := x.GetPrivateKey()
	if err != nil {
		return false, err
	}

	isMatching, err := IsCertificateMatchingPrivateKey(cert, key)
	if err != nil {
		return false, err
	}

	return isMatching, nil
}

func (x *X509CertKeyPair) CheckKeyMatchingCert() error {
	isMatching, err := x.IsKeyMatchingCert()
	if err != nil {
		return err
	}

	if !isMatching {
		return tracederrors.TracedError("key does not mach certificate")
	}

	return nil
}

func (x *X509CertKeyPair) GetPublicKey() (crypto.PublicKey, error) {
	privateKey, err := x.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	return GetPublicKeyFromPrivateKey(privateKey)
}