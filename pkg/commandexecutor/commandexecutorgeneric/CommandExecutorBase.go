package commandexecutorgeneric

import (
	"context"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CommandExecutorBase struct {
	parentCommandExecutorForBaseClass commandexecutorinterfaces.CommandExecutor
}

func NewCommandExecutorBase() (c *CommandExecutorBase) {
	return new(CommandExecutorBase)
}

func (c *CommandExecutorBase) GetParentCommandExecutorForBaseClass() (parentCommandExecutorForBaseClass commandexecutorinterfaces.CommandExecutor, err error) {
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

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsBytes(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout []byte, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return nil, err
	}

	output, err := parent.RunCommand(ctx, options)
	if err != nil {
		return nil, err
	}

	stdout, err = output.GetStdoutAsBytes()
	if err != nil {
		return nil, err
	}

	return stdout, nil
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsFloat64(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout float64, err error) {
	if options == nil {
		return -1, tracederrors.TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return -1, err
	}

	output, err := parent.RunCommand(ctx, options)
	if err != nil {
		return -1, err
	}

	stdout, err = output.GetStdoutAsFloat64()
	if err != nil {
		return -1, err
	}

	return stdout, nil
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsInt64(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout int64, err error) {
	stdoutString, err := c.RunCommandAndGetStdoutAsString(ctx, options)
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

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsLines(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdoutLines []string, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return nil, err
	}

	output, err := parent.RunCommand(ctx, options)
	if err != nil {
		return nil, err
	}

	stdoutLines, err = output.GetStdoutAsLines(options.RemoveLastLineIfEmpty)
	if err != nil {
		return nil, err
	}

	return stdoutLines, nil
}

func (c *CommandExecutorBase) RunCommandAndGetStdoutAsString(ctx context.Context, options *parameteroptions.RunCommandOptions) (stdout string, err error) {
	if options == nil {
		return "", tracederrors.TracedErrorNil("options")
	}

	parent, err := c.GetParentCommandExecutorForBaseClass()
	if err != nil {
		return "", err
	}

	stdoutBytes, err := parent.RunCommandAndGetStdoutAsBytes(ctx, options)
	if err != nil {
		return "", err
	}

	stdout = string(stdoutBytes)

	return stdout, nil
}

func (c *CommandExecutorBase) SetParentCommandExecutorForBaseClass(parentCommandExecutorForBaseClass commandexecutorinterfaces.CommandExecutor) (err error) {
	c.parentCommandExecutorForBaseClass = parentCommandExecutorForBaseClass

	return nil
}
