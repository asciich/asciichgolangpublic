package commandexecutorheadscaleoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscale"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func (c *CommandExecutorHeadscale) ListApiKeys(ctx context.Context) ([]headscaleinterfaces.ApiKey, error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandexecutorheadscale.ListApiKeys(ctx, commandExecutor)
}
