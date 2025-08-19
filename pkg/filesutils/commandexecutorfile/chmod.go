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

	chmodString, err := options.GetPermissionsString()
	if err != nil {
		return err
	}

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"chmod", chmodString, path},
		},
	)

	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Chmod '%s' for local file '%s'.", chmodString, path)

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
