package truststoreinterfaces

import (
	"context"
	"crypto/x509"
)

// TrustStore defines the contract for any trust store implementation.
type TrustStore interface {
	// AddCaCertificateFromFile installs a CA certificate into the trust store from a file path.
	AddCaCertificateFromFile(ctx context.Context, path string) error

	// AddCaCertificateFromString installs a CA certificate into the trust store from a PEM-encoded string.
	AddCaCertificateFromString(ctx context.Context, caCertPEM string) error

	// DeleteCaCertificateBySerial removes a CA certificate from the trust store identified by its serial number.
	DeleteCaCertificateBySerial(ctx context.Context, serialNumber string) error

	// DeleteCaCertificatesByCommonName removes all CA certificates from the trust store matching the given common name.
	DeleteCaCertificatesByCommonName(ctx context.Context, commonName string) error

	// ListCaCertificates returns all CA certificates currently in the trust store.
	ListCaCertificates(ctx context.Context) ([]*x509.Certificate, error)
}
