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
