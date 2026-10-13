package genericnetworkcardutils

import (
	"context"
	"fmt"
	"net"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Directory where systemd-networkd reads the administrator managed configuration files.
const SystemdNetworkdConfigDirectory = "/etc/systemd/network"

// Directory where the kernel exposes the network interfaces.
const SysClassNetDirectory = "/sys/class/net"

const managedByHeader = "# This file is managed by asciichgolangpublic networkcardutils.CreateLayer2Bridge.\n# Manual changes will be overwritten.\n"

// Linux limits interface names to 15 characters (IFNAMSIZ - 1).
// The allowed characters are intentionally restricted so names are safe to use in shell commands.
var interfaceNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,14}$`)

// SystemdNetworkdFile is a systemd-networkd configuration file to deploy.
type SystemdNetworkdFile struct {
	Path    string
	Content string
}

// CheckInterfaceName returns an error if 'interfaceName' is not a valid (and shell safe) Linux interface name.
func CheckInterfaceName(interfaceName string) error {
	if interfaceName == "" {
		return tracederrors.TracedErrorEmptyString("interfaceName")
	}

	if !interfaceNameRegex.MatchString(interfaceName) {
		return tracederrors.TracedErrorf(
			"Invalid interface name '%s'. Up to 15 characters out of [a-zA-Z0-9_.-] starting with an alphanumeric character are supported.",
			interfaceName,
		)
	}

	return nil
}

// CheckLayer2BridgeParameters validates the parameters used to create a layer 2 bridge.
func CheckLayer2BridgeParameters(bridgeName string, interfaceName string) error {
	if bridgeName == "" {
		return tracederrors.TracedErrorEmptyString("bridgeName")
	}

	if interfaceName == "" {
		return tracederrors.TracedErrorEmptyString("interfaceName")
	}

	err := CheckInterfaceName(bridgeName)
	if err != nil {
		return err
	}

	err = CheckInterfaceName(interfaceName)
	if err != nil {
		return err
	}

	if bridgeName == interfaceName {
		return tracederrors.TracedErrorf("bridgeName and interfaceName must differ but both are '%s'.", bridgeName)
	}

	return nil
}

// GetSysClassNetPath returns the sysfs path of 'interfaceName', optionally extended by 'subPaths'.
func GetSysClassNetPath(interfaceName string, subPaths ...string) (string, error) {
	err := CheckInterfaceName(interfaceName)
	if err != nil {
		return "", err
	}

	ret := path.Join(append([]string{SysClassNetDirectory, interfaceName}, subPaths...)...)

	return ret, nil
}

// NormalizeEthernetMacAddress validates 'macAddress' as 48 bit Ethernet MAC and returns it in lower case "aa:bb:cc:dd:ee:ff" notation.
func NormalizeEthernetMacAddress(macAddress string) (string, error) {
	macAddress = strings.TrimSpace(macAddress)
	if macAddress == "" {
		return "", tracederrors.TracedErrorEmptyString("macAddress")
	}

	hardwareAddr, err := net.ParseMAC(macAddress)
	if err != nil {
		return "", tracederrors.TracedErrorf("Invalid MAC address '%s': %w", macAddress, err)
	}

	if len(hardwareAddr) != 6 {
		return "", tracederrors.TracedErrorf("MAC address '%s' is not a 48 bit Ethernet MAC address.", macAddress)
	}

	ret := hardwareAddr.String()
	if ret == "00:00:00:00:00:00" {
		return "", tracederrors.TracedErrorf("MAC address '%s' is not usable for a bridge.", macAddress)
	}

	return ret, nil
}

// GetLayer2BridgeSystemdNetworkdFiles returns the systemd-networkd files needed to create the bridge 'bridgeName'
// with 'interfaceName' as bridge port. The bridge uses 'macAddress' (usually the MAC of 'interfaceName')
// so existing DHCP reservations keep working.
//
// The bind file uses the prefix '05-' to take precedence over typical existing configurations like '20-wired-<interface>.network'.
func GetLayer2BridgeSystemdNetworkdFiles(bridgeName string, interfaceName string, macAddress string) ([]SystemdNetworkdFile, error) {
	err := CheckLayer2BridgeParameters(bridgeName, interfaceName)
	if err != nil {
		return nil, err
	}

	normalizedMac, err := NormalizeEthernetMacAddress(macAddress)
	if err != nil {
		return nil, err
	}

	ret := []SystemdNetworkdFile{
		{
			Path: path.Join(SystemdNetworkdConfigDirectory, bridgeName+".netdev"),
			Content: managedByHeader + fmt.Sprintf(
				"\n[NetDev]\nName=%s\nKind=bridge\nMACAddress=%s\n",
				bridgeName, normalizedMac,
			),
		},
		{
			Path: path.Join(SystemdNetworkdConfigDirectory, fmt.Sprintf("05-%s-bind-%s.network", bridgeName, interfaceName)),
			Content: managedByHeader + fmt.Sprintf(
				"\n[Match]\nName=%s\n\n[Network]\nBridge=%s\n",
				interfaceName, bridgeName,
			),
		},
		{
			Path: path.Join(SystemdNetworkdConfigDirectory, fmt.Sprintf("10-%s.network", bridgeName)),
			Content: managedByHeader + fmt.Sprintf(
				"\n[Match]\nName=%s\n\n[Network]\nDHCP=yes\n",
				bridgeName,
			),
		},
	}

	return ret, nil
}

// WaitUntilTrue calls 'check' every 'interval' until it returns true, it returns an error or 'timeout' is reached.
func WaitUntilTrue(ctx context.Context, timeout time.Duration, interval time.Duration, description string, check func(ctx context.Context) (bool, error)) error {
	if check == nil {
		return tracederrors.TracedErrorNil("check")
	}

	if description == "" {
		return tracederrors.TracedErrorEmptyString("description")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		isTrue, err := check(timeoutCtx)
		if err != nil {
			return err
		}

		if isTrue {
			return nil
		}

		select {
		case <-timeoutCtx.Done():
			return tracederrors.TracedErrorf("Timeout after '%s' while waiting for: %s: %w", timeout, description, timeoutCtx.Err())
		case <-ticker.C:
		}
	}
}
