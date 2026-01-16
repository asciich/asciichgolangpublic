package commandexecutorgitoo

import (
	"context"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) Commit(ctx context.Context, commitOptions *gitparameteroptions.GitCommitOptions) (createdCommit gitinterfaces.GitCommit, err error) {
	if commitOptions == nil {
		return nil, tracederrors.TracedErrorNil("commitOptions")
	}

	commitCommand := []string{"commit"}

	if commitOptions.AllowEmpty {
		commitCommand = append(commitCommand, "--allow-empty")
	}

	if commitOptions.CommitAllChanges {
		commitCommand = append(commitCommand, "--all")
	}

	message, err := commitOptions.GetMessage()
	if err != nil {
		return nil, err
	}

	commitCommand = append(commitCommand, "-m", message)

	_, err = g.RunGitCommand(
		ctx,
		commitCommand,
	)
	if err != nil {
		return nil, err
	}

	createdCommit, err = g.GetCurrentCommit(ctx)
	if err != nil {
		return nil, err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	createdHash, err := createdCommit.GetHash(ctx)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Created commit '%s' in git repository '%s' on host '%s'.", createdHash, path, hostDescription)

	return createdCommit, nil
}

func (g *GitRepository) GetCurrentCommitHash(ctx context.Context) (currentCommitHash string, err error) {
	currentCommitHash, err = g.RunGitCommandAndGetStdoutAsString(ctx, []string{"rev-parse", "HEAD"})
	if err != nil {
		return "", err
	}

	currentCommitHash = strings.TrimSpace(currentCommitHash)

	return currentCommitHash, nil
}

func (g *GitRepository) GetCurrentCommit(ctx context.Context) (gitinterfaces.GitCommit, error) {
	currentCommitHash, err := g.GetCurrentCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	currentCommit, err := g.GetCommitByHash(currentCommitHash)
	if err != nil {
		return nil, err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Current commit of git repository '%s' on host '%s' has hash '%s'.", path, hostDescription, currentCommitHash)

	return currentCommit, nil
}

func (g *GitRepository) GetCommitByHash(hash string) (gitinterfaces.GitCommit, error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	gitCommit := gitgeneric.NewGitCommit()

	err := gitCommit.SetGitRepo(g)
	if err != nil {
		return nil, err
	}

	err = gitCommit.SetHash(hash)
	if err != nil {
		return nil, err
	}

	return gitCommit, nil
}

func (g *GitRepository) CommitHasParentCommitByCommitHash(hash string) (hasParentCommit bool, err error) {
	return false, tracederrors.TracedErrorNotImplemented()
}

func (g *GitRepository) GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (g *GitRepository) GetAuthorStringByCommitHash(hash string) (authorEmail string, err error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (g *GitRepository) GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (g *GitRepository) GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error) {
	return -1, tracederrors.TracedErrorNotImplemented()
}

func (g *GitRepository) GetCommitMessageByCommitHash(hash string) (commitMessage string, err error) {
	if hash == "" {
		return "", tracederrors.TracedErrorEmptyString("hash")
	}

	stdout, err := g.RunGitCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		[]string{"log", "-n", "1", "--pretty=format:%s", hash},
	)
	if err != nil {
		return "", err
	}

	commitMessage = strings.TrimSpace(stdout)

	if commitMessage == "" {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return "", err
		}

		return "", tracederrors.TracedErrorf(
			"Unable to get commit message for hash '%s' in git repository '%s' on host '%s'. commitMessage is empty string after evaluation.",
			hash,
			path,
			hostDescription,
		)
	}

	return commitMessage, nil
}

func (g *GitRepository) GetCommitParentsByCommitHash(ctx context.Context, hash string, options *parameteroptions.GitCommitGetParentsOptions) (commitParents []gitinterfaces.GitCommit, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (g *GitRepository) GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
