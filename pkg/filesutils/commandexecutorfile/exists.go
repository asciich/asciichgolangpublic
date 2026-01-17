package commandexecutorfile

import (
	"context"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func FileExists(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"test -f '%s' && echo yes || echo no",
					filePath,
				),
			},
		},
	)
	if err != nil {
		return false, err
	}

	output = strings.TrimSpace(output)
	var exists = false
	if output == "yes" {
		exists = true
	} else if output == "no" {
		exists = false
	} else {
		return false, tracederrors.TracedErrorf(
			"Unexpected output when checking for file to exist: '%s'",
			output,
		)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "File '%s' on host '%s' exists.", filePath, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "File '%s' on host '%s' does not exist.", filePath, hostDescription)
	}

	return exists, nil
}


func DirectoryExists(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, directoryPath string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if directoryPath == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"test -d '%s' && echo yes || echo no",
					directoryPath,
				),
			},
		},
	)
	if err != nil {
		return false, err
	}

	output = strings.TrimSpace(output)
	var exists = false
	if output == "yes" {
		exists = true
	} else if output == "no" {
		exists = false
	} else {
		return false, tracederrors.TracedErrorf(
			"Unexpected output when checking for directory to exist: '%s'",
			output,
		)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Directory '%s' on host '%s' exists.", directoryPath, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Directory '%s' on host '%s' does not exist.", directoryPath, hostDescription)
	}

	return exists, nil
}

func Exists(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("filePath")
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsString(
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"test -e '%s' && echo yes || echo no",
					path,
				),
			},
		},
	)
	if err != nil {
		return false, err
	}

	output = strings.TrimSpace(output)
	var exists = false
	if output == "yes" {
		exists = true
	} else if output == "no" {
		exists = false
	} else {
		return false, tracederrors.TracedErrorf(
			"Unexpected output when checking for existence: '%s'",
			output,
		)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return false, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "'%s' on host '%s' exists.", path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "'%s' on host '%s' does not exist.", path, hostDescription)
	}

	return exists, nil
}
