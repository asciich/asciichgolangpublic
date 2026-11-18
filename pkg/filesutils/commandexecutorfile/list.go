package commandexecutorfile

import (
	"context"
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ListSubDirectories(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, options *parameteroptions.ListDirectoryOptions) ([]string, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString(path)
	}

	if options == nil {
		options = &parameteroptions.ListDirectoryOptions{}
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List subdirectories of '%s' on '%s' started.", path, hostDescription)

	commandToUse := []string{"find", path, "-type", "d", "-maxdepth", "1"}
	if options.Recursive {
		commandToUse = []string{"find", path, "-type", "d"}
	}

	foundPaths, err := commandExecutor.RunCommandAndGetStdoutAsLines(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: commandToUse,
		},
	)
	if err != nil {
		return nil, err
	}

	slices.Sort(foundPaths)

	logging.LogInfoByCtxf(ctx, "List subdirectories of '%s' on '%s' finished. Found %d sub directories.", path, hostDescription, len(foundPaths))

	return foundPaths, nil
}
