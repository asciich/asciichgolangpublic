package networkcardutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
)

// Constants taken from the Linux kernel headers:
//   - include/uapi/linux/sockios.h
//   - include/uapi/linux/ethtool.h
//   - include/uapi/linux/if.h
const (
	siocEthtool = 0x8946
	ifNameSize  = 16

	ethtoolGTSO      = 0x0000001e // Get TCP segmentation offload
	ethtoolSTSO      = 0x0000001f // Set TCP segmentation offload (covers TSO, TSO6, ECN, mangleid)
	ethtoolGPERMADDR = 0x00000020 // Get permanent hardware address
	ethtoolGGSO      = 0x00000023 // Get generic segmentation offload
	ethtoolSGSO      = 0x00000024 // Set generic segmentation offload
	ethtoolGFLAGS    = 0x00000025 // Get flags bitmap (contains LRO)
	ethtoolSFLAGS    = 0x00000026 // Set flags bitmap (contains LRO)
	ethtoolGGRO      = 0x0000002b // Get generic receive offload
	ethtoolSGRO      = 0x0000002c // Set generic receive offload

	ethFlagLRO = 1 << 15

	largeReceiveOffloadName = "large-receive-offload"

	systemdNetworkDir = "/etc/systemd/network"
	sysClassNetDir    = "/sys/class/net"
)

// struct ethtool_value
type ethtoolValue struct {
	cmd  uint32
	data uint32
}

// struct ethtool_perm_addr with an inline buffer for the address.
type ethtoolPermAddr struct {
	cmd  uint32
	size uint32
	data [32]byte
}

// struct ifreq: interface name followed by the ifr_data pointer.
// The padding makes sure the struct is at least as big as the kernel's struct ifreq.
type ifreq struct {
	name [ifNameSize]byte
	data unsafe.Pointer
	_    [16]byte
}

type offloadFeature struct {
	name   string
	getCmd uint32
	setCmd uint32
}

// Offload features handled by a dedicated get/set ethtool command.
// LRO is handled separately since it is part of the ethtool flags bitmap.
var runningOffloadFeatures = []offloadFeature{
	{name: "tcp-segmentation-offload", getCmd: ethtoolGTSO, setCmd: ethtoolSTSO},
	{name: "generic-segmentation-offload", getCmd: ethtoolGGSO, setCmd: ethtoolSGSO},
	{name: "generic-receive-offload", getCmd: ethtoolGGRO, setCmd: ethtoolSGRO},
}

// Keys in the [Link] section of the systemd '.link' file which must all be set to false
// to consider offloading as permanently deactivated.
var permanentOffloadKeys = []string{
	"TCPSegmentationOffload",
	"TCP6SegmentationOffload",
	"GenericSegmentationOffload",
	"GenericReceiveOffload",
	"LargeReceiveOffload",
}

// IsOffloadingDeactivated returns true only if offloading is deactivated both running and permanent for the network card 'networkCardName':
//   - running: None of TSO, GSO, GRO or LRO is currently enabled on the network card.
//   - permanent: The systemd '.link' file '/etc/systemd/network/10-<networkCardName>-no-offload.link' exists
//     and sets all of TCPSegmentationOffload, TCP6SegmentationOffload, GenericSegmentationOffload,
//     GenericReceiveOffload and LargeReceiveOffload to false in its [Link] section.
//
// This is the state ensured by DeactivateOffloading.
// Note: This is NOT the negation of IsOffloadingActive, which only returns true if offloading is active both running and permanent.
//
// Only works on Linux.
// Does not require root privileges: The used ethtool get commands are allowed for unprivileged users
// and the '.link' file is world readable by default.
func IsOffloadingDeactivated(ctx context.Context, networkCardName string) (bool, error) {
	err := checkNetworkCardName(networkCardName)
	if err != nil {
		return false, err
	}

	err = checkRunningOnLinux("IsOffloadingDeactivated")
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Check if offloading is deactivated for network card '%s' started.", networkCardName)

	_, err = net.InterfaceByName(networkCardName)
	if err != nil {
		return false, tracederrors.TracedErrorf("Unable to get network card '%s': %w", networkCardName, err)
	}

	fd, err := openEthtoolSocket(networkCardName)
	if err != nil {
		return false, err
	}
	defer unix.Close(fd)

	activeRunning, err := isOffloadingActiveRunning(ctx, fd, networkCardName)
	if err != nil {
		return false, err
	}

	activePermanent, err := isOffloadingActivePermanent(ctx, networkCardName)
	if err != nil {
		return false, err
	}

	ret := !activeRunning && !activePermanent

	logging.LogInfoByCtxf(
		ctx,
		"Check if offloading is deactivated for network card '%s' finished. Running deactivated: '%t', permanent deactivated: '%t', offloading deactivated: '%t'.",
		networkCardName,
		!activeRunning,
		!activePermanent,
		ret,
	)

	return ret, nil
}

// DeactivateOffloading deactivates the offloading features TCP segmentation offload (TSO/TSO6),
// generic segmentation offload (GSO), generic receive offload (GRO) and large receive offload (LRO)
// for the network card 'networkCardName'.
//
// The change is applied both:
//   - running: Directly on the network card using the ethtool ioctl interface (same as 'ethtool -K <nic> tso off gso off gro off lro off').
//   - permanent: By writing a systemd '.link' file to '/etc/systemd/network/10-<networkCardName>-no-offload.link'
//     which is applied by udev every time the network card shows up (e.g. after a reboot or driver reload).
//
// This is a known workaround for e.g. the e1000e "Detected Hardware Unit Hang" issue.
//
// The function is idempotent: Already deactivated features and an already up to date '.link' file are left untouched.
//
// Only works on Linux.
// Requires root privileges: Changing offload settings needs CAP_NET_ADMIN and writing to '/etc/systemd/network' needs root.
// The root check is performed before any change is applied.
func DeactivateOffloading(ctx context.Context, networkCardName string) error {
	err := checkNetworkCardName(networkCardName)
	if err != nil {
		return err
	}

	err = checkRunningOnLinux("DeactivateOffloading")
	if err != nil {
		return err
	}

	isRoot, err := userutils.IsRunningAsRoot(ctx)
	if err != nil {
		return err
	}

	if !isRoot {
		return tracederrors.TracedErrorf("Deactivate offloading for network card '%s' requires root privileges.", networkCardName)
	}

	logging.LogInfoByCtxf(ctx, "Deactivate offloading for network card '%s' started.", networkCardName)

	iface, err := net.InterfaceByName(networkCardName)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to get network card '%s': %w", networkCardName, err)
	}

	fd, err := openEthtoolSocket(networkCardName)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	err = deactivateOffloadingRunning(ctx, fd, networkCardName)
	if err != nil {
		return err
	}

	err = deactivateOffloadingPermanent(ctx, fd, iface)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Deactivate offloading for network card '%s' finished. Offloading is deactivated running and permanent.", networkCardName)

	return nil
}

// IsOffloadingActive returns true only if offloading is active both running and permanent for the network card 'networkCardName':
//   - running: At least one of TSO, GSO, GRO or LRO is currently enabled on the network card.
//   - permanent: The systemd '.link' file '/etc/systemd/network/10-<networkCardName>-no-offload.link' is missing
//     or does not set all of TCPSegmentationOffload, TCP6SegmentationOffload, GenericSegmentationOffload,
//     GenericReceiveOffload and LargeReceiveOffload to false in its [Link] section.
//
// If offloading is deactivated running, permanent or both, false is returned.
//
// Only works on Linux.
// Does not require root privileges: The used ethtool get commands are allowed for unprivileged users
// and the '.link' file is world readable by default.
func IsOffloadingActive(ctx context.Context, networkCardName string) (bool, error) {
	err := checkNetworkCardName(networkCardName)
	if err != nil {
		return false, err
	}

	err = checkRunningOnLinux("IsOffloadingActive")
	if err != nil {
		return false, err
	}

	logging.LogInfoByCtxf(ctx, "Check if offloading is active for network card '%s' started.", networkCardName)

	_, err = net.InterfaceByName(networkCardName)
	if err != nil {
		return false, tracederrors.TracedErrorf("Unable to get network card '%s': %w", networkCardName, err)
	}

	fd, err := openEthtoolSocket(networkCardName)
	if err != nil {
		return false, err
	}
	defer unix.Close(fd)

	activeRunning, err := isOffloadingActiveRunning(ctx, fd, networkCardName)
	if err != nil {
		return false, err
	}

	activePermanent, err := isOffloadingActivePermanent(ctx, networkCardName)
	if err != nil {
		return false, err
	}

	ret := activeRunning && activePermanent

	logging.LogInfoByCtxf(
		ctx,
		"Check if offloading is active for network card '%s' finished. Running active: '%t', permanent active: '%t', offloading active: '%t'.",
		networkCardName,
		activeRunning,
		activePermanent,
		ret,
	)

	return ret, nil
}

func checkNetworkCardName(networkCardName string) error {
	if networkCardName == "" {
		return tracederrors.TracedErrorEmptyString("networkCardName")
	}

	if len(networkCardName) >= ifNameSize {
		return tracederrors.TracedErrorf("Invalid networkCardName '%s': Must be shorter than '%d' characters.", networkCardName, ifNameSize)
	}

	if strings.ContainsAny(networkCardName, "/ \t\n\r") || networkCardName == "." || networkCardName == ".." {
		return tracederrors.TracedErrorf("Invalid networkCardName '%s': Contains invalid characters.", networkCardName)
	}

	return nil
}

func checkRunningOnLinux(functionName string) error {
	if runtime.GOOS != "linux" {
		return tracederrors.TracedErrorf("%s is only supported on linux, but running on '%s'.", functionName, runtime.GOOS)
	}

	return nil
}

func openEthtoolSocket(networkCardName string) (int, error) {
	ret, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return -1, tracederrors.TracedErrorf("Unable to open socket for ethtool ioctl on network card '%s': %w", networkCardName, err)
	}

	return ret, nil
}

func ethtoolIoctl(fd int, networkCardName string, data unsafe.Pointer) error {
	var ifr ifreq
	copy(ifr.name[:ifNameSize-1], networkCardName)
	ifr.data = data

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(siocEthtool), uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		return errno
	}

	return nil
}

func getEthtoolValue(fd int, networkCardName string, cmd uint32) (uint32, error) {
	value := ethtoolValue{cmd: cmd}

	err := ethtoolIoctl(fd, networkCardName, unsafe.Pointer(&value))
	if err != nil {
		return 0, tracederrors.TracedErrorf("ethtool ioctl get command '0x%x' failed for network card '%s': %w", cmd, networkCardName, err)
	}

	return value.data, nil
}

func setEthtoolValue(fd int, networkCardName string, cmd uint32, data uint32) error {
	value := ethtoolValue{cmd: cmd, data: data}

	err := ethtoolIoctl(fd, networkCardName, unsafe.Pointer(&value))
	if err != nil {
		return tracederrors.TracedErrorf("ethtool ioctl set command '0x%x' with data '0x%x' failed for network card '%s': %w", cmd, data, networkCardName, err)
	}

	return nil
}

func isOffloadingActiveRunning(ctx context.Context, fd int, networkCardName string) (bool, error) {
	ret := false

	for _, feature := range runningOffloadFeatures {
		enabled, err := getEthtoolValue(fd, networkCardName, feature.getCmd)
		if err != nil {
			return false, err
		}

		if enabled != 0 {
			logging.LogInfoByCtxf(ctx, "Offload feature '%s' of network card '%s' is active (running).", feature.name, networkCardName)
			ret = true
		} else {
			logging.LogInfoByCtxf(ctx, "Offload feature '%s' of network card '%s' is deactivated (running).", feature.name, networkCardName)
		}
	}

	flags, err := getEthtoolValue(fd, networkCardName, ethtoolGFLAGS)
	if err != nil {
		return false, err
	}

	if flags&ethFlagLRO != 0 {
		logging.LogInfoByCtxf(ctx, "Offload feature '%s' of network card '%s' is active (running).", largeReceiveOffloadName, networkCardName)
		ret = true
	} else {
		logging.LogInfoByCtxf(ctx, "Offload feature '%s' of network card '%s' is deactivated (running).", largeReceiveOffloadName, networkCardName)
	}

	return ret, nil
}

func isOffloadingActivePermanent(ctx context.Context, networkCardName string) (bool, error) {
	linkFilePath := getLinkFilePath(networkCardName)

	content, err := os.ReadFile(linkFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.LogInfoByCtxf(ctx, "Link file '%s' does not exist. Offloading for network card '%s' is active (permanent).", linkFilePath, networkCardName)
			return true, nil
		}

		return false, tracederrors.TracedErrorf("Unable to read link file '%s': %w", linkFilePath, err)
	}

	linkSection := parseIniSection(string(content), "Link")

	ret := false
	for _, key := range permanentOffloadKeys {
		value, found := linkSection[key]
		if !found {
			logging.LogInfoByCtxf(ctx, "Key '%s' is missing in [Link] section of '%s'. Offloading for network card '%s' is active (permanent).", key, linkFilePath, networkCardName)
			ret = true
			continue
		}

		if !isSystemdBooleanFalse(value) {
			logging.LogInfoByCtxf(ctx, "Key '%s' in [Link] section of '%s' is set to '%s' instead of false. Offloading for network card '%s' is active (permanent).", key, linkFilePath, value, networkCardName)
			ret = true
		}
	}

	if !ret {
		logging.LogInfoByCtxf(ctx, "Link file '%s' deactivates all offload features for network card '%s' (permanent).", linkFilePath, networkCardName)
	}

	return ret, nil
}

// parseIniSection returns the key/value pairs of the given section of a systemd unit style ini file.
// Comments ('#' and ';') and empty lines are ignored. If a key is set multiple times the last value wins (same as systemd).
func parseIniSection(content string, section string) map[string]string {
	ret := map[string]string{}
	currentSection := ""

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			continue
		}

		if currentSection != section {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		ret[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return ret
}

// isSystemdBooleanFalse returns true for all values systemd parses as boolean false.
func isSystemdBooleanFalse(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "0", "no", "n", "false", "f", "off":
		return true
	default:
		return false
	}
}

func deactivateOffloadingRunning(ctx context.Context, fd int, networkCardName string) error {
	for _, feature := range runningOffloadFeatures {
		enabled, err := getEthtoolValue(fd, networkCardName, feature.getCmd)
		if err != nil {
			return err
		}

		if enabled == 0 {
			logging.LogInfoByCtxf(ctx, "Offload feature '%s' of network card '%s' is already deactivated (running).", feature.name, networkCardName)
			continue
		}

		err = setEthtoolValue(fd, networkCardName, feature.setCmd, 0)
		if err != nil {
			return err
		}

		enabled, err = getEthtoolValue(fd, networkCardName, feature.getCmd)
		if err != nil {
			return err
		}

		if enabled != 0 {
			return tracederrors.TracedErrorf("Offload feature '%s' of network card '%s' is still active after deactivating it. Maybe it is fixed by the driver.", feature.name, networkCardName)
		}

		logging.LogChangedByCtxf(ctx, "Offload feature '%s' of network card '%s' deactivated (running).", feature.name, networkCardName)
	}

	return deactivateLargeReceiveOffloadRunning(ctx, fd, networkCardName)
}

func deactivateLargeReceiveOffloadRunning(ctx context.Context, fd int, networkCardName string) error {
	flags, err := getEthtoolValue(fd, networkCardName, ethtoolGFLAGS)
	if err != nil {
		return err
	}

	if flags&ethFlagLRO == 0 {
		logging.LogInfoByCtxf(ctx, "Offload feature '%s' of network card '%s' is already deactivated (running).", largeReceiveOffloadName, networkCardName)
		return nil
	}

	// Keep all other flags (VLAN, NTUPLE, RXHASH) as they are and only clear LRO:
	err = setEthtoolValue(fd, networkCardName, ethtoolSFLAGS, flags&^ethFlagLRO)
	if err != nil {
		return err
	}

	flags, err = getEthtoolValue(fd, networkCardName, ethtoolGFLAGS)
	if err != nil {
		return err
	}

	if flags&ethFlagLRO != 0 {
		return tracederrors.TracedErrorf("Offload feature '%s' of network card '%s' is still active after deactivating it. Maybe it is fixed by the driver.", largeReceiveOffloadName, networkCardName)
	}

	logging.LogChangedByCtxf(ctx, "Offload feature '%s' of network card '%s' deactivated (running).", largeReceiveOffloadName, networkCardName)

	return nil
}

// getPermanentMACAddress returns the permanent (burned in) MAC address of the network card.
// An empty string is returned if the network card does not provide a permanent MAC address (e.g. virtual devices).
func getPermanentMACAddress(ctx context.Context, fd int, networkCardName string) (string, error) {
	req := ethtoolPermAddr{cmd: ethtoolGPERMADDR}
	req.size = uint32(len(req.data))

	err := ethtoolIoctl(fd, networkCardName, unsafe.Pointer(&req))
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to get permanent MAC address of network card '%s': %w", networkCardName, err)
	}

	if req.size == 0 || int(req.size) > len(req.data) {
		logging.LogInfoByCtxf(ctx, "Network card '%s' reported no usable permanent MAC address (size=%d).", networkCardName, req.size)
		return "", nil
	}

	addr := req.data[:req.size]
	if bytes.Equal(addr, make([]byte, len(addr))) {
		logging.LogInfoByCtxf(ctx, "Network card '%s' has an all zero permanent MAC address.", networkCardName)
		return "", nil
	}

	ret := net.HardwareAddr(addr).String()

	return ret, nil
}

// getDriverName returns the kernel driver name of the network card or an empty string if the device has no driver (virtual devices).
func getDriverName(ctx context.Context, networkCardName string) (string, error) {
	driverLink := filepath.Join(sysClassNetDir, networkCardName, "device", "driver")

	driverPath, err := filepath.EvalSymlinks(driverLink)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.LogInfoByCtxf(ctx, "Network card '%s' has no driver symlink '%s'. Driver match is omitted in the link file.", networkCardName, driverLink)
			return "", nil
		}

		return "", tracederrors.TracedErrorf("Unable to evaluate driver symlink '%s': %w", driverLink, err)
	}

	ret := filepath.Base(driverPath)

	return ret, nil
}

func getLinkFilePath(networkCardName string) string {
	return filepath.Join(systemdNetworkDir, "10-"+networkCardName+"-no-offload.link")
}

func generateLinkFileContent(ctx context.Context, fd int, iface *net.Interface) (string, error) {
	if iface == nil {
		return "", tracederrors.TracedErrorNil("iface")
	}

	networkCardName := iface.Name

	macKey := "PermanentMACAddress"
	macAddress, err := getPermanentMACAddress(ctx, fd, networkCardName)
	if err != nil {
		return "", err
	}

	if macAddress == "" {
		macKey = "MACAddress"
		macAddress = iface.HardwareAddr.String()
		logging.LogInfoByCtxf(ctx, "Falling back to current MAC address '%s' to match network card '%s' in the link file.", macAddress, networkCardName)
	}

	if macAddress == "" {
		return "", tracederrors.TracedErrorf("Network card '%s' has no MAC address. Unable to generate a stable match for the link file.", networkCardName)
	}

	driverName, err := getDriverName(ctx, networkCardName)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n", getLinkFilePath(networkCardName))
	sb.WriteString("# Managed by netoworkcardutils.DeactivateOffloading. Manual changes will be overwritten.\n")
	sb.WriteString("# Deactivates TSO/GSO/GRO/LRO, e.g. as workaround for e1000e \"Detected Hardware Unit Hang\".\n")
	sb.WriteString("#\n")
	sb.WriteString("# Only the first matching .link file is applied by udev. Therefore the interface name is\n")
	sb.WriteString("# pinned here as well, otherwise the default naming from 99-default.link would be lost.\n")
	sb.WriteString("\n")
	sb.WriteString("[Match]\n")
	fmt.Fprintf(&sb, "%s=%s\n", macKey, macAddress)
	if driverName != "" {
		fmt.Fprintf(&sb, "Driver=%s\n", driverName)
	}
	sb.WriteString("\n")
	sb.WriteString("[Link]\n")
	fmt.Fprintf(&sb, "Name=%s\n", networkCardName)
	sb.WriteString("AlternativeNamesPolicy=database onboard slot path mac\n")
	sb.WriteString("MACAddressPolicy=persistent\n")
	for _, key := range permanentOffloadKeys {
		fmt.Fprintf(&sb, "%s=false\n", key)
	}

	return sb.String(), nil
}

func deactivateOffloadingPermanent(ctx context.Context, fd int, iface *net.Interface) error {
	if iface == nil {
		return tracederrors.TracedErrorNil("iface")
	}

	networkCardName := iface.Name
	linkFilePath := getLinkFilePath(networkCardName)

	dirInfo, err := os.Stat(systemdNetworkDir)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to access systemd network directory '%s'. Is systemd/udev used on this host? %w", systemdNetworkDir, err)
	}

	if !dirInfo.IsDir() {
		return tracederrors.TracedErrorf("'%s' exists but is not a directory.", systemdNetworkDir)
	}

	content, err := generateLinkFileContent(ctx, fd, iface)
	if err != nil {
		return err
	}

	existing, err := os.ReadFile(linkFilePath)
	if err == nil {
		if string(existing) == content {
			logging.LogInfoByCtxf(ctx, "Link file '%s' to deactivate offloading for network card '%s' is already up to date (permanent).", linkFilePath, networkCardName)
			return nil
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return tracederrors.TracedErrorf("Unable to read existing link file '%s': %w", linkFilePath, err)
	}

	err = os.WriteFile(linkFilePath, []byte(content), 0o644)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to write link file '%s': %w", linkFilePath, err)
	}

	logging.LogChangedByCtxf(ctx, "Link file '%s' written to deactivate offloading for network card '%s' (permanent).", linkFilePath, networkCardName)

	return nil
}
