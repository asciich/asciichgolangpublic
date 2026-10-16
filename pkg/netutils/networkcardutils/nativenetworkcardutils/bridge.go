package nativenetworkcardutils

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/godbus/dbus/v5"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/networkcardutils/genericnetworkcardutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
)

const applyTimeout = 30 * time.Second
const applyPollInterval = 500 * time.Millisecond

// CreateLayer2Bridge creates the layer 2 Linux bridge 'bridgeName' and attaches 'interfaceName' as bridge port.
//
// systemd-networkd is used for the configuration (see genericnetworkcardutils.GetLayer2BridgeSystemdNetworkdFiles):
//   - The bridge takes over the MAC address of 'interfaceName' so DHCP reservations keep working.
//   - The bridge gets its IP configuration using DHCP.
//
// The configuration is applied instantly (systemd-networkd reload) and is persistent across reboots.
//
// Requires root privileges and an active and enabled systemd-networkd.
// This function is idempotent.
func CreateLayer2Bridge(ctx context.Context, bridgeName string, interfaceName string) error {
	err := genericnetworkcardutils.CheckLayer2BridgeParameters(bridgeName, interfaceName)
	if err != nil {
		return err
	}

	isRoot, err := userutils.IsRunningAsRoot(ctx)
	if err != nil {
		return err
	}

	if !isRoot {
		return tracederrors.TracedErrorf("Create layer 2 bridge '%s' with interface '%s' requires root privileges.", bridgeName, interfaceName)
	}

	logging.LogInfoByCtxf(ctx, "Create layer 2 bridge '%s' with interface '%s' started.", bridgeName, interfaceName)

	conn, err := dbus.ConnectSystemBus(dbus.WithContext(ctx))
	if err != nil {
		return tracederrors.TracedErrorf("Unable to connect to the system D-Bus: %w", err)
	}
	defer conn.Close()

	err = checkSystemdNetworkdActiveAndEnabled(ctx, conn)
	if err != nil {
		return err
	}

	err = checkInterfaceBridgeable(ctx, bridgeName, interfaceName)
	if err != nil {
		return err
	}

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to get network interface '%s': %w", interfaceName, err)
	}

	files, err := genericnetworkcardutils.GetLayer2BridgeSystemdNetworkdFiles(bridgeName, interfaceName, iface.HardwareAddr.String())
	if err != nil {
		return err
	}

	err = nativefiles.CreateDirectory(ctx, genericnetworkcardutils.SystemdNetworkdConfigDirectory, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	filesChanged := false
	for _, file := range files {
		changed, err := ensureFileContent(ctx, file.Path, file.Content)
		if err != nil {
			return err
		}

		filesChanged = filesChanged || changed
	}

	isApplied, err := isLayer2BridgeApplied(ctx, bridgeName, interfaceName)
	if err != nil {
		return err
	}

	if !filesChanged && isApplied {
		logging.LogInfoByCtxf(ctx, "Layer 2 bridge '%s' with interface '%s' is already configured and active.", bridgeName, interfaceName)
	} else {
		err = applySystemdNetworkdConfig(ctx, conn, iface.Index)
		if err != nil {
			return err
		}

		err = genericnetworkcardutils.WaitUntilTrue(
			ctx,
			applyTimeout,
			applyPollInterval,
			"interface '"+interfaceName+"' attached to bridge '"+bridgeName+"'",
			func(ctx context.Context) (bool, error) {
				return isLayer2BridgeApplied(ctx, bridgeName, interfaceName)
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

func checkSystemdNetworkdActiveAndEnabled(ctx context.Context, conn *dbus.Conn) error {
	if conn == nil {
		return tracederrors.TracedErrorNil("conn")
	}

	const unitName = "systemd-networkd.service"

	var unitPath dbus.ObjectPath
	err := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1").
		CallWithContext(ctx, "org.freedesktop.systemd1.Manager.LoadUnit", 0, unitName).
		Store(&unitPath)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to load unit '%s' over D-Bus: %w", unitName, err)
	}

	unit := conn.Object("org.freedesktop.systemd1", unitPath)

	activeState, err := getStringProperty(unit, "org.freedesktop.systemd1.Unit.ActiveState")
	if err != nil {
		return err
	}

	if activeState != "active" {
		return tracederrors.TracedErrorf("'%s' must be active to create a layer 2 bridge, but its state is '%s'.", unitName, activeState)
	}

	unitFileState, err := getStringProperty(unit, "org.freedesktop.systemd1.Unit.UnitFileState")
	if err != nil {
		return err
	}

	if unitFileState != "enabled" {
		return tracederrors.TracedErrorf("'%s' must be enabled for a persistent layer 2 bridge, but its unit file state is '%s'.", unitName, unitFileState)
	}

	return nil
}

func getStringProperty(object dbus.BusObject, propertyName string) (string, error) {
	variant, err := object.GetProperty(propertyName)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to get D-Bus property '%s': %w", propertyName, err)
	}

	ret, ok := variant.Value().(string)
	if !ok {
		return "", tracederrors.TracedErrorf("D-Bus property '%s' is not a string: '%v'", propertyName, variant.Value())
	}

	return ret, nil
}

func applySystemdNetworkdConfig(ctx context.Context, conn *dbus.Conn, interfaceIndex int) error {
	networkd := conn.Object("org.freedesktop.network1", "/org/freedesktop/network1")

	// Creates missing netdevs (the bridge) and reconfigures interfaces with new or changed .network files.
	err := networkd.CallWithContext(ctx, "org.freedesktop.network1.Manager.Reload", 0).Err
	if err != nil {
		return tracederrors.TracedErrorf("Reload systemd-networkd failed: %w", err)
	}

	// Ensures the bridge port is reconfigured even if only the runtime state was missing.
	err = networkd.CallWithContext(ctx, "org.freedesktop.network1.Manager.ReconfigureLink", 0, int32(interfaceIndex)).Err
	if err != nil {
		return tracederrors.TracedErrorf("Reconfigure link with index '%d' using systemd-networkd failed: %w", interfaceIndex, err)
	}

	return nil
}
func checkInterfaceBridgeable(ctx context.Context, bridgeName string, interfaceName string) error {
	interfacePath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName)
	if err != nil {
		return err
	}

	interfaceExists, err := nativefiles.Exists(ctx, interfacePath)
	if err != nil {
		return err
	}

	if !interfaceExists {
		return tracederrors.TracedErrorf("Network interface '%s' does not exist.", interfaceName)
	}

	wirelessPath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName, "wireless")
	if err != nil {
		return err
	}

	if nativefiles.IsDir(ctx, wirelessPath) {
		return tracederrors.TracedErrorf("Network interface '%s' is a wireless interface. Wireless client interfaces can not be used as layer 2 bridge port.", interfaceName)
	}

	bridgePath, err := genericnetworkcardutils.GetSysClassNetPath(bridgeName)
	if err != nil {
		return err
	}

	bridgeExists, err := nativefiles.Exists(ctx, bridgePath)
	if err != nil {
		return err
	}

	if bridgeExists {
		var existingIsBridge bool
		existingIsBridge, err = isBridge(ctx, bridgeName)
		if err != nil {
			return err
		}

		if !existingIsBridge {
			return tracederrors.TracedErrorf("Interface '%s' already exists but is not a bridge.", bridgeName)
		}
	}

	master, err := getMasterName(interfaceName)
	if err != nil {
		return err
	}

	if master != "" && master != bridgeName {
		return tracederrors.TracedErrorf("Network interface '%s' is already attached to '%s' and can not be attached to bridge '%s'.", interfaceName, master, bridgeName)
	}

	return nil
}

func isBridge(ctx context.Context, interfaceName string) (bool, error) {
	bridgeDirPath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName, "bridge")
	if err != nil {
		return false, err
	}

	return nativefiles.IsDir(ctx, bridgeDirPath), nil
}

// getMasterName returns the name of the master (e.g. bridge) of 'interfaceName' or an empty string if there is no master.
func getMasterName(interfaceName string) (string, error) {
	masterLinkPath, err := genericnetworkcardutils.GetSysClassNetPath(interfaceName, "master")
	if err != nil {
		return "", err
	}

	target, err := os.Readlink(masterLinkPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No master link means the interface is not attached to a bridge. This is a valid state, not an error.
			return "", nil
		}

		return "", tracederrors.TracedErrorf("Unable to read master link '%s': %w", masterLinkPath, err)
	}

	return filepath.Base(target), nil
}

func isLayer2BridgeApplied(ctx context.Context, bridgeName string, interfaceName string) (bool, error) {
	bridgeExists, err := isBridge(ctx, bridgeName)
	if err != nil {
		return false, err
	}

	if !bridgeExists {
		return false, nil
	}

	master, err := getMasterName(interfaceName)
	if err != nil {
		return false, err
	}

	return master == bridgeName, nil
}

func ensureFileContent(ctx context.Context, filePath string, content string) (bool, error) {
	if nativefiles.IsFile(ctx, filePath) {
		currentContent, err := nativefiles.ReadAsString(ctx, filePath, &filesoptions.ReadOptions{})
		if err != nil {
			return false, err
		}

		if currentContent == content {
			logging.LogInfoByCtxf(ctx, "File '%s' already has the desired content.", filePath)
			return false, nil
		}
	}

	err := nativefiles.WriteString(ctx, filePath, content)
	if err != nil {
		return false, err
	}

	logging.LogChangedByCtxf(ctx, "Updated systemd-networkd config file '%s'.", filePath)

	return true, nil
}
