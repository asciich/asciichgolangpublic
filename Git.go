package asciichgolangpublic

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitService struct {
}

func Git() (git *GitService) {
	return NewGitService()
}

func NewGitService() (g *GitService) {
	return new(GitService)
}

func (g *GitService) GetRepositoryRootPathByPath(ctx context.Context, path string) (repoRootPath string, err error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	repoRootPath, err = commandexecutorbashoo.Bash().RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"git", "-C", path, "rev-parse", "--show-toplevel"},
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

	exists, err := repoRootDir.Exists(contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return "", err
	}

	if !exists {
		return "", tracederrors.TracedErrorf(
			"internal error: repoRootDir '%s' points to an non existent path after evaluation",
			repoRootPath,
		)
	}

	logging.LogInfoByCtxf(
		ctx,
		"Found git repository root directory '%s' for local path '%s'.",
		repoRootPath,
		path,
	)

	return repoRootPath, nil
}
