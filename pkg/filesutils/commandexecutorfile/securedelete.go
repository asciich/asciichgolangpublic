package commandexecutorfile

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// SecurelyDelete securely deletes a file by overwriting it before deletion using shred.
func SecurelyDelete(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, path string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	// Check if file exists
	exists, err := Exists(ctx, commandExecutor, path)
	if err != nil {
		return err
	}

	if !exists {
		logging.LogInfoByCtxf(ctx, "File '%s' on '%s' already absent. Skip secure delete.", path, hostDescription)
		return nil
	}

	// Use shred to securely delete the file
	command := []string{"shred", "-u", "-n", "1", "-z", path}

	logging.LogInfoByCtxf(ctx, "Securely deleting '%s' on '%s' using shred.", path, hostDescription)

	_, err = commandExecutor.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: command,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to securely delete '%s' on '%s': %w", path, hostDescription, err)
	}

	logging.LogChangedByCtxf(ctx, "Securely deleted '%s' on '%s'.", path, hostDescription)

	return nil
}
