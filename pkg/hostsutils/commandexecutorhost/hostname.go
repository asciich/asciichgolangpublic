package commandexecutorhost

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// SetHostName sets the OS hostname of the host to the given value using
// "hostnamectl set-hostname". The operation is skipped when the hostname is
// already set (idempotent).
func (c *CommandExecutorHost) SetHostName(ctx context.Context, hostname string) error {
	if hostname == "" {
		return tracederrors.TracedErrorEmptyString("hostname")
	}

	hostDescription, err := c.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set hostname to '%s' on '%s' started.", hostname, hostDescription)

	currentHostname, err := c.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"hostnamectl", "--static", "hostname"},
		},
	)
	if err != nil {
		return err
	}
	currentHostname = strings.TrimSpace(currentHostname)

	if currentHostname == hostname {
		logging.LogInfoByCtxf(ctx, "Hostname on '%s' is already set to '%s'. Skip setting hostname.", hostDescription, hostname)
		logging.LogInfoByCtxf(ctx, "Set hostname to '%s' on '%s' finished.", hostname, hostDescription)
		return nil
	}

	_, err = c.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"sudo", "hostnamectl", "set-hostname", hostname},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Set hostname from '%s' to '%s' on '%s'.", currentHostname, hostname, hostDescription)
	logging.LogInfoByCtxf(ctx, "Set hostname to '%s' on '%s' finished.", hostname, hostDescription)

	return nil
}
