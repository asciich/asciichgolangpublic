package asciichgolangpublic

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
		return nil, TracedErrorNil("options")
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
			ShellLineHandler().MustJoin([]string{
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
		LogGoErrorFatal(err)
	}

	return commandOutput
}
