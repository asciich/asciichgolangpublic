package commandexecutor

import (
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type CommandExecutorBase struct {
	parentCommandExecutorForBaseClass CommandExecutor
}

func NewCommandExecutorBase() (c *CommandExecutorBase) {
	return new(CommandExecutorBase)
}

func (c *CommandExecutorBase) GetParentCommandExecutorForBaseClass() (parentCommandExecutorForBaseClass CommandExecutor, err error) {
	if c.parentCommandExecutorForBaseClass == nil {
		return nil, tracederrors.TracedError("parent for CommandExecutorBase not set")
	}

	return c.parentCommandExecutorForBaseClass, nil
}

func (c *CommandExecutorBase) IsRunningOnLocalhost() (isRunningOnLocalhost bool, err error) {
	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return false, err
	}

	hostDescriotion, err := parent.GetHostDescription()
	if err != nil {
		return false, err
	}

	return hostDescriotion == "localhost", nil
}

func (c *CommandExecutorBase) MustGetParentCommandExecutorForBaseClass() (parentCommandExecutorForBaseClass CommandExecutor) {
	parentCommandExecutorForBaseClass, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return parentCommandExecutorForBaseClass
}

func (c *CommandExecutorBase) MustIsRunningOnLocalhost() (isRunningOnLocalhost bool) {
	isRunningOnLocalhost, err := c.IsRunningOnLocalhost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isRunningOnLocalhost
}

func (c *CommandExecutorBase) MustRunCommandAndGetStdoutAsBytes(options *parameteroptions.RunCommandOptions) (stdout []byte) {
	stdout, err := c.RunCommandAndGetStdoutAsBytes(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorBase) MustRunCommandAndGetStdoutAsFloat64(options *parameteroptions.RunCommandOptions) (stdout float64) {
	stdout, err := c.RunCommandAndGetStdoutAsFloat64(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorBase) MustRunCommandAndGetStdoutAsInt64(options *parameteroptions.RunCommandOptions) (stdout int64) {
	stdout, err := c.RunCommandAndGetStdoutAsInt64(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorBase) MustRunCommandAndGetStdoutAsLines(options *parameteroptions.RunCommandOptions) (stdoutLines []string) {
	stdoutLines, err := c.RunCommandAndGetStdoutAsLines(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdoutLines
}

func (c *CommandExecutorBase) MustRunCommandAndGetStdoutAsString(options *parameteroptions.RunCommandOptions) (stdout string) {
	stdout, err := c.RunCommandAndGetStdoutAsString(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdout
}

func (c *CommandExecutorBase) MustSetParentCommandExecutorForBaseClass(parentCommandExecutorForBaseClass CommandExecutor) {
	err := c.SetParentCommandExecutorForBaseClass(parentCommandExecutorForBaseClass)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsBytes(options *parameteroptions.RunCommandOptions) (stdout []byte, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
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

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsFloat64(options *parameteroptions.RunCommandOptions) (stdout float64, err error) {
	if options == nil {
		return -1, tracederrors.TracedErrorNil("options")
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

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsInt64(options *parameteroptions.RunCommandOptions) (stdout int64, err error) {
	stdoutString, err := c.RunCommandAndGetStdoutAsString(options)
	if err != nil {
		return -1, err
	}

	stdoutString = strings.TrimSpace(stdoutString)

	stdout, err = strconv.ParseInt(stdoutString, 10, 64)
	if err != nil {
		return -1, err
	}

	return stdout, nil
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsLines(options *parameteroptions.RunCommandOptions) (stdoutLines []string, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return nil, err
	}

	output, err := parent.RunCommand(options)
	if err != nil {
		return nil, err
	}

	stdoutLines, err = output.GetStdoutAsLines(options.RemoveLastLineIfEmpty)
	if err != nil {
		return nil, err
	}

	return stdoutLines, nil
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsString(options *parameteroptions.RunCommandOptions) (stdout string, err error) {
	if options == nil {
		return "", tracederrors.TracedErrorNil("options")
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

func (c *CommandExecutorBase) SetParentCommandExecutorForBaseClass(parentCommandExecutorForBaseClass CommandExecutor) (err error) {
	c.parentCommandExecutorForBaseClass = parentCommandExecutorForBaseClass

	return nil
}
