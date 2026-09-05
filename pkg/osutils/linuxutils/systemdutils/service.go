package systemdutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsServiceActive returns true if the systemd service with the given name is
// currently active (running) on the system reachable through the given
// commandExecutor.
func IsServiceActive(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, serviceName string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if serviceName == "" {
		return false, tracederrors.TracedErrorEmptyString("serviceName")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				`systemctl is-active '` + serviceName + `' >/dev/null 2>&1 && echo yes || echo no`,
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
			"Unexpected output while checking if systemd service '%s' is active: '%s'",
			serviceName,
			trimmed,
		)
	}
}

// IsServiceEnabled returns true if the systemd service with the given name is
// enabled (started on boot) on the system reachable through the given
// commandExecutor.
func IsServiceEnabled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, serviceName string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if serviceName == "" {
		return false, tracederrors.TracedErrorEmptyString("serviceName")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				`systemctl is-enabled '` + serviceName + `' >/dev/null 2>&1 && echo yes || echo no`,
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
			"Unexpected output while checking if systemd service '%s' is enabled: '%s'",
			serviceName,
			trimmed,
		)
	}
}

// StartService starts the systemd service with the given name on the system
// reachable through the given commandExecutor if it is not already active.
func StartService(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, serviceName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if serviceName == "" {
		return tracederrors.TracedErrorEmptyString("serviceName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Start systemd service '%s' on '%s' started.", serviceName, hostDescription)

	active, err := IsServiceActive(ctx, commandExecutor, serviceName)
	if err != nil {
		return err
	}

	if active {
		logging.LogInfoByCtxf(ctx, "Systemd service '%s' is already active on '%s'. Skip starting service.", serviceName, hostDescription)
	} else {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"systemctl", "start", serviceName},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Started systemd service '%s' on '%s'.", serviceName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Start systemd service '%s' on '%s' finished.", serviceName, hostDescription)

	return nil
}

// StopService stops the systemd service with the given name on the system
// reachable through the given commandExecutor if it is currently active.
func StopService(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, serviceName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if serviceName == "" {
		return tracederrors.TracedErrorEmptyString("serviceName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Stop systemd service '%s' on '%s' started.", serviceName, hostDescription)

	active, err := IsServiceActive(ctx, commandExecutor, serviceName)
	if err != nil {
		return err
	}

	if active {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"systemctl", "stop", serviceName},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Stopped systemd service '%s' on '%s'.", serviceName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Systemd service '%s' is already stopped on '%s'. Skip stopping service.", serviceName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Stop systemd service '%s' on '%s' finished.", serviceName, hostDescription)

	return nil
}

// EnableService enables the systemd service with the given name on the system
// reachable through the given commandExecutor if it is not already enabled.
func EnableService(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, serviceName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if serviceName == "" {
		return tracederrors.TracedErrorEmptyString("serviceName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Enable systemd service '%s' on '%s' started.", serviceName, hostDescription)

	enabled, err := IsServiceEnabled(ctx, commandExecutor, serviceName)
	if err != nil {
		return err
	}

	if enabled {
		logging.LogInfoByCtxf(ctx, "Systemd service '%s' is already enabled on '%s'. Skip enabling service.", serviceName, hostDescription)
	} else {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"systemctl", "enable", serviceName},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Enabled systemd service '%s' on '%s'.", serviceName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Enable systemd service '%s' on '%s' finished.", serviceName, hostDescription)

	return nil
}

// DisableService disables the systemd service with the given name on the system
// reachable through the given commandExecutor if it is currently enabled.
func DisableService(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, serviceName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if serviceName == "" {
		return tracederrors.TracedErrorEmptyString("serviceName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Disable systemd service '%s' on '%s' started.", serviceName, hostDescription)

	enabled, err := IsServiceEnabled(ctx, commandExecutor, serviceName)
	if err != nil {
		return err
	}

	if enabled {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"systemctl", "disable", serviceName},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Disabled systemd service '%s' on '%s'.", serviceName, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Systemd service '%s' is already disabled on '%s'. Skip disabling service.", serviceName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Disable systemd service '%s' on '%s' finished.", serviceName, hostDescription)

	return nil
}

// EnableAndStartService enables and starts the systemd service with the given
// name on the system reachable through the given commandExecutor. The service is
// only touched when it is not already both enabled and active.
func EnableAndStartService(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, serviceName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if serviceName == "" {
		return tracederrors.TracedErrorEmptyString("serviceName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Enable and start systemd service '%s' on '%s' started.", serviceName, hostDescription)

	enabled, err := IsServiceEnabled(ctx, commandExecutor, serviceName)
	if err != nil {
		return err
	}

	active, err := IsServiceActive(ctx, commandExecutor, serviceName)
	if err != nil {
		return err
	}

	if enabled && active {
		logging.LogInfoByCtxf(ctx, "Systemd service '%s' is already enabled and active on '%s'. Skip enabling and starting service.", serviceName, hostDescription)
	} else {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"systemctl", "enable", "--now", serviceName},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Enabled and started systemd service '%s' on '%s'.", serviceName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Enable and start systemd service '%s' on '%s' finished.", serviceName, hostDescription)

	return nil
}

// RestartService restarts the systemd service with the given name on the system
// reachable through the given commandExecutor.
func RestartService(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, serviceName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if serviceName == "" {
		return tracederrors.TracedErrorEmptyString("serviceName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Restart systemd service '%s' on '%s' started.", serviceName, hostDescription)

	_, err = commandExecutor.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"systemctl", "restart", serviceName},
		},
	)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Restarted systemd service '%s' on '%s'.", serviceName, hostDescription)

	logging.LogInfoByCtxf(ctx, "Restart systemd service '%s' on '%s' finished.", serviceName, hostDescription)

	return nil
}
