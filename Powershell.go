package asciichgolangpublic

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/shell/shelllinehandler"
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

func (b *PowerShellService) RunCommand(options *RunCommandOptions) (commandOutput *CommandOutput, err error) {
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
		powerShellCommand = []string{
			"powershell",
			shelllinehandler.MustJoin([]string{
				"Start-Process",
				"powershell",
				"-Verb",
				"runAs",
				joinedCommand,
			}),
		}
	}

	optionsToUse.Command = powerShellCommand

	commandOutput, err = Exec().RunCommand(optionsToUse)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (p *PowerShellService) MustRunCommand(options *RunCommandOptions) (commandOutput *CommandOutput) {
	commandOutput, err := p.RunCommand(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (p *PowerShellService) MustRunOneLiner(oneLiner string, verbose bool) (output *CommandOutput) {
	output, err := p.RunOneLiner(oneLiner, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return output
}

func (p *PowerShellService) MustRunOneLinerAndGetStdoutAsString(oneLiner string, verbose bool) (stdout string) {
	stdout, err := p.RunOneLinerAndGetStdoutAsString(oneLiner, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdout
}

func (p *PowerShellService) RunOneLiner(oneLiner string, verbose bool) (output *CommandOutput, err error) {
	if oneLiner == "" {
		return nil, tracederrors.TracedErrorEmptyString("oneLiner")
	}

	output, err = p.RunCommand(
		&RunCommandOptions{
			Command:            []string{oneLiner},
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
		},
	)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (p *PowerShellService) RunOneLinerAndGetStdoutAsString(oneLiner string, verbose bool) (stdout string, err error) {
	output, err := p.RunOneLiner(oneLiner, verbose)
	if err != nil {
		return "", err
	}

	stdout, err = output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return stdout, nil
}
