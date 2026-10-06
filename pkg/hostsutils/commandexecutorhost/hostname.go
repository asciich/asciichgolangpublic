package commandexecutorhost

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/commandexecutorhostsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsoptions"
)

// SetHostName sets the OS hostname of the host to the given value using
// "hostnamectl set-hostname". The operation is skipped when the hostname is
// already set (idempotent).
func (c *CommandExecutorHost) SetHostName(ctx context.Context, hostname string, options *hostsutilsoptions.SetHostnameOptions) error {
	return commandexecutorhostsutils.SetHostName(ctx, c, hostname, options)
}
