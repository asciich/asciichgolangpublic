package x509utils

import (
	"crypto"
	"crypto/x509"
)

type X509CertificateHandler interface {
	CreateRootCaCertificate(options *X509CreateCertificateOptions) (caCert *x509.Certificate, caPrivateKey crypto.PrivateKey, err error)
	CreateIntermediateCertificate(options *X509CreateCertificateOptions) (intermediateCert *x509.Certificate, intermediateCertPrivateKey crypto.PrivateKey, err error)
	CreateSignedIntermediateCertificate(options *X509CreateCertificateOptions, caCert *x509.Certificate, caPrivateKey crypto.PrivateKey, verbose bool) (intermediateCert *x509.Certificate, intermediateCertPrivateKey crypto.PrivateKey, err error)
	CreateSignedEndEndityCertificate(options *X509CreateCertificateOptions, caCert *x509.Certificate, caPrivateKey crypto.PrivateKey, verbose bool) (endEndityCert *x509.Certificate, endEndityCertPrivateKey crypto.PrivateKey, err error)
}
