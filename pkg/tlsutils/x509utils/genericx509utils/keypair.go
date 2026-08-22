package genericx509utils

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/cryptoutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type X509CertKeyPair struct {
	Cert *x509.Certificate
	Key  crypto.PrivateKey
}

func (x *X509CertKeyPair) GetCertificateAsPEMString() (string, error) {
	cert, err := x.GetX509Certificate()
	if err != nil {
		return "", err
	}

	return WriteCertificateAsPEMString(cert)
}

func (x *X509CertKeyPair) GetCertificateAsPEMBytes() ([]byte, error) {
	cert, err := x.GetX509Certificate()
	if err != nil {
		return nil, err
	}

	return WriteCertificateAsPEMBytes(cert)
}

func (x *X509CertKeyPair) GetX509Certificate() (*x509.Certificate, error) {
	if x.Cert == nil {
		return nil, tracederrors.TracedError("Cert not set")
	}

	return GetX509CertificateDeepCopy(x.Cert), nil
}

func (x *X509CertKeyPair) GetPrivateKeyAsPEMString() (string, error) {
	privateKey, err := x.GetPrivateKey()
	if err != nil {
		return "", err
	}

	return cryptoutils.EncodePrivateKeyAsPEMString(privateKey)
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

func (x *X509CertKeyPair) CheckKeyMatchingCertificate() error {
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

	return cryptoutils.GetPublicKeyFromPrivateKey(privateKey)
}

func (x *X509CertKeyPair) WriteCertificatePemToFile(ctx context.Context, toWrite filesinterfaces.File) error {
	if toWrite == nil {
		return tracederrors.TracedErrorNil("toWrite")
	}

	certPemBytes, err := x.GetCertificateAsPEMBytes()
	if err != nil {
		return err
	}

	return toWrite.WriteBytes(ctx, certPemBytes, nil)
}

func (x *X509CertKeyPair) WritePrivateKeyToFile(ctx context.Context, toWrite filesinterfaces.File) error {
	if toWrite == nil {
		return tracederrors.TracedErrorNil("toWrite")
	}

	privateKey, err := x.GetPrivateKey()
	if err != nil {
		return err
	}

	privateKeyBytes, err := encodePrivateKeyAsPEMBytes(privateKey)
	if err != nil {
		return err
	}

	return toWrite.WriteBytes(ctx, privateKeyBytes, nil)
}

func encodePrivateKeyAsPEMBytes(privateKey crypto.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, tracederrors.TracedErrorNil("privateKey")
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to marshal private key: '%w'", err)
	}

	var buf bytes.Buffer
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	err = pem.Encode(io.Writer(&buf), pemBlock)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to PEM encode private key: '%w'", err)
	}

	return buf.Bytes(), nil
}
