package commandexecutorhostsutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// SetHostName sets the OS hostname of the host to the given value using
// "hostnamectl set-hostname". The operation is skipped when the hostname is
// already set (idempotent).
func SetHostName(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, hostname string, options *hostsutilsoptions.SetHostnameOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if hostname == "" {
		return tracederrors.TracedErrorEmptyString("hostname")
	}

	if options == nil {
		options = &hostsutilsoptions.SetHostnameOptions{}
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set hostname to '%s' on '%s' started.", hostname, hostDescription)

	currentHostname, err := commandExecutor.RunCommandAndGetStdoutAsString(
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

	cmd := []string{"hostnamectl", "set-hostname", hostname}
	if options.UseSudo {
		cmd = append([]string{"sudo"}, cmd...)
	}

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: cmd,
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Set hostname from '%s' to '%s' on '%s'.", currentHostname, hostname, hostDescription)
	logging.LogInfoByCtxf(ctx, "Set hostname to '%s' on '%s' finished.", hostname, hostDescription)

	return nil
}
