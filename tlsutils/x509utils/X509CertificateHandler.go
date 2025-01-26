package x509utils

import (
	"crypto"
	"crypto/x509"
)


type X509CertificateHandler interface {
	CreateRootCertificate(options *X509CreateCertificateOptions) (cert *x509.Certificate, privateKey crypto.PrivateKey , err error)	
}

