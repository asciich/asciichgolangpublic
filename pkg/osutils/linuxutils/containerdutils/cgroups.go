package containerdutils

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils/systemdutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const configPath = "/etc/containerd/config.toml"

// IsCGroupEnabled returns true if the systemd cgroup driver is currently enabled
// in the containerd config file on the system reachable through the given
// commandExecutor.
//
// Following the constitution we do NOT rely on the exit code of a single
// command. Instead we let the shell decide and echo a well known value so we
// can be sure the command was actually executed.
func IsCGroupEnabled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				`grep -q 'SystemdCgroup = true' '` + configPath + `' 2>/dev/null && echo yes || echo no`,
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
			"Unexpected output while checking if containerd systemd cgroup driver is enabled: '%s'",
			trimmed,
		)
	}
}

// EnableCGroup enables the systemd cgroup driver in the containerd config on the
// system reachable through the given commandExecutor and restarts containerd so
// the change takes effect.
//
// This function is implemented in an idempotent way: the config is only
// regenerated and containerd is only restarted when the systemd cgroup driver is
// currently not enabled. If it is already enabled nothing is changed.
func EnableCGroup(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Enable containerd systemd cgroup driver on '%s' started.", hostDescription)

	enabled, err := IsCGroupEnabled(ctx, commandExecutor)
	if err != nil {
		return err
	}

	if enabled {
		logging.LogInfoByCtxf(ctx, "Containerd systemd cgroup driver is already enabled on '%s'. Skip enabling.", hostDescription)
	} else {
		// Generate the default containerd config and read it from stdout.
		defaultConfig, err := commandExecutor.RunCommandAndGetStdoutAsString(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"containerd", "config", "default"},
			},
		)
		if err != nil {
			return err
		}

		// Enable the systemd cgroup driver by replacing the setting in Go.
		config := strings.ReplaceAll(defaultConfig, "SystemdCgroup = false", "SystemdCgroup = true")

		err = commandexecutorfile.CreateDirectory(ctx, commandExecutor, filepath.Dir(configPath), &filesoptions.CreateOptions{})
		if err != nil {
			return err
		}

		err = commandexecutorfile.WriteBytes(ctx, commandExecutor, configPath, []byte(config), &filesoptions.WriteOptions{})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Enabled containerd systemd cgroup driver on '%s'.", hostDescription)

		err = systemdutils.RestartService(ctx, commandExecutor, "containerd")
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Enable containerd systemd cgroup driver on '%s' finished.", hostDescription)

	return nil
}
