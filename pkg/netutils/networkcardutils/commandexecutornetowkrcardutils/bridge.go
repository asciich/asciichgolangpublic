package commandexecutornetworkcardutils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils/genericnetworkcardutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/userutils/commandexecutoruserutils"
)

const applyTimeout = 30 * time.Second
const applyPollInterval = 500 * time.Millisecond

// CreateLayer2Bridge creates the layer 2 Linux bridge 'bridgeName' and attaches 'interfaceName' as bridge port
// on the machine reached by 'commandExecutor'.
//
// Behaves like nativenetworkcardutils.CreateLayer2Bridge:
// systemd-networkd based, instantly applied, persistent across reboots and idempotent.
//
// Requires root privileges on the target machine and an active and enabled systemd-networkd.
//
// Warning: If 'commandExecutor' connects (e.g. via SSH) over 'interfaceName' the connection is shortly interrupted
// while the interface is moved into the bridge.
func CreateLayer2Bridge(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, bridgeName string, interfaceName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	err := genericnetworkcardutils.CheckLayer2BridgeParameters(bridgeName, interfaceName)
	if err != nil {
		return err
	}

	isRoot, err := commandexecutoruserutils.IsRunningAsRoot(ctx, commandExecutor)
	if err != nil {
		return err
	}

	if !isRoot {
		return tracederrors.TracedErrorf("Create layer 2 bridge '%s' with interface '%s' requires root privileges.", bridgeName, interfaceName)
	}

	logging.LogInfoByCtxf(ctx, "Create layer 2 bridge '%s' with interface '%s' started.", bridgeName, interfaceName)

	err = checkSystemdNetworkdActiveAndEnabled(ctx, commandExecutor)
	if err != nil {
		return err
	}

	err = checkInterfaceBridgeable(ctx, commandExecutor, bridgeName, interfaceName)
	if err != nil {
		return err
	}

	macAddressPath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName, "address")
	if err != nil {
		return err
	}

	macAddress, err := commandexecutorfile.ReadAsString(ctx, commandExecutor, macAddressPath)
	if err != nil {
		return err
	}

	files, err := genericnetworkcardutils.GetLayer2BridgeSystemdNetworkdFiles(bridgeName, interfaceName, macAddress)
	if err != nil {
		return err
	}

	err = commandexecutorfile.CreateDirectory(ctx, commandExecutor, genericnetworkcardutils.SystemdNetworkdConfigDirectory, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	filesChanged := false
	for _, file := range files {
		changed, err := ensureFileContent(ctx, commandExecutor, file.Path, file.Content)
		if err != nil {
			return err
		}

		filesChanged = filesChanged || changed
	}

	isApplied, err := isLayer2BridgeApplied(ctx, commandExecutor, bridgeName, interfaceName)
	if err != nil {
		return err
	}

	if !filesChanged && isApplied {
		logging.LogInfoByCtxf(ctx, "Layer 2 bridge '%s' with interface '%s' is already configured and active.", bridgeName, interfaceName)
	} else {
		// Creates missing netdevs (the bridge) and reconfigures interfaces with new or changed .network files.
		_, err = runShellCommand(ctx, commandExecutor, "networkctl reload")
		if err != nil {
			return err
		}

		// Ensures the bridge port is reconfigured even if only the runtime state was missing.
		_, err = runShellCommand(ctx, commandExecutor, fmt.Sprintf("networkctl reconfigure '%s'", interfaceName))
		if err != nil {
			return err
		}

		err = genericnetworkcardutils.WaitUntilTrue(
			ctx,
			applyTimeout,
			applyPollInterval,
			"interface '"+interfaceName+"' attached to bridge '"+bridgeName+"'",
			func(ctx context.Context) (bool, error) {
				return isLayer2BridgeApplied(ctx, commandExecutor, bridgeName, interfaceName)
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Layer 2 bridge '%s' with interface '%s' applied to the running system.", bridgeName, interfaceName)
	}

	logging.LogInfoByCtxf(ctx, "Create layer 2 bridge '%s' with interface '%s' finished.", bridgeName, interfaceName)

	return nil
}

func runShellCommand(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, shellCommand string) (string, error) {
	return commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"sh", "-c", shellCommand},
		},
	)
}

// runYesNoCommand evaluates 'condition' without relying on the exit code of the executed command itself.
func runYesNoCommand(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, condition string) (bool, error) {
	stdout, err := runShellCommand(ctx, commandExecutor, condition+" && echo yes || echo no")
	if err != nil {
		return false, err
	}

	switch strings.TrimSpace(stdout) {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, tracederrors.TracedErrorf("Unexpected output '%s' when evaluating '%s'.", stdout, condition)
	}
}

func checkSystemdNetworkdActiveAndEnabled(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	const unitName = "systemd-networkd.service"

	isActive, err := runYesNoCommand(ctx, commandExecutor, "systemctl is-active --quiet "+unitName)
	if err != nil {
		return err
	}

	if !isActive {
		return tracederrors.TracedErrorf("'%s' must be active to create a layer 2 bridge.", unitName)
	}

	isEnabled, err := runYesNoCommand(ctx, commandExecutor, "systemctl is-enabled --quiet "+unitName)
	if err != nil {
		return err
	}

	if !isEnabled {
		return tracederrors.TracedErrorf("'%s' must be enabled for a persistent layer 2 bridge.", unitName)
	}

	return nil
}

func checkInterfaceBridgeable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, bridgeName string, interfaceName string) error {
	interfacePath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName)
	if err != nil {
		return err
	}

	exists, err := commandexecutorfile.Exists(ctx, commandExecutor, interfacePath)
	if err != nil {
		return err
	}

	if !exists {
		return tracederrors.TracedErrorf("Network interface '%s' does not exist.", interfaceName)
	}

	wirelessPath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName, "wireless")
	if err != nil {
		return err
	}

	isWireless, err := commandexecutorfile.DirectoryExists(ctx, commandExecutor, wirelessPath)
	if err != nil {
		return err
	}

	if isWireless {
		return tracederrors.TracedErrorf("Network interface '%s' is a wireless interface. Wireless client interfaces can not be used as layer 2 bridge port.", interfaceName)
	}

	bridgePath, err := genericnetworkcardutils.GetSysClassNetPath(bridgeName)
	if err != nil {
		return err
	}

	bridgeExists, err := commandexecutorfile.Exists(ctx, commandExecutor, bridgePath)
	if err != nil {
		return err
	}

	if bridgeExists {
		isBridgeInterface, err := isBridge(ctx, commandExecutor, bridgeName)
		if err != nil {
			return err
		}

		if !isBridgeInterface {
			return tracederrors.TracedErrorf("Interface '%s' already exists but is not a bridge.", bridgeName)
		}
	}

	master, err := getMasterName(ctx, commandExecutor, interfaceName)
	if err != nil {
		return err
	}

	if master != "" && master != bridgeName {
		return tracederrors.TracedErrorf("Network interface '%s' is already attached to '%s' and can not be attached to bridge '%s'.", interfaceName, master, bridgeName)
	}

	return nil
}

func isBridge(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, interfaceName string) (bool, error) {
	bridgeDirPath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName, "bridge")
	if err != nil {
		return false, err
	}

	return commandexecutorfile.DirectoryExists(ctx, commandExecutor, bridgeDirPath)
}

// getMasterName returns the name of the master (e.g. bridge) of 'interfaceName' or an empty string if there is no master.
func getMasterName(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, interfaceName string) (string, error) {
	masterLinkPath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName, "master")
	if err != nil {
		return "", err
	}

	// The "master=" prefix ensures the command was really executed, even if there is no master.
	stdout, err := runShellCommand(
		ctx,
		commandExecutor,
		fmt.Sprintf(
			`if [ -L '%s' ]; then printf 'master=%%s\n' "$(basename "$(readlink '%s')")"; else printf 'master=\n'; fi`,
			masterLinkPath, masterLinkPath,
		),
	)
	if err != nil {
		return "", err
	}

	line := strings.TrimSpace(stdout)
	if !strings.HasPrefix(line, "master=") {
		return "", tracederrors.TracedErrorf("Unexpected output '%s' when reading master of interface '%s'.", stdout, interfaceName)
	}

	return strings.TrimPrefix(line, "master="), nil
}

func isLayer2BridgeApplied(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, bridgeName string, interfaceName string) (bool, error) {
	bridgeExists, err := isBridge(ctx, commandExecutor, bridgeName)
	if err != nil {
		return false, err
	}

	if !bridgeExists {
		return false, nil
	}

	master, err := getMasterName(ctx, commandExecutor, interfaceName)
	if err != nil {
		return false, err
	}

	return master == bridgeName, nil
}

func ensureFileContent(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, filePath string, content string) (bool, error) {
	exists, err := commandexecutorfile.FileExists(ctx, commandExecutor, filePath)
	if err != nil {
		return false, err
	}

	if exists {
		currentContent, err := commandexecutorfile.ReadAsString(ctx, commandExecutor, filePath)
		if err != nil {
			return false, err
		}

		if currentContent == content {
			logging.LogInfoByCtxf(ctx, "File '%s' already has the desired content.", filePath)
			return false, nil
		}
	}

	err = commandexecutorfile.WriteString(ctx, commandExecutor, filePath, content, &filesoptions.WriteOptions{})
	if err != nil {
		return false, err
	}

	logging.LogChangedByCtxf(ctx, "Updated systemd-networkd config file '%s'.", filePath)

	return true, nil
}
