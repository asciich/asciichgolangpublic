package truststoreutils

import (
	"context"
	"crypto/x509"

	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils/nativetruststoreoo"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/truststoreutils/truststoreinterfaces"
)

func NewSystemWideTrustStore() (truststoreinterfaces.TrustStore, error) {
	return nativetruststoreoo.NewNativeTrustStore("")
}

func AddCaCertificateFromFile(ctx context.Context, path string) error {
	trustStore, err := NewSystemWideTrustStore()
	if err != nil {
		return err
	}

	return trustStore.AddCaCertificateFromFile(ctx, path)
}

func AddCaCertificateFromString(ctx context.Context, caCertPEM string) error {
	trustStore, err := NewSystemWideTrustStore()
	if err != nil {
		return err
	}

	return trustStore.AddCaCertificateFromString(ctx, caCertPEM)
}

func DeleteCaCertificateBySerial(ctx context.Context, serialNumber string) error {
	trustStore, err := NewSystemWideTrustStore()
	if err != nil {
		return err
	}

	return trustStore.DeleteCaCertificateBySerial(ctx, serialNumber)
}

func DeleteCaCertificatesByCommonName(ctx context.Context, commonName string) error {
	trustStore, err := NewSystemWideTrustStore()
	if err != nil {
		return err
	}

	return trustStore.DeleteCaCertificatesByCommonName(ctx, commonName)
}

func ListCaCertificates(ctx context.Context) ([]*x509.Certificate, error) {
	trustStore, err := NewSystemWideTrustStore()
	if err != nil {
		return nil, err
	}

	return trustStore.ListCaCertificates(ctx)
}
