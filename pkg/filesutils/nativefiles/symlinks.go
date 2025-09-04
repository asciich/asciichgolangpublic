package nativefiles

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ResolveSymlink(ctx context.Context, symlink string) (string, error) {
	if symlink == "" {
		return "", tracederrors.TracedErrorEmptyString("symlink")
	}

	resolved, err := filepath.EvalSymlinks(symlink)
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to resolve symlink: %w", err)
	}

	return resolved, nil
}

func IsSymlinkTo(ctx context.Context, symlink string, target string) (bool, error) {
	if symlink == "" {
		return false, tracederrors.TracedErrorEmptyString("symlink")
	}

	isSymlink, err := IsSymlink(contextutils.WithSilent(ctx), symlink)
	if err != nil {
		return false, err
	}

	if !isSymlink {
		logging.LogInfoByCtxf(ctx, "'%s' is not a symlink and can therefore not point to '%s'.", symlink, target)
		return false, nil
	}

	resolved, err := ResolveSymlink(ctx, symlink)
	if err != nil {
		return false, err
	}

	isLinkTo := resolved == target

	if isLinkTo {
		logging.LogInfoByCtxf(ctx, "'%s' is a symlink to '%s'.", symlink, target)
	} else {
		logging.LogInfoByCtxf(ctx, "'%s' is not a symlinkt to '%s', pointing to '%s' instead.", symlink, target, resolved)
	}

	return isLinkTo, nil
}

func CreateSymlink(ctx context.Context, target string, symlink string) error {
	if target == "" {
		return tracederrors.TracedErrorEmptyString("target")
	}

	if symlink == "" {
		return tracederrors.TracedErrorEmptyString("symlink")
	}

	if !Exists(contextutils.WithSilent(ctx), target) {
		return tracederrors.TracedErrorf("Failed to create symlink '%s'. Target '%s' does not exist.", symlink, target)
	}

	var createSymlink bool
	if Exists(contextutils.WithSilent(ctx), symlink) {
		isSymlinkTo, err := IsSymlinkTo(contextutils.WithSilent(ctx), symlink, target)
		if err != nil {
			return err
		}
		if isSymlinkTo {
			logging.LogInfoByCtxf(ctx, "'%s' is already a symlink to '%s'.", symlink, target)
		} else {
			err = Delete(ctx, symlink, &filesoptions.DeleteOptions{})
			if err != nil {
				return err
			}
			logging.LogChangedByCtxf(ctx, "Removed symlink '%s' pointing to wrong target. Will be recreated correctly.", symlink)
			createSymlink = true
		}
	} else {
		createSymlink = true
	}

	if createSymlink {
		err := os.Symlink(target, symlink)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to create symlink '%s': %w", symlink, err)
		}
		logging.LogChangedByCtxf(ctx, "Symlink '%s' created to point to '%s'.", symlink, target)
	}

	return nil
}

func IsSymlink(ctx context.Context, pathToCheck string) (bool, error) {
	if pathToCheck == "" {
		return false, tracederrors.TracedErrorEmptyString("pathToCheck")
	}

	fileInfo, err := os.Lstat(pathToCheck)
	if err != nil {
		return false, tracederrors.TracedErrorf("Failed to lstat '%s' to detect if symbolic link: %w", pathToCheck, err)
	}

	isSymlink := fileInfo.Mode()&fs.ModeSymlink != 0

	if isSymlink {
		logging.LogInfoByCtxf(ctx, "%s is a symbolic link.", pathToCheck)
	} else {
		logging.LogInfoByCtxf(ctx, "%s is not a symbolic link.", pathToCheck)
	}

	return isSymlink, nil
}

func IsSymlinkToDirectory(ctx context.Context, pathToCheck string) (bool, error) {
	if pathToCheck == "" {
		return false, tracederrors.TracedErrorEmptyString("pathToCheck")
	}

	isSymlink, err := IsSymlink(contextutils.WithSilent(ctx), pathToCheck)
	if err != nil {
		return false, err
	}

	if !isSymlink {
		logging.LogInfoByCtxf(ctx, "%s is not a symbolic link.", pathToCheck)
		return false, err
	}

	resolved, err := ResolveSymlink(contextutils.WithSilent(ctx), pathToCheck)
	if err != nil {
		return false, err
	}

	if IsDir(contextutils.WithSilent(ctx), resolved) {
		logging.LogInfoByCtxf(ctx, "%s is a symbolic link to the directory '%s'.", pathToCheck, resolved)
		return true, nil
	}

	logging.LogInfoByCtxf(ctx, "%s is a symbolic link to '%s' which is not to a directory.", pathToCheck, resolved)
	return false, nil
}
