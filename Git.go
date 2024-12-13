package asciichgolangpublic

import "strings"

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
		return "", TracedErrorEmptyString("path")
	}

	repoRootPath, err = Bash().RunCommandAndGetStdoutAsString(
		&RunCommandOptions{
			Command:            []string{"git", "-C", path, "rev-parse", "--show-toplevel"},
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
		},
	)
	if err != nil {
		return "", err
	}

	repoRootPath = strings.TrimSpace(repoRootPath)

	repoRootDir, err := GetLocalDirectoryByPath(repoRootPath)
	if err != nil {
		return "", err
	}

	exists, err := repoRootDir.Exists(verbose)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", TracedErrorf(
			"internal error: repoRootDir '%s' points to an non existent path after evaluation",
			repoRootPath,
		)
	}

	if verbose {
		LogInfof(
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
		LogGoErrorFatal(err)
	}

	return repoRootPath
}
