package nativegitoo

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (n *NativeGitRepository) Pull(ctx context.Context) (err error) {
	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pull git repository '%s' on '%s' started.", path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return tracederrors.TracedErrorf("Get worktree for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	err = worktree.PullContext(ctx, &git.PullOptions{})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' is already up to date.", path, hostDescription)
		} else {
			return tracederrors.TracedErrorf("Pull git repository '%s' on '%s' failed: %w", path, hostDescription, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Pull git repository '%s' on '%s' finished.", path, hostDescription)

	return nil
}

func (n *NativeGitRepository) Fetch(ctx context.Context) (err error) {
	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Fetch git repository '%s' on '%s' started.", path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	err = repo.FetchContext(ctx, &git.FetchOptions{})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logging.LogInfoByCtxf(ctx, "Git repository '%s' on '%s' is already up to date.", path, hostDescription)
		} else {
			return tracederrors.TracedErrorf("Fetch git repository '%s' on '%s' failed: %w", path, hostDescription, err)
		}
	}

	logging.LogInfoByCtxf(ctx, "Fetch git repository '%s' on '%s' finished.", path, hostDescription)

	return nil
}
