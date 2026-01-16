package commandexecutorgitoo

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) GetGitStatusOutput(ctx context.Context) (output string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
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
