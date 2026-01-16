package commandexecutorgitoo

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) RunGitCommand(ctx context.Context, gitCommand []string) (commandOutput *commandoutput.CommandOutput, err error) {
	if len(gitCommand) <= 0 {
		return nil, tracederrors.TracedError("gitCommand has no elements")
	}

	path, err := g.GetPath()
	if err != nil {
		return nil, err
	}

	commandExecutor, err := g.GetCommandExecutor()
	if err != nil {
		return nil, err
	}

	commandToUse := append([]string{"git", "-C", path}, gitCommand...)

	return commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: commandToUse,
		},
	)
}

func (g *GitRepository) RunGitCommandAndGetStdoutAsString(ctx context.Context, command []string) (stdout string, err error) {
	commandOutput, err := g.RunGitCommand(ctx, command)
	if err != nil {
		return "", err
	}

	return commandOutput.GetStdoutAsString()
}

func (g *GitRepository) RunGitCommandAndGetStdoutAsLines(ctx context.Context, command []string) (lines []string, err error) {
	if command == nil {
		return nil, tracederrors.TracedErrorNil("command")
	}

	output, err := g.RunGitCommand(ctx, command)
	if err != nil {
		return nil, err
	}

	lines, err = output.GetStdoutAsLines(true)
	if err != nil {
		return nil, err
	}

	return lines, nil
}
