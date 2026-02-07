package commandexecutorheadscaleoo

import (
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorHeadscale struct {
	commandExecutor commandexecutorinterfaces.CommandExecutor
}

func (c *CommandExecutorHeadscale) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	if c.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}

	return c.commandExecutor, nil
}

func New(commandExectutor commandexecutorinterfaces.CommandExecutor) (headscaleinterfaces.HeadScale, error) {
	if commandExectutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	return &CommandExecutorHeadscale{
		commandExecutor: commandExectutor,
	}, nil
}

func NewOnLocalhost() (headscaleinterfaces.HeadScale, error) {
	exec := commandexecutorexecoo.Exec()
	return New(exec)
}
