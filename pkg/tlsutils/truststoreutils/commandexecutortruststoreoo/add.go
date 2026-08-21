package commandexecutortruststoreoo

import (
	"context"
	"encoding/pem"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (c *CommandExecutorTrustStore) AddCaCertificateFromFile(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	certBytes, err := commandexecutorfile.ReadAsBytes(c.commandExecutor, path)
	if err != nil {
		return fmt.Errorf("failed to read certificate file: %w", err)
	}

	block, _ := pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("invalid PEM certificate in file")
	}

	cmd, err := c.getInstallCommand(path)
	if err != nil {
		return err
	}

	_, err = c.commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"sh", "-c", cmd},
	})
	if err != nil {
		return fmt.Errorf("failed to install certificate: %w", err)
	}

	return nil
}
