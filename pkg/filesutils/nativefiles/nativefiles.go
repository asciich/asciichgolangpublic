package nativefiles

import (
	"context"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Create(ctx context.Context, path string) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if IsFile(contextutils.WithSilent(ctx), path) {
		logging.LogInfoByCtxf(ctx, "File '%s' already exists. Skip create.", path)
	} else {
		err := CreateDirectory(ctx, filepath.Dir(path))
		if err != nil {
			return err
		}

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

func CreateDirectory(ctx context.Context, path string) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if IsDir(contextutils.WithSilent(ctx), path) {
		logging.LogInfoByCtxf(ctx, "Directory '%s' already exists. Skip create.", path)
	} else {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to create directory '%s': %w", path, err)
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

// Delete a file or directory.
// Directories are deleted recursively.
func Delete(ctx context.Context, pathToDelete string, options *filesoptions.DeleteOptions) error {
	if pathToDelete == "" {
		return tracederrors.TracedErrorEmptyString("pathToDelete")
	}

	ctxSilent := contextutils.WithSilent(ctx)
	if Exists(ctxSilent, pathToDelete) {
		isDir := IsDir(ctxSilent, pathToDelete)
		var isFile bool
		if !isDir {
			isFile = IsFile(ctxSilent, pathToDelete)
			if !isFile {
				return tracederrors.TracedErrorf("Path to delete '%s' is pointing to something existing but not a file nor a directory.", pathToDelete)
			}
		}

		if isDir {
			err := os.RemoveAll(pathToDelete)
			if err != nil {
				return tracederrors.TracedErrorf("Delete '%s' failed: %w", pathToDelete, err)
			}
		} else {
			if options != nil && options.UseSudo {
				_, err := commandexecutorexec.RunCommand(ctx, &parameteroptions.RunCommandOptions{
					Command: []string{"sudo", "rm", pathToDelete},
				})
				if err != nil {
					return tracederrors.TracedErrorf("Delete '%s' using sudo failed: %w", pathToDelete, err)
				}
			} else {
				err := os.Remove(pathToDelete)
				if err != nil {
					return tracederrors.TracedErrorf("Delete '%s' failed: %w", pathToDelete, err)
				}
			}
		}

		if isDir {
			logging.LogChangedByCtxf(ctx, "Deleted directory '%s'.", pathToDelete)
		} else if isFile {
			logging.LogChangedByCtxf(ctx, "Deleted file '%s'.", pathToDelete)
		} else {
			return tracederrors.TracedErrorf("Unknown deletion for '%s'", pathToDelete)
		}
	} else {
		logging.LogInfoByCtxf(ctx, "File '%s' already absent. Skip delete.", pathToDelete)
	}

	return nil
}
