package commandexecutorfile

import (
	"context"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CreateFile(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, options *filesoptions.CreateOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if options == nil {
		options = &filesoptions.CreateOptions{}
	}

	exists, err := FileExists(ctx, commandExecutor, path)
	if err != nil {
		return err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "File '%s' on '%s' already exists. Skip file creation.", path, hostDescription)
	} else {
		// Ensure the parent directory exists before creating the file.
		parentDir := filepath.Dir(path)
		err = CreateDirectory(ctx, commandExecutor, parentDir, &filesoptions.CreateOptions{UseSudo: options.UseSudo})
		if err != nil {
			return err
		}

		command := []string{"touch", path}
		if options.UseSudo {
			command = append([]string{"sudo"}, command...)
		}

		_, err := commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Created file '%s' on '%s'.", path, hostDescription)
	}

	return nil
}

func CreateDirectory(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string, options *filesoptions.CreateOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	if options == nil {
		options = &filesoptions.CreateOptions{}
	}

	exists, err := DirectoryExists(ctx, commandExecutor, path)
	if err != nil {
		return err
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Directory '%s' on '%s' already exists. Skip directory creation.", path, hostDescription)
	} else {
		command := []string{"mkdir", "-p", path}
		if options.UseSudo {
			command = append([]string{"sudo"}, command...)
		}

		_, err := commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Created directory '%s' on '%s'.", path, hostDescription)
	}

	return nil
}
