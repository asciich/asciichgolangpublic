package asciichgolangpublic

type CommandExecutorBase struct {
	parentCommandExecutorForBaseClass CommandExecutor
}

func NewCommandExecutorBase() (c *CommandExecutorBase) {
	return new(CommandExecutorBase)
}

func (c *CommandExecutorBase) GetParentCommandExecutorForBaseClass() (parentCommandExecutorForBaseClass CommandExecutor, err error) {

	return c.parentCommandExecutorForBaseClass, nil
}

func (c *CommandExecutorBase) MustGetParentCommandExecutorForBaseClass() (parentCommandExecutorForBaseClass CommandExecutor) {
	parentCommandExecutorForBaseClass, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return parentCommandExecutorForBaseClass
}

func (c *CommandExecutorBase) MustRunCommandAndGetStdoutAsBytes(options *RunCommandOptions) (stdout []byte) {
	stdout, err := c.RunCommandAndGetStdoutAsBytes(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorBase) MustRunCommandAndGetStdoutAsFloat64(options *RunCommandOptions) (stdout float64) {
	stdout, err := c.RunCommandAndGetStdoutAsFloat64(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorBase) MustRunCommandAndGetStdoutAsString(options *RunCommandOptions) (stdout string) {
	stdout, err := c.RunCommandAndGetStdoutAsString(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorBase) MustRunCommandandGetStdoutAsLines(options *RunCommandOptions) (stdoutLines []string) {
	stdoutLines, err := c.RunCommandandGetStdoutAsLines(options)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return stdoutLines
}

func (c *CommandExecutorBase) MustSetParentCommandExecutorForBaseClass(parentCommandExecutorForBaseClass CommandExecutor) {
	err := c.SetParentCommandExecutorForBaseClass(parentCommandExecutorForBaseClass)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsBytes(options *RunCommandOptions) (stdout []byte, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return nil, err
	}

	output, err := parent.RunCommand(options)
	if err != nil {
		return nil, err
	}

	stdout, err = output.GetStdoutAsBytes()
	if err != nil {
		return nil, err
	}

	return stdout, nil
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsFloat64(options *RunCommandOptions) (stdout float64, err error) {
	if options == nil {
		return -1, TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return -1, err
	}

	output, err := parent.RunCommand(options)
	if err != nil {
		return -1, err
	}

	stdout, err = output.GetStdoutAsFloat64()
	if err != nil {
		return -1, err
	}

	return stdout, nil
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsString(options *RunCommandOptions) (stdout string, err error) {
	if options == nil {
		return "", TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return "", err
	}

	stdoutBytes, err := parent.RunCommandAndGetStdoutAsBytes(options)
	if err != nil {
		return "", err
	}

	stdout = string(stdoutBytes)

	return stdout, nil
}

func (c *CommandExecutorBase) RunCommandandGetStdoutAsLines(options *RunCommandOptions) (stdoutLines []string, err error) {
	if options == nil {
		return nil, TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return nil, err
	}

	output, err := parent.RunCommand(options)
	if err != nil {
		return nil, err
	}

	stdoutLines, err = output.GetStdoutAsLines()
	if err != nil {
		return nil, err
	}

	return stdoutLines, nil
}

func (c *CommandExecutorBase) SetParentCommandExecutorForBaseClass(parentCommandExecutorForBaseClass CommandExecutor) (err error) {
	c.parentCommandExecutorForBaseClass = parentCommandExecutorForBaseClass

	return nil
}
