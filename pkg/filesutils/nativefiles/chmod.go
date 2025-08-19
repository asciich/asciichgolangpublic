package nativefiles

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetAccessPermissions(path string) (int, error) {
	if path == "" {
		return 0, tracederrors.TracedErrorEmptyString("path")
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, tracederrors.TracedErrorf("Unable to get fileInfo of '%s': %w", path, err)
	}

	perm := fileInfo.Mode().Perm()

	return int(perm), nil
}

func GetAccessPermissionsString(path string) (string, error) {
	permissions, err := GetAccessPermissions(path)
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

	current, err := GetAccessPermissions(path)
	if err != nil {
		return err
	}

	toSetString, err := options.GetPermissionsString()
	if err != nil {
		return err
	}

	if current == toSet {
		logging.LogInfoByCtxf(ctx, "Access permissions of '%s' are already set to '%s'", path, toSetString)
	} else {
		err = os.Chmod(path, os.FileMode(toSet))
		if err != nil {
			return tracederrors.TracedErrorf("Failed to set access permissions of '%s' to '%s': %w", path, toSetString, err)
		}

		logging.LogChangedByCtxf(ctx, "Access permissions of '%s' set to '%s'", path, toSetString)
	}

	return nil
}
