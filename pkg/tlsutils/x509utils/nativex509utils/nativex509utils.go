// Package nativex509utils provides native Go implementation for X509 certificate operations.
package nativex509utils

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/x509options"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ReadCertificateFromFile(ctx context.Context, pathToRead string) (cert *x509.Certificate, err error) {
	if pathToRead == "" {
		return nil, tracederrors.TracedErrorEmptyString("pathToRead")
	}

	content, err := nativefiles.ReadAsBytes(ctx, pathToRead)
	if err != nil {
		return nil, err
	}

	return genericx509utils.ReadCertFromBytes(content)
}

func ReadCertificateFromFileOrStdin(ctx context.Context, pathToRead string) (cert *x509.Certificate, err error) {
	if pathToRead == "" {
		return nil, tracederrors.TracedErrorEmptyString("pathToRead")
	}

	if pathToRead == "-" {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to read from stdin: %w", err)
		}
		return genericx509utils.ReadCertFromBytes(content)
	}

	return ReadCertificateFromFile(ctx, pathToRead)
}

func GeneratePrivateKey(ctx context.Context) (privateKey crypto.PrivateKey, err error) {
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

func GetServersCertificateChain(ctx context.Context, hostname string, port int) (certs []*x509.Certificate, err error) {
	if hostname == "" {
		return nil, tracederrors.TracedErrorEmptyString("hostname")
	}

	if port <= 0 {
		return nil, tracederrors.TracedErrorf("port must be positive but got '%d'", port)
	}

	address := fmt.Sprintf("%s:%d", hostname, port)

	conn, err := tls.Dial("tcp", address, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to dial '%s': %w", address, err)
	}
	defer conn.Close()

	certs = conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, tracederrors.TracedErrorf("No certificates returned by server '%s'", address)
	}

	return certs, nil
}
