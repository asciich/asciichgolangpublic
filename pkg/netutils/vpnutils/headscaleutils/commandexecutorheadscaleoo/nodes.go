package commandexecutorheadscaleoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscale"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
)

func (c *CommandExecutorHeadscale) ListNodeNames(ctx context.Context) ([]string, error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandexecutorheadscale.ListNodeNames(ctx, commandExecutor)
}

func (c *CommandExecutorHeadscale) ListNodes(ctx context.Context) ([]headscaleinterfaces.Node, error) {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandexecutorheadscale.ListNodes(ctx, commandExecutor)
}

func (c *CommandExecutorHeadscale) DeleteNodeByName(ctx context.Context, nodeName string) error {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	return commandexecutorheadscale.DeleteNodeByName(ctx, commandExecutor, nodeName)
}

func (c *CommandExecutorHeadscale) DeleteNodeById(ctx context.Context, id uint64) error {
	commandExecutor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	return commandexecutorheadscale.DeleteNodeById(ctx, commandExecutor, id)
}
