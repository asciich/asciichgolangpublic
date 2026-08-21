package commandexecutortruststoreoo

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func (c *CommandExecutorTrustStore) ListCaCertificates(ctx context.Context) ([]*x509.Certificate, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	certPaths := c.getResolvedCertPaths()
	var certificates []*x509.Certificate

	for _, certPath := range certPaths {
		certs, err := c.readCertificatesFromPath(ctx, certPath)
		if err != nil {
			continue
		}
		certificates = append(certificates, certs...)
	}

	return certificates, nil
}

// getResolvedCertPaths returns only the directories where the system stores
// the final/resolved trusted certificates (output of update-ca-certificates).
// These are used for listing only.
func (c *CommandExecutorTrustStore) getResolvedCertPaths() []string {
	return []string{
		"/etc/ssl/certs",
		"/etc/pki/tls/certs",
	}
}

func (c *CommandExecutorTrustStore) readCertificatesFromPath(ctx context.Context, path string) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate

	exists, err := commandexecutorfile.Exists(ctx, c.commandExecutor, path)
	if err != nil || !exists {
		return nil, fmt.Errorf("path does not exist: %s", path)
	}

	isDir, err := c.isPathDirectory(ctx, path)
	if err != nil {
		return nil, err
	}

	if isDir {
		entries, err := c.readDirEntries(ctx, path)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			ext := strings.ToLower(getFileExt(entry))
			if ext != ".crt" && ext != ".cer" && ext != ".pem" && ext != ".ca" {
				continue
			}

			certPath := path + "/" + entry
			cert, err := c.loadCertificateFromFile(ctx, certPath)
			if err != nil {
				continue
			}
			certificates = append(certificates, cert)
		}
	} else {
		cert, err := c.loadCertificateFromFile(ctx, path)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, cert)
	}

	return certificates, nil
}

func (c *CommandExecutorTrustStore) isPathDirectory(ctx context.Context, path string) (bool, error) {
	cmd := fmt.Sprintf("test -d %s && echo yes || echo no", path)
	output, err := c.commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"sh", "-c", cmd},
	})
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(output) == "yes", nil
}

func (c *CommandExecutorTrustStore) readDirEntries(ctx context.Context, path string) ([]string, error) {
	cmd := fmt.Sprintf("ls -1 %s 2>/dev/null", path)
	output, err := c.commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"sh", "-c", cmd},
	})
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var entries []string
	for _, line := range lines {
		if line != "" {
			entries = append(entries, line)
		}
	}
	return entries, nil
}

func (c *CommandExecutorTrustStore) loadCertificateFromFile(ctx context.Context, path string) (*x509.Certificate, error) {
	data, err := commandexecutorfile.ReadAsBytes(c.commandExecutor, path)
	if err != nil {
		return nil, err
	}

	remaining := data
	for {
		var block *pem.Block
		block, remaining = pem.Decode(remaining)
		if block == nil {
			break
		}

		if block.Type == "CERTIFICATE" {
			x509Cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				continue
			}
			return x509Cert, nil
		}
	}

	return nil, fmt.Errorf("no valid certificate found in file: %s", path)
}

func getFileExt(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) > 0 {
		return "." + parts[len(parts)-1]
	}
	return ""
}
