package x509utils

import (
	"context"
	"crypto"
)

type X509CertificateHandler interface {
	CreateRootCaCertificate(ctx context.Context, options *X509CreateCertificateOptions) (caCertAndKey *X509CertKeyPair, err error)
	CreateIntermediateCertificate(ctx context.Context, options *X509CreateCertificateOptions) (intermediateCert *X509CertKeyPair, err error)
	CreateSelfSignedCertificate(ctx context.Context, options *X509CreateCertificateOptions) (selfSignesCertAndKey *X509CertKeyPair, err error)
	CreateSignedIntermediateCertificate(ctx context.Context, options *X509CreateCertificateOptions, rootCaCertAndKey *X509CertKeyPair) (intermediateCertAndKey *X509CertKeyPair, err error)
	CreateSignedEndEndityCertificate(ctx context.Context, options *X509CreateCertificateOptions, caCertAndKey *X509CertKeyPair) (endEndityCertAndKey *X509CertKeyPair, err error)

	GeneratePrivateKey(ctx context.Context) (privateKey crypto.PrivateKey, err error)
}
