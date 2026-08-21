package commandexecutortruststoreoo

import (
	"context"
	"crypto/x509"
	"fmt"
	"math/big"
	"runtime"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func (c *CommandExecutorTrustStore) DeleteCaCertificateBySerial(ctx context.Context, serialNumber string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	certs, err := c.ListCaCertificates(ctx)
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}

	cleanSerial := strings.ReplaceAll(serialNumber, ":", "")
	targetSerialInt, ok := new(big.Int).SetString(cleanSerial, 16)
	if !ok {
		return fmt.Errorf("invalid serial number format: %s", serialNumber)
	}

	var found bool
	for _, cert := range certs {
		if cert.SerialNumber != nil && cert.SerialNumber.Cmp(targetSerialInt) == 0 {
			err = c.deleteCertificate(ctx, cert)
			if err != nil {
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

func (c *CommandExecutorTrustStore) DeleteCaCertificatesByCommonName(ctx context.Context, commonName string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	certs, err := c.ListCaCertificates(ctx)
	if err != nil {
		return fmt.Errorf("failed to list certificates: %w", err)
	}

	var found int
	for _, cert := range certs {
		if cert.Subject.CommonName == commonName {
			err = c.deleteCertificate(ctx, cert)
			if err != nil {
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

func (c *CommandExecutorTrustStore) deleteCertificate(ctx context.Context, cert *x509.Certificate) error {
	// Find and remove the certificate from the input directories
	inputPaths := c.getInputCertPaths()

	var removed bool
	for _, inputPath := range inputPaths {
		exists, err := commandexecutorfile.Exists(ctx, c.commandExecutor, inputPath)
		if err != nil || !exists {
			continue
		}

		entries, err := c.readDirEntries(ctx, inputPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			ext := strings.ToLower(getFileExt(entry))
			if ext != ".crt" && ext != ".cer" && ext != ".pem" && ext != ".ca" {
				continue
			}

			certPath := inputPath + "/" + entry
			fileCert, err := c.loadCertificateFromFile(ctx, certPath)
			if err != nil {
				continue
			}

			if fileCert.SerialNumber.Cmp(cert.SerialNumber) == 0 &&
				fileCert.Subject.CommonName == cert.Subject.CommonName {
				err = commandexecutorfile.Delete(ctx, c.commandExecutor, certPath, nil)
				if err != nil {
					return fmt.Errorf("failed to remove certificate file %s: %w", certPath, err)
				}
				removed = true
			}
		}
	}

	if !removed {
		return fmt.Errorf("certificate file not found in input directories for CN=%s", cert.Subject.CommonName)
	}

	// Run update-ca-certificates to refresh the resolved trust store
	cmd, err := c.getUpdateCommand()
	if err != nil {
		return err
	}

	_, err = c.commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"sh", "-c", cmd},
	})
	if err != nil {
		return fmt.Errorf("failed to update trust store: %w", err)
	}

	return nil
}

// getInputCertPaths returns directories where user-added certificates are placed
// before running update-ca-certificates. These are used for add/remove operations.
func (c *CommandExecutorTrustStore) getInputCertPaths() []string {
	return []string{
		"/usr/local/share/ca-certificates",
		"/etc/ca-certificates/trust-source/anchors",
	}
}

func (c *CommandExecutorTrustStore) getUpdateCommand() (string, error) {
	sudoPrefix := ""
	if c.useSudo {
		sudoPrefix = "sudo "
	}

	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("%supdate-ca-certificates --fresh", sudoPrefix), nil
	case "darwin":
		return "", nil // macOS doesn't need a separate update step
	case "windows":
		return "", nil // Windows doesn't need a separate update step
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
