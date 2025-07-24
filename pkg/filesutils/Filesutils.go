package filesutils

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Create(ctx context.Context, path string) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if IsFile(contextutils.WithSilent(ctx), path) {
		logging.LogInfoByCtxf(ctx, "File '%s' already exists. Skip create.", path)
	} else {
		file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to create file '%s': %w", path, err)
		}

		err = file.Close()
		if err != nil {
			return tracederrors.TracedErrorf("Failed to close created file '%s': %w", path, err)
		}

		logging.LogChangedByCtxf(ctx, "Created file '%s'.", path)
	}

	return nil
}

func IsFile(ctx context.Context, pathToCheck string) bool {
	stat, err := os.Stat(pathToCheck)
	if err != nil {
		logging.LogInfoByCtxf(ctx, "'%s' is not a file", pathToCheck)
		return false
	}

	if stat.IsDir() {
		logging.LogInfoByCtxf(ctx, "'%s' is a directory, not a file.", pathToCheck)
		return false
	}

	logging.LogInfoByCtxf(ctx, "'%s' is a file.", pathToCheck)
	return true
}

func IsDir(ctx context.Context, pathToCheck string) bool {
	stat, err := os.Stat(pathToCheck)
	if err != nil {
		logging.LogInfoByCtxf(ctx, "'%s' is not a dirextory", pathToCheck)
		return false
	}

	if stat.IsDir() {
		logging.LogInfoByCtxf(ctx, "'%s' is a directory.", pathToCheck)
		return true
	}

	logging.LogInfoByCtxf(ctx, "'%s' is a file, not a directory.", pathToCheck)
	return false
}

func Exists(ctx context.Context, pathToCheck string) bool {
	_, err := os.Stat(pathToCheck)
	if err != nil {
		logging.LogInfoByCtxf(ctx, "'%s' does not exist in file system.", pathToCheck)
		return false
	}

	logging.LogInfoByCtxf(ctx, "'%s' exist in file system.", pathToCheck)
	return true
}
