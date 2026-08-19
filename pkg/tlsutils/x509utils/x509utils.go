package x509utils

import (
	"context"
	"crypto"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/nativex509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
)

func ReadCertificateFromFile(ctx context.Context, pathToRead string) (*x509.Certificate, error) {
	return nativex509utils.ReadCertificateFromFile(ctx, pathToRead)
}

func GeneratePrivateKey(ctx context.Context) (crypto.PrivateKey, error) {
	return genericx509utils.GeneratePrivateKey(ctx)
}

func CreateRootCaCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
	return genericx509utils.CreateRootCaCertificate(ctx, options)
}

func CreateIntermediateCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
	return genericx509utils.CreateIntermediateCertificate(ctx, options)
}

func CreateSelfSignedCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions) (*genericx509utils.X509CertKeyPair, error) {
	return genericx509utils.CreateSelfSignedCertificate(ctx, options)
}

func CreateSignedIntermediateCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions, rootCaCertAndKey *genericx509utils.X509CertKeyPair) (*genericx509utils.X509CertKeyPair, error) {
	return genericx509utils.CreateSignedIntermediateCertificate(ctx, options, rootCaCertAndKey)
}

func CreateSignedEndEntityCertificate(ctx context.Context, options *x509options.X509CreateCertificateOptions, caCertAndKey *genericx509utils.X509CertKeyPair) (*genericx509utils.X509CertKeyPair, error) {
	return genericx509utils.CreateSignedEndEntityCertificate(ctx, options, caCertAndKey)
}
