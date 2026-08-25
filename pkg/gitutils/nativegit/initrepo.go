package nativegit

import (
	"context"
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func InitializeEmptyGitRepository(ctx context.Context, path string) (ret *git.Repository, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Initialize empty git repository in '%s' started.", path)

	ret, err = git.PlainInit(path, false)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			ret, err = git.PlainOpen(path)
			if err != nil {
				return nil, tracederrors.TracedErrorf("Open existing git repository in '%s' failed: %w", path, err)
			}
		} else {
			return nil, tracederrors.TracedErrorf("Initialize empty git repository in '%s' failed: %w", path, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Initialize empty git repository in '%s' finished.", path)

	return ret, nil
}

func GetRepositoryRootPathByPath(ctx context.Context, path string) (rootPath string, err error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Get repository root path for '%s' started.", path)

	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to open git repository at '%s': %w", path, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to get worktree for repository at '%s': %w", path, err)
	}

	rootPath = worktree.Filesystem.Root()

	logging.LogInfoByCtxf(ctx, "Get repository root path for '%s' finished: '%s'.", path, rootPath)

	return rootPath, nil
}
