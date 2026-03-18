package commandexecutorbash

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RunCommand(ctx context.Context, options *parameteroptions.RunCommandOptions) (commandOutput *commandoutput.CommandOutput, err error) {
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

	return commandexecutorexec.RunCommand(ctx, optionsToUse)
}

func RunOneLiner(ctx context.Context, oneLiner string) (output *commandoutput.CommandOutput, err error) {
	if oneLiner == "" {
		return nil, tracederrors.TracedErrorEmptyString("oneLiner")
	}

	output, err = RunCommand(
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

func RunOneLinerAndGetStdoutAsLines(ctx context.Context, oneLiner string) (stdoutLines []string, err error) {
	output, err := RunOneLiner(ctx, oneLiner)
	if err != nil {
		return nil, err
	}

	stdoutLines, err = output.GetStdoutAsLines(false)
	if err != nil {
		return nil, err
	}

	return stdoutLines, nil
}

func RunOneLinerAndGetStdoutAsString(ctx context.Context, oneLiner string) (stdout string, err error) {
	output, err := RunOneLiner(ctx, oneLiner)
	if err != nil {
		return "", err
	}

	stdout, err = output.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return stdout, nil
}

func RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
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

	return commandexecutorexec.RunCommandAndGetStdoutAsIoReadCloser(ctx, optionsToUse)
}

func RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
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

	return commandexecutorexec.RunCommandAndGetStdinAsIoWriteCloser(ctx, optionsToUse)
}
