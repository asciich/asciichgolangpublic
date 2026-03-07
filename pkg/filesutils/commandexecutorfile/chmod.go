package commandexecutorfile

import (
	"context"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/unixfilepermissionsutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Chmod(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, options *filesoptions.ChmodOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExectuor")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	newPermissions, err := options.GetPermissions()
	if err != nil {
		return err
	}

	currentPermissions, err := GetAccessPermissions(commandExecutor, path)
	if err != nil {
		return err
	}

	newPermissionsString, err := unixfilepermissionsutils.GetPermissionString(newPermissions)
	if err != nil {
		return err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	if newPermissions == currentPermissions {
		logging.LogInfoByCtxf(ctx, "File permissions for '%s' on '%s' are already set to '%s'.", path, hostDescription, newPermissionsString)
	} else {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"chmod", newPermissionsString, path},
			},
		)
		
		logging.LogChangedByCtxf(ctx, "Changed file permissions for '%s' on '%s' to '%s'.",path, hostDescription, newPermissionsString)
	}

	if err != nil {
		return err
	}

	

	return nil
}

func GetAccessPermissions(commandexecutor commandexecutorinterfaces.CommandExecutor, path string) (int, error) {
	if commandexecutor == nil {
		return 0, tracederrors.TracedErrorNil("commandexecutor")
	}

	if path == "" {
		return 0, tracederrors.TracedErrorEmptyString("path")
	}

	output, err := commandexecutor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{"stat", "-c", "%a", path},
		},
	)
	if err != nil {
		return 0, err
	}

	output = strings.TrimSpace(output)

	permissions64, err := strconv.ParseInt(output, 8, 32)
	if err != nil {
		return 0, tracederrors.TracedErrorf("Unable to parse permission string '%s': %w", output, err)
	}

	return int(permissions64), nil
}

func GetAccessPermissionsString(commandexecutor commandexecutorinterfaces.CommandExecutor, path string) (string, error) {
	permissions, err := GetAccessPermissions(commandexecutor, path)
	if err != nil {
		return "", err
	}

	return unixfilepermissionsutils.GetPermissionString(permissions)
}
