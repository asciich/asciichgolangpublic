package asciichgolangpublic

type BashService struct {
	CommandExecutorBase
}

func Bash() (b *BashService) {
	return NewBashService()
}

func NewBashService() (b *BashService) {
	b = new(BashService)
	b.SetParentCommandExecutorForBaseClass(b)
	return b
}

func (b *BashService) MustRunCommand(options *RunCommandOptions) (commandOutput *CommandOutput) {
	commandOutput, err := b.RunCommand(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commandOutput
}

func (b *BashService) RunCommand(options *RunCommandOptions) (commandOutput *CommandOutput, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	optionsToUse := options.GetDeepCopy()

	joinedCommand, err := optionsToUse.GetJoinedCommand()
	if err != nil {
		return nil, err
	}

	bashCommand := []string{
		"bash",
		"-c",
		joinedCommand,
	}
	optionsToUse.Command = bashCommand

	commandOutput, err = Exec().RunCommand(optionsToUse)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}
