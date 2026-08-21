package nativetruststoreoo

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/smallstep/truststore"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils/genericx509utils"
)

// NativeTrustStore is an object-oriented native trust store implementation using smallstep/truststore.
type NativeTrustStore struct {
	trustStorePath string
}

// NewNativeTrustStore creates a new native trust store instance.
// Pass "" for the default OS trust store, or a custom path for a specific trust store location.
func NewNativeTrustStore(trustStorePath string) (*NativeTrustStore, error) {
	return &NativeTrustStore{
		trustStorePath: trustStorePath,
	}, nil
}

// AddCaCertificateFromFile installs a CA certificate into the trust store from a file path.
func (n *NativeTrustStore) AddCaCertificateFromFile(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	cert, err := truststore.ReadCertificate(path)
	if err != nil {
		return fmt.Errorf("failed to read certificate from file: %w", err)
	}

	if err := truststore.Install(cert, n.getOptions()...); err != nil {
		return fmt.Errorf("failed to install certificate: %w", err)
	}

	return nil
}

// AddCaCertificateFromString installs a CA certificate into the trust store from a PEM-encoded string.
func (n *NativeTrustStore) AddCaCertificateFromString(ctx context.Context, caCertPEM string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	cert, err := genericx509utils.ReadCertFromString(caCertPEM)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	if err := truststore.Install(cert, n.getOptions()...); err != nil {
		return fmt.Errorf("failed to install certificate: %w", err)
	}

	return nil
}

// DeleteCaCertificateBySerial removes a CA certificate from the trust store identified by its serial number.
func (n *NativeTrustStore) DeleteCaCertificateBySerial(ctx context.Context, serialNumber string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	certs, err := n.ListCaCertificates(ctx)
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}

	var found bool
	for _, cert := range certs {
		serial, err := genericx509utils.GetSerialNumberAsHexColonSeparated(cert)
		if err != nil {
			continue
		}

		if serial == serialNumber {
			if err := truststore.Uninstall(cert, n.getOptions()...); err != nil {
				return fmt.Errorf("failed to uninstall certificate: %w", err)
			}
			found = true
		}
	}

	if !found {
		return fmt.Errorf("certificate with serial number %s not found", serialNumber)
	}

	return nil
}

// DeleteCaCertificatesByCommonName removes all CA certificates from the trust store matching the given common name.
func (n *NativeTrustStore) DeleteCaCertificatesByCommonName(ctx context.Context, commonName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	certs, err := n.ListCaCertificates(ctx)
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}

	var found int
	for _, cert := range certs {
		if cert.Subject.CommonName == commonName {
			if err := truststore.Uninstall(cert, n.getOptions()...); err != nil {
				return fmt.Errorf("failed to uninstall certificate: %w", err)
			}
			found++
		}
	}

	if found == 0 {
		return fmt.Errorf("no certificates found with common name %s", commonName)
	}

	return nil
}

// ListCaCertificates returns all CA certificates currently in the trust store.
func (n *NativeTrustStore) ListCaCertificates(ctx context.Context) ([]*x509.Certificate, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var certPaths []string

	if n.trustStorePath != "" {
		certPaths = []string{n.trustStorePath}
	} else {
		certPaths = n.getDefaultCertPaths()
	}

	var certificates []*x509.Certificate

	for _, certPath := range certPaths {
		certs, err := n.readCertificatesFromPath(certPath)
		if err != nil {
			continue
		}
		certificates = append(certificates, certs...)
	}

	// Deduplicate by SHA-256 fingerprint
	seen := make(map[string]struct{})
	var unique []*x509.Certificate
	for _, cert := range certificates {
		fp := sha256.Sum256(cert.Raw)
		key := hex.EncodeToString(fp[:])
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			unique = append(unique, cert)
		}
	}

	return unique, nil
}

func (n *NativeTrustStore) getOptions() []truststore.Option {
	var opts []truststore.Option
	return opts
}

func (n *NativeTrustStore) getDefaultCertPaths() []string {
	switch runtime.GOOS {
	case "linux":
		return []string{
			"/etc/ssl/certs",
			"/etc/pki/tls/certs",
			"/usr/local/share/ca-certificates",
			"/etc/ca-certificates/trust-source/anchors",
		}
	case "darwin":
		return []string{
			"/etc/ssl/certs",
			"/System/Library/Keychains",
			"/Library/Keychains",
		}
	case "windows":
		return []string{}
	default:
		return []string{}
	}
}

func (n *NativeTrustStore) readCertificatesFromPath(certPath string) ([]*x509.Certificate, error) {
	info, err := os.Stat(certPath)
	if err != nil {
		return nil, err
	}

	var certificates []*x509.Certificate

	if info.IsDir() {
		entries, err := os.ReadDir(certPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filePath := filepath.Join(certPath, entry.Name())
			certs, err := n.loadCertificatesFromFile(filePath)
			if err != nil {
				continue
			}

			certificates = append(certificates, certs...)
		}
	} else {
		certs, err := n.loadCertificatesFromFile(certPath)
		if err != nil {
			return nil, err
		}

		certificates = append(certificates, certs...)
	}

	return certificates, nil
}

func (n *NativeTrustStore) loadCertificatesFromFile(path string) ([]*x509.Certificate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	certs, err := genericx509utils.ReadCertsFromBytes(data)
	if err != nil {
		return nil, err
	}

	return certs, nil
}
