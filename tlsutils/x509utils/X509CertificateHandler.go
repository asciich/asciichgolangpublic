package x509utils

import (
	"crypto"
	"crypto/x509"
)

type X509CertificateHandler interface {
	CreateRootCaCertificate(options *X509CreateCertificateOptions) (caCert *x509.Certificate, caPrivateKey crypto.PrivateKey, err error)
}
