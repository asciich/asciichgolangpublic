package systemdutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// DaemonReload reloads the systemd manager configuration on the system
// reachable through the given commandExecutor (equivalent to
// `systemctl daemon-reload`). This makes systemd pick up newly added or
// changed unit files.
func DaemonReload(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Reload systemd manager configuration on '%s' started.", hostDescription)

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"systemctl", "daemon-reload"},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Reloaded systemd manager configuration on '%s'.", hostDescription)

	logging.LogInfoByCtxf(ctx, "Reload systemd manager configuration on '%s' finished.", hostDescription)

	return nil
}
