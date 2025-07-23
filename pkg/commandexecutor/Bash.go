package commandexecutor

import (
	"context"

	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type BashService struct {
	CommandExecutorBase
}

// Can be used to run commands in bash on localhost.
func Bash() (b *BashService) {
	return NewBashService()
}

func NewBashService() (b *BashService) {
	b = new(BashService)
	b.SetParentCommandExecutorForBaseClass(b)
	return b
}

func (b *BashService) GetDeepCopy() (deepCopy commandexecutorinterfaces.CommandExecutor) {
	d := NewBashService()

	*d = *b

	deepCopy = d

	return deepCopy
}

func (b *BashService) GetHostDescription() (hostDescription string, err error) {
	return "localhost", err
}

func (b *BashService) RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandexecutorgeneric.CommandOutput, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
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

	commandOutput, err = Exec().RunCommand(ctx, optionsToUse)
	if err != nil {
		return nil, err
	}

	return commandOutput, nil
}

func (b *BashService) RunOneLiner(ctx context.Context, oneLiner string) (output *commandexecutorgeneric.CommandOutput, err error) {
	if oneLiner == "" {
		return nil, tracederrors.TracedErrorEmptyString("oneLiner")
	}

	output, err = b.RunCommand(
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

func (b *BashService) RunOneLinerAndGetStdoutAsLines(ctx context.Context, oneLiner string) (stdoutLines []string, err error) {
	output, err := b.RunOneLiner(ctx, oneLiner)
	if err != nil {
		return nil, err
	}

	stdoutLines, err = output.GetStdoutAsLines(false)
	if err != nil {
		return nil, err
	}

	return stdoutLines, nil
}

func (b *BashService) RunOneLinerAndGetStdoutAsString(ctx context.Context, oneLiner string) (stdout string, err error) {
	output, err := b.RunOneLiner(ctx, oneLiner)
	if err != nil {
		return "", err
	}

	stdout, err = output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return stdout, nil
}
