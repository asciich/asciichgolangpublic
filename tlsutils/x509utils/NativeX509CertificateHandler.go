package x509utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type NativeX509CertificateHandler struct {
}

func NewNativeX509CertificateHandler() (handler *NativeX509CertificateHandler) {
	return new(NativeX509CertificateHandler)
}

func GetNativeX509CertificateHandler() (Handler X509CertificateHandler) {
	return NewNativeX509CertificateHandler()
}

func (n *NativeX509CertificateHandler) SignCertificate(certToSign *x509.Certificate, certWhichSigns *x509.Certificate, certToSignKey crypto.PublicKey, certWhichSignsPrivateKey crypto.PrivateKey, verbose bool) (signedCert *x509.Certificate, err error) {
	if certToSign == nil {
		return nil, tracederrors.TracedErrorNil("certToSign")
	}

	if certWhichSigns == nil {
		return nil, tracederrors.TracedErrorNil("certWhichSigns")
	}

	if certToSignKey == nil {
		return nil, tracederrors.TracedErrorNil("certToSignKey")
	}

	if certWhichSignsPrivateKey == nil {
		return nil, tracederrors.TracedErrorNil("certWhichSignsPrivateKey")
	}

	certToSign.Issuer = pkix.Name{}

	signedCertData, err := x509.CreateCertificate(rand.Reader, certToSign, certWhichSigns, certToSignKey, certWhichSignsPrivateKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to sign cert: %w", err)
	}

	signedCert, err = LoadCertificateFromDerBytes(signedCertData)
	if err != nil {
		return nil, err
	}

	if verbose {
		logging.LogChangedf(
			"Signed certificate '%s' by '%s'.",
			certToSign.Subject,
			certWhichSigns.Subject,
		)
	}

	return signedCert, nil
}

func (n *NativeX509CertificateHandler) GeneratePrivateKey() (privateKey crypto.PrivateKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to create private key: %w", err)
	}

	return privateKey, nil
}

func (n *NativeX509CertificateHandler) generateAndAddKey(cert *x509.Certificate) (certWithKey *x509.Certificate, privateKey crypto.PrivateKey, err error) {
	if cert == nil {
		return nil, nil, tracederrors.TracedErrorNil("cert")
	}

	generatedKey, err := n.GeneratePrivateKey()
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := GetPublicKeyFromPrivateKey(generatedKey)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, publicKey, generatedKey)
	if err != nil {
		return nil, nil, tracederrors.TracedErrorf("Unable to create private key: %w", err)
	}

	certWithKey, err = LoadCertificateFromDerBytes(certBytes)
	if err != nil {
		return nil, nil, err
	}

	return certWithKey, generatedKey, nil
}

func (n *NativeX509CertificateHandler) CreateSignedEndEndityCertificate(options *X509CreateCertificateOptions, intermediateCert *x509.Certificate, intermediatePrivateKey crypto.PrivateKey, verbose bool) (endEndityCert *x509.Certificate, endEndityKey crypto.PrivateKey, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	if intermediateCert == nil {
		return nil, nil, tracederrors.TracedErrorNil("intermediateCert")
	}

	if intermediatePrivateKey == nil {
		return nil, nil, tracederrors.TracedErrorNil("intermediatePrivateKey")
	}

	endEndityCert, endEndityKey, err = n.CreateEndEndityCertificate(options)
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := GetPublicKeyFromPrivateKey(endEndityKey)
	if err != nil {
		return nil, nil, err
	}

	endEndityCert, err = n.SignCertificate(endEndityCert, intermediateCert, pubKey, intermediatePrivateKey, verbose)
	if err != nil {
		return nil, nil, err
	}

	return endEndityCert, endEndityKey, err
}

func (n *NativeX509CertificateHandler) CreateSignedIntermediateCertificate(options *X509CreateCertificateOptions, caCert *x509.Certificate, caPrivateKey crypto.PrivateKey, verbose bool) (intermediateCert *x509.Certificate, intermediateCertPrivateKey crypto.PrivateKey, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	if caCert == nil {
		return nil, nil, tracederrors.TracedErrorNil("caCert")
	}

	if caPrivateKey == nil {
		return nil, nil, tracederrors.TracedErrorNil("caPrivateKey")
	}

	intermediateCert, intermediateCertPrivateKey, err = n.CreateIntermediateCertificate(options)
	if err != nil {
		return nil, nil, err
	}

	pubKey, err := GetPublicKeyFromPrivateKey(intermediateCertPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	intermediateCert, err = n.SignCertificate(intermediateCert, caCert, pubKey, caPrivateKey, verbose)
	if err != nil {
		return nil, nil, err
	}

	return intermediateCert, intermediateCertPrivateKey, err
}

func (n *NativeX509CertificateHandler) CreateEndEndityCertificate(options *X509CreateCertificateOptions) (cert *x509.Certificate, privateKey crypto.PrivateKey, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	subject, err := options.GetSubjectAsPkixName()
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := options.GetSerialNumberOrDefaultIfUnsetAsStringBigInt()
	if err != nil {
		return nil, nil, err
	}

	validityDuration, err := options.GetValidityDuration()
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()

	ca := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               *subject,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(*validityDuration),
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	cert, privateKey, err = n.generateAndAddKey(ca)
	if err != nil {
		return nil, nil, err
	}

	return cert, privateKey, nil
}

func (n *NativeX509CertificateHandler) CreateIntermediateCertificate(options *X509CreateCertificateOptions) (cert *x509.Certificate, privateKey crypto.PrivateKey, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	subject, err := options.GetSubjectAsPkixName()
	if err != nil {
		return nil, nil, err
	}

	validityDuration, err := options.GetValidityDuration()
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()

	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               *subject,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(*validityDuration),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	cert, privateKey, err = n.generateAndAddKey(ca)
	if err != nil {
		return nil, nil, err
	}

	return cert, privateKey, nil
}

func (n *NativeX509CertificateHandler) CreateRootCaCertificate(options *X509CreateCertificateOptions) (cert *x509.Certificate, privateKey crypto.PrivateKey, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	subject, err := options.GetSubjectAsPkixName()
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := options.GetSerialNumberOrDefaultIfUnsetAsStringBigInt()
	if err != nil {
		return nil, nil, err
	}

	validityDuration, err := options.GetValidityDuration()
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()

	ca := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               *subject,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(*validityDuration),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	cert, privateKey, err = n.generateAndAddKey(ca)
	if err != nil {
		return nil, nil, err
	}

	return cert, privateKey, nil
}

func (n *NativeX509CertificateHandler) CreateSelfSignedCertificate(options *X509CreateCertificateOptions) (selfSignedCert *x509.Certificate, selfSignedCertPrivateKey crypto.PrivateKey, err error) {
	if options == nil {
		return nil, nil, tracederrors.TracedErrorNil("options")
	}

	subject, err := options.GetSubjectAsPkixName()
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := options.GetSerialNumberOrDefaultIfUnsetAsStringBigInt()
	if err != nil {
		return nil, nil, err
	}

	validityDuration, err := options.GetValidityDuration()
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()

	selfSigned := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               *subject,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(*validityDuration),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	selfSignedCert, selfSignedCertPrivateKey, err = n.generateAndAddKey(selfSigned)
	if err != nil {
		return nil, nil, err
	}

	return selfSignedCert, selfSignedCertPrivateKey, nil
}
