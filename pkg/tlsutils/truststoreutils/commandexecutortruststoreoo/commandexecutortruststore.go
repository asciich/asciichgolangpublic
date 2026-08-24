// Package commandexecutortruststoreoo provides command-based implementation for trust store operations.
package commandexecutortruststoreoo

import (
	"context"
	"encoding/pem"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorTrustStore struct {
	commandExecutor commandexecutorinterfaces.CommandExecutor
	useSudo         bool
}

func NewCommandExecutorTrustStore(commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool) (*CommandExecutorTrustStore, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}
	return &CommandExecutorTrustStore{
		commandExecutor: commandExecutor,
		useSudo:         useSudo,
	}, nil
}

func (c *CommandExecutorTrustStore) AddCaCertificateFromString(ctx context.Context, caCertPEM string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if caCertPEM == "" {
		return tracederrors.TracedErrorEmptyString("caCertPEM")
	}

	block, _ := pem.Decode([]byte(caCertPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("invalid PEM certificate string")
	}

	// Write the cert to a temporary file on the target machine
	tempPath := "/tmp/truststore-cert-install.crt"
	err := commandexecutorfile.WriteBytes(ctx, c.commandExecutor, tempPath, []byte(caCertPEM), &filesoptions.WriteOptions{})
	if err != nil {
		return fmt.Errorf("failed to write certificate to temp file: %w", err)
	}

	// Install from the temp file
	err = c.AddCaCertificateFromFile(ctx, tempPath)
	if err != nil {
		return err
	}

	// Clean up temp file
	_, _ = c.commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"rm", "-f", tempPath},
	})

	return nil
}
