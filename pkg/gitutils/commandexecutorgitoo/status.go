package commandexecutorgitoo

import (
	"context"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func (g *GitRepository) GetGitStatusOutput(ctx context.Context) (output string, err error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Get git status for repository '%s' on '%s' started.", path, hostDescription)

	output, err = g.RunGitCommandAndGetStdoutAsString(
		ctx,
		[]string{"status", "--porcelain"},
	)
	if err != nil {
		return "", err
	}

	output = strings.TrimSpace(output)

	logging.LogInfoByCtxf(ctx, "Get git status for repository '%s' on '%s' finished.", path, hostDescription)

	return output, nil
}

func (g *GitRepository) HasUncommittedChanges(ctx context.Context) (hasUncommitedChanges bool, err error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	commandExecutor, err := g.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	commandOutput, err := commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"cd '%s' && git diff && git diff --cached && git status --porcelain",
					path,
				),
			},
		},
	)
	if err != nil {
		return false, err
	}

	isEmpty, err := commandOutput.IsStdoutAndStderrEmpty()
	if err != nil {
		return false, err
	}

	if !isEmpty {
		hasUncommitedChanges = true
	}

	if hasUncommitedChanges {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' has uncommited changes.", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' has no uncommited changes.", path, hostDescription)
	}

	return hasUncommitedChanges, nil
}
