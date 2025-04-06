package x509utils

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
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

// To generate certificates the cert library provided by golang uses a cert struct as a template for the cert to create.
// This function will generate this template by the given options.
func getCertTemplate(options *X509CreateCertificateOptions) (certTemplate *x509.Certificate, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	subject, err := options.GetSubjectAsPkixName()
	if err != nil {
		return nil, err
	}

	serialNumber, err := options.GetSerialNumberOrDefaultIfUnsetAsStringBigInt()
	if err != nil {
		return nil, err
	}

	validityDuration, err := options.GetValidityDuration()
	if err != nil {
		return nil, err
	}

	sans, err := options.GetAdditionalSansOrEmptySliceIfUnset()
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()

	certTemplate = &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               *subject,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(*validityDuration),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		DNSNames:              sans,
		BasicConstraintsValid: true,
	}

	return certTemplate, nil
}

func (n *NativeX509CertificateHandler) SignCertificate(ctx context.Context, certToSignAndKey *X509CertKeyPair, signingCertAndKey *X509CertKeyPair) (signedCert *x509.Certificate, err error) {
	if certToSignAndKey == nil {
		return nil, tracederrors.TracedErrorNil("certToSignAndKey")
	}

	if signingCertAndKey == nil {
		return nil, tracederrors.TracedErrorNil("signingCertAndKey")
	}

	err = certToSignAndKey.CheckKeyMatchingCert()
	if err != nil {
		return nil, err
	}

	err = signingCertAndKey.CheckKeyMatchingCert()
	if err != nil {
		return nil, err
	}

	certToSign, err := certToSignAndKey.GetX509Certificate()
	if err != nil {
		return nil, err
	}

	certToSign.Issuer = pkix.Name{}

	certWhichSigns, err := signingCertAndKey.GetX509Certificate()
	if err != nil {
		return nil, err
	}

	certToSignPublicKey, err := certToSignAndKey.GetPublicKey()
	if err != nil {
		return nil, err
	}

	certWhichSignsPrivateKey, err := signingCertAndKey.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	signedCertData, err := x509.CreateCertificate(rand.Reader, certToSign, certWhichSigns, certToSignPublicKey, certWhichSignsPrivateKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to sign cert: %w", err)
	}

	signedCert, err = LoadCertificateFromDerBytes(signedCertData)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Signed certificate '%s' by '%s'.", FormatForLogging(certToSign), FormatForLogging(certWhichSigns))

	return signedCert, nil
}

func (n *NativeX509CertificateHandler) GeneratePrivateKey(ctx context.Context) (privateKey crypto.PrivateKey, err error) {
	privateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to create private key: %w", err)
	}

	return privateKey, nil
}

func (n *NativeX509CertificateHandler) generateAndAddKey(ctx context.Context, cert *x509.Certificate) (certKeyPair *X509CertKeyPair, err error) {
	if cert == nil {
		return nil, tracederrors.TracedErrorNil("cert")
	}

	generatedKey, err := n.GeneratePrivateKey(ctx)
	if err != nil {
		return nil, err
	}

	publicKey, err := GetPublicKeyFromPrivateKey(generatedKey)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, publicKey, generatedKey)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to create private key: %w", err)
	}

	certWithKey, err := LoadCertificateFromDerBytes(certBytes)
	if err != nil {
		return nil, err
	}

	certKeyPair = &X509CertKeyPair{
		Cert: certWithKey,
		Key:  generatedKey,
	}

	err = certKeyPair.CheckKeyMatchingCert()
	if err != nil {
		return nil, err
	}

	return certKeyPair, nil
}

func (n *NativeX509CertificateHandler) CreateSignedEndEndityCertificate(ctx context.Context, options *X509CreateCertificateOptions, intermediateCertAndKey *X509CertKeyPair) (certKeyPair *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if intermediateCertAndKey == nil {
		return nil, tracederrors.TracedErrorNil("intermediateCertAndKey")
	}

	err = intermediateCertAndKey.CheckKeyMatchingCert()
	if err != nil {
		return nil, err
	}

	endEndityCertAndKey, err := n.CreateEndEndityCertificate(ctx, options)
	if err != nil {
		return nil, err
	}

	endEndityCertAndKey.Cert, err = n.SignCertificate(ctx, endEndityCertAndKey, intermediateCertAndKey)
	if err != nil {
		return nil, err
	}

	return endEndityCertAndKey, err
}

func (n *NativeX509CertificateHandler) CreateSignedIntermediateCertificate(ctx context.Context, options *X509CreateCertificateOptions, caCertAndKey *X509CertKeyPair) (intermediateCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if caCertAndKey == nil {
		return nil, tracederrors.TracedErrorNil("caCertAndKey")
	}

	intermediateCertAndKey, err = n.CreateIntermediateCertificate(ctx, options)
	if err != nil {
		return nil, err
	}

	intermediateCertAndKey.Cert, err = n.SignCertificate(ctx, intermediateCertAndKey, caCertAndKey)
	if err != nil {
		return nil, err
	}

	return intermediateCertAndKey, err
}

func (n *NativeX509CertificateHandler) CreateEndEndityCertificate(ctx context.Context, options *X509CreateCertificateOptions) (endEndityCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	certTemplate, err := getCertTemplate(options)
	if err != nil {
		return nil, err
	}

	endEndityCertAndKey, err = n.generateAndAddKey(ctx, certTemplate)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Generated root CA certificate '%s'", FormatForLogging(endEndityCertAndKey.Cert))

	return endEndityCertAndKey, nil
}

func (n *NativeX509CertificateHandler) CreateIntermediateCertificate(ctx context.Context, options *X509CreateCertificateOptions) (intermediateCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	certTemplate, err := getCertTemplate(options)
	if err != nil {
		return nil, err
	}

	certTemplate.IsCA = true

	intermediateCertAndKey, err = n.generateAndAddKey(ctx, certTemplate)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Generated root CA certificate '%s'", FormatForLogging(intermediateCertAndKey.Cert))

	return intermediateCertAndKey, nil
}

func (n *NativeX509CertificateHandler) CreateRootCaCertificate(ctx context.Context, options *X509CreateCertificateOptions) (rootCaCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	certTemplate, err := getCertTemplate(options)
	if err != nil {
		return nil, err
	}

	certTemplate.IsCA = true

	rootCaCertAndKey, err = n.generateAndAddKey(ctx, certTemplate)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Generated root CA certificate '%s'", FormatForLogging(rootCaCertAndKey.Cert))

	return rootCaCertAndKey, nil
}

func (n *NativeX509CertificateHandler) CreateSelfSignedCertificate(ctx context.Context, options *X509CreateCertificateOptions) (selfSignedCertAndKey *X509CertKeyPair, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	certTemplate, err := getCertTemplate(options)
	if err != nil {
		return nil, err
	}

	selfSignedCertAndKey, err = n.generateAndAddKey(ctx, certTemplate)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Generated self signed certificate '%s'", FormatForLogging(selfSignedCertAndKey.Cert))

	return selfSignedCertAndKey, nil
}
