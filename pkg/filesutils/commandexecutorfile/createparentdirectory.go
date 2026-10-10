package commandexecutorfile

import (
	"context"
	"path"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// CreateParentDirectory ensures the parent directory of 'filePath' exists (including all
// missing parents, like 'mkdir -p').
//
// - Nothing is done if the parent directory already exists (idempotent).
// - Fails if the parent path exists but is not a directory.
// - 'path' (not 'filepath') is used since the target host is not necessarily the local OS.
func CreateParentDirectory(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string, options *filesoptions.CreateOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	if options == nil {
		options = &filesoptions.CreateOptions{}
	}

	dirPath := path.Dir(filePath)

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	// 'test -d' is POSIX and available in BusyBox. Exit code != 0 means "not a directory".
	output, err := commandExecutor.RunCommand(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command:           []string{"test", "-d", dirPath},
			AllowAllExitCodes: true,
		},
	)
	if err != nil {
		return err
	}

	if output.IsExitSuccess() {
		logging.LogInfoByCtxf(ctx, "Parent directory '%s' of '%s' on '%s' already exists.", dirPath, filePath, hostDescription)
		return nil
	}

	command := []string{"mkdir", "-p", "--", dirPath}
	if options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	_, err = commandExecutor.RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: command,
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Parent directory '%s' of '%s' created on '%s'.", dirPath, filePath, hostDescription)

	return nil
}
