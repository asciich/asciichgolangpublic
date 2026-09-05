package linuxutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsSwapEnabled returns true if there is currently at least one active swap
// device/file on the system reachable through the given commandExecutor.
//
// Following the constitution we do NOT rely on the exit code of a single
// command. Instead we let the shell decide and echo a well known value so we
// can be sure the command was actually executed.
func IsSwapEnabled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				`test -n "$(swapon --show --noheadings 2>/dev/null)" && echo yes || echo no`,
			},
		},
	)
	if err != nil {
		return false, err
	}

	trimmed := strings.TrimSpace(stdout)

	switch trimmed {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, tracederrors.TracedErrorf(
			"Unexpected output while checking if swap is enabled: '%s'",
			trimmed,
		)
	}
}

// TurnSwapOff turns off swap on the system reachable through the given
// commandExecutor if it is needed.
//
// This function is implemented in an idempotent way: swap is only turned off
// when it is currently enabled. If swap is already off nothing is changed.
func TurnSwapOff(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Turn swap off on '%s' started.", hostDescription)

	enabled, err := IsSwapEnabled(ctx, commandExecutor)
	if err != nil {
		return err
	}

	if enabled {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"swapoff", "-a"},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Turned swap off on '%s'.", hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Swap is already turned off on '%s'. Skip turning swap off.", hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Turn swap off on '%s' finished.", hostDescription)

	return nil
}
