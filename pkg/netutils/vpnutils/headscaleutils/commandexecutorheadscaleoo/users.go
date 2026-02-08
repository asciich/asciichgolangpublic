package commandexecutorheadscaleoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscale"
)

func (c *CommandExecutorHeadscale) CreateUser(ctx context.Context, userName string) error {
	commandExectuor, err := c.GetCommandExecutor()
	if err != nil {
		return err
	}

	return commandexecutorheadscale.CreateUser(ctx, commandExectuor, userName)
}

func (c *CommandExecutorHeadscale) GetUserId(ctx context.Context, userName string) (int, error) {
	commandExectuor, err := c.GetCommandExecutor()
	if err != nil {
		return 0, err
	}

	return commandexecutorheadscale.GetUserId(ctx, commandExectuor, userName)
}

func (c *CommandExecutorHeadscale) ListUserNames(ctx context.Context) ([]string, error) {
	commandExectuor, err := c.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	return commandexecutorheadscale.ListUserNames(ctx, commandExectuor)
}

func (c *CommandExecutorHeadscale) GeneratePreauthKeyForUser(ctx context.Context, userName string) (string, error) {
	commandExectuor, err := c.GetCommandExecutor()
	if err != nil {
		return "", err
	}

	return commandexecutorheadscale.GeneratePreauthKeyForUser(ctx, commandExectuor, userName)
}
