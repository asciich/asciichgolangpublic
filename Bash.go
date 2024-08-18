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

func (b *BashService) MustRunOneLiner(oneLiner string, verbose bool) (output *CommandOutput) {
	output, err := b.RunOneLiner(oneLiner, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return output
}

func (b *BashService) MustRunOneLinerAndGetStdoutAsLines(oneLiner string, verbose bool) (stdoutLines []string) {
	stdoutLines, err := b.RunOneLinerAndGetStdoutAsLines(oneLiner, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdoutLines
}

func (b *BashService) MustRunOneLinerAndGetStdoutAsString(oneLiner string, verbose bool) (stdout string) {
	stdout, err := b.RunOneLinerAndGetStdoutAsString(oneLiner, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
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

func (b *BashService) RunOneLiner(oneLiner string, verbose bool) (output *CommandOutput, err error) {
	if oneLiner == "" {
		return nil, TracedErrorEmptyString("oneLiner")
	}

	output, err = b.RunCommand(
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

func (b *BashService) RunOneLinerAndGetStdoutAsLines(oneLiner string, verbose bool) (stdoutLines []string, err error) {
	output, err := b.RunOneLiner(oneLiner, verbose)
	if err != nil {
		return nil, err
	}

	stdoutLines, err = output.GetStdoutAsLines()
	if err != nil {
		return nil, err
	}

	return stdoutLines, nil
}

func (b *BashService) RunOneLinerAndGetStdoutAsString(oneLiner string, verbose bool) (stdout string, err error) {
	output, err := b.RunOneLiner(oneLiner, verbose)
	if err != nil {
		return "", err
	}

	stdout, err = output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return stdout, nil
}
