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

// IsKernelModuleLoaded returns true if the kernel module with the given name is
// currently loaded on the system reachable through the given commandExecutor.
func IsKernelModuleLoaded(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, moduleName string) (bool, error) {
	if commandExecutor == nil {
		return false, tracederrors.TracedErrorNil("commandExecutor")
	}

	if moduleName == "" {
		return false, tracederrors.TracedErrorEmptyString("moduleName")
	}

	stdout, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"sh", "-c",
				fmt.Sprintf(`lsmod | awk '{print $1}' | grep -qx '%s' && echo yes || echo no`, moduleName),
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
			"Unexpected output while checking if kernel module '%s' is loaded: '%s'",
			moduleName,
			trimmed,
		)
	}
}

// LoadKernelModule loads the kernel module with the given name on the system
// reachable through the given commandExecutor if it is not already loaded.
func LoadKernelModule(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, moduleName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if moduleName == "" {
		return tracederrors.TracedErrorEmptyString("moduleName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Load kernel module '%s' on '%s' started.", moduleName, hostDescription)

	loaded, err := IsKernelModuleLoaded(ctx, commandExecutor, moduleName)
	if err != nil {
		return err
	}

	if loaded {
		logging.LogInfoByCtxf(ctx, "Kernel module '%s' is already loaded on '%s'. Skip loading kernel module.", moduleName, hostDescription)
	} else {
		_, err = commandExecutor.RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: []string{"modprobe", moduleName},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Loaded kernel module '%s' on '%s'.", moduleName, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Load kernel module '%s' on '%s' finished.", moduleName, hostDescription)

	return nil
}

// LoadKernelModules loads all given kernel modules on the system reachable
// through the given commandExecutor if they are not already loaded.
func LoadKernelModules(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, moduleNames []string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if moduleNames == nil {
		return tracederrors.TracedErrorNil("moduleNames")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Load '%d' kernel modules on '%s' started.", len(moduleNames), hostDescription)

	for _, moduleName := range moduleNames {
		err = LoadKernelModule(ctx, commandExecutor, moduleName)
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Load '%d' kernel modules on '%s' finished.", len(moduleNames), hostDescription)

	return nil
}
