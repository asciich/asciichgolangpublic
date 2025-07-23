package commandexecutor

import (
	"context"

	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type PowerShellService struct {
	CommandExecutorBase
}

func NewPowerShell() (p *PowerShellService) {
	return new(PowerShellService)
}

func NewPowerShellService() (p *PowerShellService) {
	return new(PowerShellService)
}

func PowerShell() (p *PowerShellService) {
	return NewPowerShell()
}

func (b *PowerShellService) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *CommandOutput, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	optionsToUse := options.GetDeepCopy()

	joinedCommand, err := optionsToUse.GetJoinedCommand()
	if err != nil {
		return nil, err
	}

	powerShellCommand := []string{
		"powershell",
		joinedCommand,
	}

	if optionsToUse.RunAsRoot {
		joined, err := shelllinehandler.Join([]string{
			"Start-Process",
			"powershell",
			"-Verb",
			"runAs",
			joinedCommand,
		})
		if err != nil {
			return nil, err
		}

		powerShellCommand = []string{
			"powershell",
			joined,
		}
	}

	optionsToUse.Command = powerShellCommand

	commandOutput, err = Exec().RunCommand(ctx, optionsToUse)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (p *PowerShellService) RunOneLiner(ctx context.Context, oneLiner string) (output *CommandOutput, err error) {
	if oneLiner == "" {
		return nil, tracederrors.TracedErrorEmptyString("oneLiner")
	}

	output, err = p.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{oneLiner},
		},
	)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (p *PowerShellService) RunOneLinerAndGetStdoutAsString(ctx context.Context, oneLiner string) (stdout string, err error) {
	output, err := p.RunOneLiner(ctx, oneLiner)
	if err != nil {
		return "", err
	}

	stdout, err = output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return stdout, nil
}
