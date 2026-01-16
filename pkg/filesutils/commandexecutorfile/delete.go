package commandexecutorfile

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/pathsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Delete(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string, options *filesoptions.DeleteOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	if options == nil {
		options = &filesoptions.DeleteOptions{}
	}

	if !pathsutils.IsAbsolutePath(filePath) {
		return tracederrors.TracedErrorf(
			"For security reasons deleting a is only implemented for absolute paths but got '%s'",
			filePath,
		)
	}

	exists, err := Exists(ctx, commandExecutor, filePath)
	if err != nil {
		return err
	}

	command := []string{"rm", filePath}

	if options != nil && options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	if exists {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf("'%s' on '%s' deleted.", filePath, hostDescription)
	} else {
		logging.LogInfof("'%s' on '%s' already absent.", filePath, hostDescription)
	}

	return nil
}

func DeleteDirectory(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string, options *filesoptions.DeleteOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if filePath == "" {
		return tracederrors.TracedErrorEmptyString("filePath")
	}

	if options == nil {
		options = &filesoptions.DeleteOptions{}
	}

	if !pathsutils.IsAbsolutePath(filePath) {
		return tracederrors.TracedErrorf(
			"For security reasons deleting a is only implemented for absolute paths but got '%s'",
			filePath,
		)
	}

	exists, err := Exists(ctx, commandExecutor, filePath)
	if err != nil {
		return err
	}

	command := []string{"rm", "-rf", filePath}

	if options != nil && options.UseSudo {
		command = append([]string{"sudo"}, command...)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	if exists {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedf("Directory '%s' on '%s' deleted.", filePath, hostDescription)
	} else {
		logging.LogInfof("Directory '%s' on '%s' already absent.", filePath, hostDescription)
	}

	return nil
}
