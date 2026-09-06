package linuxutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// GetSysctlValue returns the current runtime value of the sysctl parameter with
// the given key on the system reachable through the given commandExecutor.
func GetSysctlValue(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, key string) (string, error) {
	if commandExecutor == nil {
		return "", tracederrors.TracedErrorNil("commandExecutor")
	}

	if key == "" {
		return "", tracederrors.TracedErrorEmptyString("key")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"sysctl", "-n", key},
		},
	)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout), nil
}

// SetSysctlValue sets the sysctl parameter with the given key to the given value
// at runtime on the system reachable through the given commandExecutor if it is
// not already set to the desired value.
func SetSysctlValue(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, key string, value string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if key == "" {
		return tracederrors.TracedErrorEmptyString("key")
	}

	if value == "" {
		return tracederrors.TracedErrorEmptyString("value")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set sysctl parameter '%s' to '%s' on '%s' started.", key, value, hostDescription)

	currentValue, err := GetSysctlValue(ctx, commandExecutor, key)
	if err != nil {
		return err
	}

	if currentValue == value {
		logging.LogInfoByCtxf(ctx, "Sysctl parameter '%s' is already set to '%s' on '%s'. Skip setting sysctl parameter.", key, value, hostDescription)
	} else {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"sysctl", "-w", fmt.Sprintf("%s=%s", key, value)},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Set sysctl parameter '%s' to '%s' on '%s'.", key, value, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Set sysctl parameter '%s' to '%s' on '%s' finished.", key, value, hostDescription)

	return nil
}

// SetSysctlValues sets all given sysctl parameters to their desired values at
// runtime on the system reachable through the given commandExecutor if they are
// not already set to the desired values.
func SetSysctlValues(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, values map[string]string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if values == nil {
		return tracederrors.TracedErrorNil("values")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set '%d' sysctl parameters on '%s' started.", len(values), hostDescription)

	for key, value := range values {
		err = SetSysctlValue(ctx, commandExecutor, key, value)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Set '%d' sysctl parameters on '%s' finished.", len(values), hostDescription)

	return nil
}
