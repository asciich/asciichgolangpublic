package nativegitoo

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (n *NativeGitRepository) CloneRepositoryByPathOrUrl(ctx context.Context, pathOrUrl string) (err error) {
	if pathOrUrl == "" {
		return tracederrors.TracedErrorEmptyString("pathOrUrl")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Cloning git repository '%s' to '%s' on '%s' started.", pathOrUrl, path, hostDescription)

	isInitialized, err := n.IsInitialized(ctx)
	if err != nil {
		return err
	}

	if isInitialized {
		logging.LogInfoByCtxf(ctx, "'%s' is already an initialized git repository on host '%s'. Skip clone.", path, hostDescription)
	} else {
		_, err = git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
			URL: pathOrUrl,
		})
		if err != nil {
			if err == transport.ErrEmptyRemoteRepository {
				_, err = git.PlainInit(path, false)
				if err != nil {
					return tracederrors.TracedErrorf("Init git repository '%s' on '%s' after empty remote clone failed: %w", path, hostDescription, err)
				}

				err = n.SetRemoteUrl(ctx, pathOrUrl)
				if err != nil {
					return err
				}
			} else {
				return tracederrors.TracedErrorf("Clone git repository '%s' to '%s' on '%s' failed: %w", pathOrUrl, path, hostDescription, err)
			}
		}
	}

	logging.LogInfoByCtxf(ctx, "Cloning git repository '%s' to '%s' on '%s' finished.", pathOrUrl, path, hostDescription)

	return nil
}

func (n *NativeGitRepository) IsInitialized(ctx context.Context) (isInitialized bool, err error) {
	path, err := n.GetPath()
	if err != nil {
		return false, err
	}

	_, err = git.PlainOpen(path)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return false, nil
		}
		return false, tracederrors.TracedErrorf("Check if git repository '%s' is initialized failed: %w", path, err)
	}

	return true, nil
}
