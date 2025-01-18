package asciichgolangpublic

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitService struct {
}

func Git() (git *GitService) {
	return NewGitService()
}

func NewGitService() (g *GitService) {
	return new(GitService)
}

func (g *GitService) GetRepositoryRootPathByPath(path string, verbose bool) (repoRootPath string, err error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	repoRootPath, err = commandexecutor.Bash().RunCommandAndGetStdoutAsString(
		&parameteroptions.RunCommandOptions{
			Command:            []string{"git", "-C", path, "rev-parse", "--show-toplevel"},
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
		},
	)
	if err != nil {
		return "", err
	}

	repoRootPath = strings.TrimSpace(repoRootPath)

	repoRootDir, err := files.GetLocalDirectoryByPath(repoRootPath)
	if err != nil {
		return "", err
	}

	exists, err := repoRootDir.Exists(verbose)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", tracederrors.TracedErrorf(
			"internal error: repoRootDir '%s' points to an non existent path after evaluation",
			repoRootPath,
		)
	}

	if verbose {
		logging.LogInfof(
			"Found git repository root directory '%s' for local path '%s'.",
			repoRootPath,
			path,
		)
	}

	return repoRootPath, nil
}

func (g *GitService) MustGetRepositoryRootPathByPath(path string, verbose bool) (repoRootPath string) {
	repoRootPath, err := g.GetRepositoryRootPathByPath(path, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return repoRootPath
}
