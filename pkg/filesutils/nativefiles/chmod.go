package nativefiles

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetAccessPermissions(path string, useSudo bool) (int, error) {
	if path == "" {
		return 0, tracederrors.TracedErrorEmptyString("path")
	}

	var perm int
	var err error
	if useSudo {
		perm, err = commandexecutorfile.GetAccessPermissions(commandexecutorexecoo.Exec(), path, useSudo)
		if err != nil {
			return 0, err
		}
	} else {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return 0, tracederrors.TracedErrorf("Unable to get fileInfo of '%s': %w", path, err)
		}

		perm = int(fileInfo.Mode().Perm())
	}

	return perm, nil
}

func GetAccessPermissionsString(path string) (string, error) {
	permissions, err := GetAccessPermissions(path, false)
	if err != nil {
		return "", err
	}

	return unixfilepermissionsutils.GetPermissionString(permissions)
}

func Chmod(ctx context.Context, path string, options *filesoptions.ChmodOptions) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	toSet, err := options.GetPermissions()
	if err != nil {
		return err
	}

	current, err := GetAccessPermissions(path, options.UseSudo)
	if err != nil {
		return err
	}

	toSetString, err := options.GetPermissionsString()
	if err != nil {
		return err
	}

	if current == toSet {
		logging.LogInfoByCtxf(ctx, "Access permissions of '%s' are already set to '%s' = '%o'", path, toSetString, toSet)
	} else {
		if options.UseSudo {
			_, err = commandexecutorexec.RunCommand(ctx, &parameteroptions.RunCommandOptions{
				Command: []string{"sudo", "chmod", toSetString, path},
			})
			if err != nil {
				return err
			}
		} else {
			err = os.Chmod(path, os.FileMode(toSet))
			if err != nil {
				return tracederrors.TracedErrorf("Failed to set access permissions of '%s' to '%s' = '%o': %w", path, toSetString, toSet, err)
			}
		}

		logging.LogChangedByCtxf(ctx, "Access permissions of '%s' set to '%s' = '%o'", path, toSetString, toSet)
	}

	return nil
}
