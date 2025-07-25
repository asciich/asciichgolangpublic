package nativefiles

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

func Delete(ctx context.Context, pathToDelete string) error {
	if pathToDelete == "" {
		return tracederrors.TracedErrorEmptyString("pathToDelete")
	}

	if Exists(contextutils.WithSilent(ctx), pathToDelete) {
		err := os.Remove(pathToDelete)
		if err != nil {
			return tracederrors.TracedErrorf("Delete file '%s' failed: %w", pathToDelete, err)
		}

		logging.LogChangedByCtxf(ctx, "Deleted file '%s'.", pathToDelete)
	} else {
		logging.LogInfoByCtxf(ctx, "File '%s' already absent. Skip delete.", pathToDelete)
	}

	return nil
}

func WriteString(ctx context.Context, pathToWrite string, content string) error {
	if pathToWrite == "" {
		return tracederrors.TracedErrorEmptyString("pathToWrite")
	}

	err := os.WriteFile(pathToWrite, []byte(content), 0644)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to write to file '%s': %w", pathToWrite, err)
	}

	logging.LogChangedByCtxf(ctx, "Wrote content to file '%s'.", pathToWrite)

	return nil
}

func ReadAsString(ctx context.Context, pathToRead string) (string, error) {
	if pathToRead == "" {
		return "", tracederrors.TracedErrorEmptyString("pathToRead")
	}

	content, err := os.ReadFile(pathToRead)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to read file '%s': %w", pathToRead, err)
	}

	logging.LogInfoByCtxf(ctx, "Read content of file '%s'.", content)

	return string(content), nil
}
