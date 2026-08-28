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
	if hash == "" {
		return false, tracederrors.TracedErrorEmptyString("hash")
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(contextutils.ContextVerbose(), "Check if commit '%s' has a parent commit in git repository '%s' on '%s' started.", hash, path, hostDescription)

	// Use git rev-parse <hash>^@ to get the parents of the commit
	// If the commit has no parents (root commit), this will return nothing
	// If the commit has parents, this will return the parent hash(es)
	stdout, err := g.RunGitCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		[]string{"rev-parse", hash + "^@"},
	)
	if err != nil {
		// If the command fails, the commit might be a root commit with no parents
		// Try another approach: check if rev-parse <hash>^ returns error
		_, err2 := g.RunGitCommandAndGetStdoutAsString(
			contextutils.ContextSilent(),
			[]string{"rev-parse", hash + "^"},
		)
		if err2 != nil {
			// Both commands failed, this is a root commit with no parents
			logging.LogInfoByCtxf(contextutils.ContextVerbose(), "Check if commit '%s' has a parent commit in git repository '%s' on '%s' finished. hasParentCommit=false (root commit).", hash, path, hostDescription)
			return false, nil
		}
		return false, err
	}

	stdout = strings.TrimSpace(stdout)

	// If stdout is empty, the commit has no parents
	// If stdout contains parent hashes, the commit has parents
	hasParentCommit = stdout != ""

	logging.LogInfoByCtxf(contextutils.ContextVerbose(), "Check if commit '%s' has a parent commit in git repository '%s' on '%s' finished. hasParentCommit=%v.", hash, path, hostDescription, hasParentCommit)

	return hasParentCommit, nil
}

func (g *GitRepository) GetAuthorEmailByCommitHash(hash string) (authorEmail string, err error) {
	if hash == "" {
		return "", tracederrors.TracedErrorEmptyString("hash")
	}

	stdout, err := g.RunGitCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		[]string{"log", "-n", "1", "--pretty=format:%ae", hash},
	)
	if err != nil {
		return "", err
	}

	authorEmail = strings.TrimSpace(stdout)

	if authorEmail == "" {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return "", err
		}

		return "", tracederrors.TracedErrorf(
			"Unable to get author email for hash '%s' in git repository '%s' on host '%s'. authorEmail is empty string after evaluation.",
			hash,
			path,
			hostDescription,
		)
	}

	return authorEmail, nil
}

func (g *GitRepository) GetAuthorStringByCommitHash(hash string) (authorString string, err error) {
	if hash == "" {
		return "", tracederrors.TracedErrorEmptyString("hash")
	}

	stdout, err := g.RunGitCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		[]string{"log", "-n", "1", "--pretty=format:%an <%ae>", hash},
	)
	if err != nil {
		return "", err
	}

	authorString = strings.TrimSpace(stdout)

	if authorString == "" {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return "", err
		}

		return "", tracederrors.TracedErrorf(
			"Unable to get author string for hash '%s' in git repository '%s' on host '%s'. authorString is empty string after evaluation.",
			hash,
			path,
			hostDescription,
		)
	}

	return authorString, nil
}

func (g *GitRepository) GetCommitAgeDurationByCommitHash(hash string) (ageDuration *time.Duration, err error) {
	commitTime, err := g.GetCommitTimeByCommitHash(hash)
	if err != nil {
		return nil, err
	}

	duration := time.Since(*commitTime)

	return &duration, nil
}

func (g *GitRepository) GetCommitAgeSecondsByCommitHash(hash string) (ageSeconds float64, err error) {
	ageDuration, err := g.GetCommitAgeDurationByCommitHash(hash)
	if err != nil {
		return -1, err
	}

	ageSeconds = ageDuration.Seconds()

	return ageSeconds, nil
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
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	if options == nil {
		options = parameteroptions.NewGitCommitGetParentsOptions()
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Get parent commits for commit '%s' in git repository '%s' on '%s' started.", hash, path, hostDescription)

	// Use git rev-parse <hash>^@ to get all parent hashes of the commit
	parentHashesStr, err := g.RunGitCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		[]string{"rev-parse", hash + "^@"},
	)
	if err != nil {
		// If the command fails, the commit might be a root commit with no parents
		logging.LogInfoByCtxf(ctx, "Get parent commits for commit '%s' in git repository '%s' on '%s' finished. No parent commits found (root commit).", hash, path, hostDescription)
		return []gitinterfaces.GitCommit{}, nil
	}

	parentHashesStr = strings.TrimSpace(parentHashesStr)

	if parentHashesStr == "" {
		logging.LogInfoByCtxf(ctx, "Get parent commits for commit '%s' in git repository '%s' on '%s' finished. No parent commits found.", hash, path, hostDescription)
		return []gitinterfaces.GitCommit{}, nil
	}

	parentHashes := strings.Split(parentHashesStr, "\n")
	commitParents = []gitinterfaces.GitCommit{}

	for _, parentHash := range parentHashes {
		parentHash = strings.TrimSpace(parentHash)
		if parentHash == "" {
			continue
		}

		parentCommit, err := g.GetCommitByHash(parentHash)
		if err != nil {
			return nil, err
		}
		commitParents = append(commitParents, parentCommit)

		// If IncludeParentsOfParents is true, recursively get parents of this parent
		if options.GetIncludeParentsOfParents() {
			grandparents, err := g.GetCommitParentsByCommitHash(ctx, parentHash, &parameteroptions.GitCommitGetParentsOptions{
				IncludeParentsOfParents: false, // Prevent infinite recursion
			})
			if err != nil {
				return nil, err
			}
			commitParents = append(commitParents, grandparents...)
		}
	}

	logging.LogInfoByCtxf(ctx, "Get parent commits for commit '%s' in git repository '%s' on '%s' finished. Found %d parent commit(s).", hash, path, hostDescription, len(commitParents))

	return commitParents, nil
}

func (g *GitRepository) GetCommitTimeByCommitHash(hash string) (commitTime *time.Time, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	stdout, err := g.RunGitCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		[]string{"log", "-n", "1", "--pretty=format:%ai", hash},
	)
	if err != nil {
		return nil, err
	}

	stdout = strings.TrimSpace(stdout)

	if stdout == "" {
		path, hostDescription, err := g.GetPathAndHostDescription()
		if err != nil {
			return nil, err
		}

		return nil, tracederrors.TracedErrorf(
			"Unable to get commit time for hash '%s' in git repository '%s' on host '%s'. stdout is empty string after evaluation.",
			hash,
			path,
			hostDescription,
		)
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700", stdout)
	if err != nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to parse commit time '%s' for hash '%s': %w",
			stdout,
			hash,
			err,
		)
	}

	return &parsedTime, nil
}
