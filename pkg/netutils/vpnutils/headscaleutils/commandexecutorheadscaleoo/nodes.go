package commandexecutorheadscaleoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscale"
)


func (c *CommandExecutorHeadscale) ListNodeNames(ctx context.Context) ([]string, error) {
	commandExectuor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandexecutorheadscale.ListNodeNames(ctx, commandExectuor)
}
