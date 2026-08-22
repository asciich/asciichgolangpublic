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
		err := CreateDirectory(ctx, filepath.Dir(path), &filesoptions.CreateOptions{})
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

func CreateDirectory(ctx context.Context, path string, options *filesoptions.CreateOptions) error {
	if options == nil {
		options = &filesoptions.CreateOptions{}
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if IsDir(contextutils.WithSilent(ctx), path) {
		logging.LogInfoByCtxf(ctx, "Directory '%s' already exists. Skip create.", path)
	} else {
		if options.UseSudo {
			logging.LogInfoByCtxf(ctx, "Going to create directory '%s' using sudo. Please enter your password if asked by sudo.", path)
			_, err := commandexecutorexec.RunCommand(ctx, &parameteroptions.RunCommandOptions{
				Command: []string{"sudo", "mkdir", "-p", path},
			})
			if err != nil {
				return tracederrors.TracedErrorf("Failed to create directory '%s' using sudo: %w", path, err)
			} else {
				logging.LogChangedByCtxf(ctx, "Created directory '%s' using sudo.", path)
			}
		} else {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return tracederrors.TracedErrorf("Failed to create directory '%s': %w", path, err)
			}
			logging.LogChangedByCtxf(ctx, "Created directory '%s'.", path)
		}
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

func CheckAllExists(ctx context.Context, pathsToCheck []string) error {
	if len(pathsToCheck) <= 0 {
		return tracederrors.TracedError("pathsToCheck is empty")
	}

	for _, p := range pathsToCheck {
		err := CheckExists(ctx, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckExists(ctx context.Context, pathToCheck string) error {
	exists := Exists(ctx, pathToCheck)
	if exists {
		return nil
	}

	return tracederrors.TracedErrorf("Path '%s' does not exist.", pathToCheck)
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

// AppendBytes appends bytes to a file.
func AppendBytes(ctx context.Context, path string, toWrite []byte) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if toWrite == nil {
		return tracederrors.TracedErrorNil("toWrite")
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to open file '%s' for append: %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(toWrite)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to append to file '%s': %w", path, err)
	}

	logging.LogChangedByCtxf(ctx, "Appended bytes to file '%s'.", path)

	return nil
}

// AppendString appends a string to a file.
func AppendString(ctx context.Context, path string, toWrite string) error {
	return AppendBytes(ctx, path, []byte(toWrite))
}

// Chown changes the ownership of a file or directory.
func Chown(ctx context.Context, path string, options *parameteroptions.ChownOptions) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	userAndGroup, err := options.GetUserAndOptionallyGroupForChownCommand()
	if err != nil {
		return err
	}

	if options.UseSudo {
		logging.LogInfoByCtxf(ctx, "Chown '%s' to '%s' using sudo started.", path, userAndGroup)
		_, err := commandexecutorexec.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"sudo", "chown", userAndGroup, path},
		})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to chown '%s' to '%s' using sudo: %w", path, userAndGroup, err)
		}
	} else {
		logging.LogInfoByCtxf(ctx, "Chown '%s' to '%s' started.", path, userAndGroup)
		// For non-sudo chown, we'd need to use syscall.Chown which requires root or file owner
		// For now, we'll use the command executor approach even for local
		_, err := commandexecutorexec.RunCommand(ctx, &parameteroptions.RunCommandOptions{
			Command: []string{"chown", userAndGroup, path},
		})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to chown '%s' to '%s': %w", path, userAndGroup, err)
		}
	}

	logging.LogChangedByCtxf(ctx, "Changed ownership of '%s' to '%s'.", path, userAndGroup)

	return nil
}

// GetBaseName returns the base name of a path.
func GetBaseName(path string) (string, error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	return filepath.Base(path), nil
}

// GetParentDirectoryPath returns the parent directory path.
func GetParentDirectoryPath(path string) (string, error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	parentPath := filepath.Dir(path)
	if parentPath == "" {
		return "", tracederrors.TracedError("parentPath is empty string after evaluation.")
	}

	return parentPath, nil
}

// MoveToPath moves a file to a destination path.
func MoveToPath(ctx context.Context, src string, destPath string, useSudo bool) error {
	if src == "" {
		return tracederrors.TracedErrorEmptyString("src")
	}

	if destPath == "" {
		return tracederrors.TracedErrorEmptyString("destPath")
	}

	options := &filesoptions.MoveOptions{
		UseSudo: useSudo,
	}

	return Move(ctx, src, destPath, options)
}
