package kvmutilsgeneric

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
)

// Max length of a virtio-blk serial. Longer serials are truncated by QEMU which would break
// the stable '/dev/disk/by-id/virtio-<serial>' path inside the VM.
const MaxVirtioBlkSerialLength = 20

// Allow list for block device paths. Since the paths are rendered into the domain XML using
// text/template (no XML escaping) characters like ', ", <, > and & must never be accepted.
// ':', '+', '@', '=' and '\' are needed for '/dev/disk/by-*' paths (e.g. 'pci-0000:00:17.0-ata-1').
var blockDeviceSourcePathRegex = regexp.MustCompile(`^/dev/[A-Za-z0-9._:+@=\\/-]+$`)

var virtioBlkSerialRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// BlockDeviceAttachment is a host block device including the assigned target device name in the VM.
type BlockDeviceAttachment struct {
	SourcePath string
	TargetDev  string
	Serial     string
	ReadOnly   bool
}

// GetVirtioTargetDevName returns the virtio target device name for the additional disk with the given index.
// Index 0 returns 'vdb' since 'vda' is reserved for the root disk.
// The naming follows the Linux scheme: ..., vdz, vdaa, vdab, ...
func GetVirtioTargetDevName(index int) (string, error) {
	if index < 0 {
		return "", tracederrors.TracedErrorf("index must be >= 0 but is '%d'.", index)
	}

	n := index + 1 // Skip 'vda' which is reserved for the root disk.

	suffix := ""
	for {
		suffix = string(rune('a'+n%26)) + suffix
		n = n/26 - 1
		if n < 0 {
			break
		}
	}

	ret := "vd" + suffix

	return ret, nil
}

// CheckBlockDeviceSourcePath validates the format of a host block device path.
// This is a pure format check. Whether the path really is a block device on the host
// has to be checked by the native/commandexecutor implementations.
func CheckBlockDeviceSourcePath(sourcePath string) error {
	if sourcePath == "" {
		return tracederrors.TracedErrorEmptyString("sourcePath")
	}

	if !strings.HasPrefix(sourcePath, "/dev/") {
		return tracederrors.TracedErrorf("Block device source path '%s' must start with '/dev/'.", sourcePath)
	}

	// Rejects '..', '//', trailing '/' and '/dev/' itself.
	if filepath.Clean(sourcePath) != sourcePath {
		return tracederrors.TracedErrorf("Block device source path '%s' is not a clean path. Expected '%s'.", sourcePath, filepath.Clean(sourcePath))
	}

	if !blockDeviceSourcePathRegex.MatchString(sourcePath) {
		return tracederrors.TracedErrorf("Block device source path '%s' contains unsupported characters.", sourcePath)
	}

	return nil
}

// CheckVirtioBlkSerial validates a virtio-blk serial. An empty serial is valid since the serial is optional.
func CheckVirtioBlkSerial(serial string) error {
	if serial == "" {
		return nil
	}

	if len(serial) > MaxVirtioBlkSerialLength {
		return tracederrors.TracedErrorf("Serial '%s' is '%d' characters long but max '%d' are supported by virtio-blk.", serial, len(serial), MaxVirtioBlkSerialLength)
	}

	if !virtioBlkSerialRegex.MatchString(serial) {
		return tracederrors.TracedErrorf("Serial '%s' contains unsupported characters. Only [A-Za-z0-9_-] are allowed.", serial)
	}

	return nil
}

// GetBlockDeviceAttachments validates the given block devices and assigns the target device
// names ('vdb', 'vdc', ...) in the given order.
//
// Duplicates are detected by string comparison only. Different paths pointing to the same device
// (e.g. '/dev/sda' and '/dev/disk/by-id/...') or overlapping devices (e.g. '/dev/sda' and '/dev/sda1')
// have to be detected on the host by the native/commandexecutor implementations.
func GetBlockDeviceAttachments(blockDevices []kvmutilsoptions.AdditionalBlockDevice) ([]BlockDeviceAttachment, error) {
	ret := []BlockDeviceAttachment{}

	seenSourcePaths := map[string]bool{}
	seenSerials := map[string]bool{}

	for i, blockDevice := range blockDevices {
		err := CheckBlockDeviceSourcePath(blockDevice.SourcePath)
		if err != nil {
			return nil, err
		}

		err = CheckVirtioBlkSerial(blockDevice.Serial)
		if err != nil {
			return nil, err
		}

		if seenSourcePaths[blockDevice.SourcePath] {
			return nil, tracederrors.TracedErrorf("Block device '%s' is specified more than once.", blockDevice.SourcePath)
		}
		seenSourcePaths[blockDevice.SourcePath] = true

		if blockDevice.Serial != "" {
			if seenSerials[blockDevice.Serial] {
				return nil, tracederrors.TracedErrorf("Serial '%s' is used for more than one block device.", blockDevice.Serial)
			}
			seenSerials[blockDevice.Serial] = true
		}

		targetDev, err := GetVirtioTargetDevName(i)
		if err != nil {
			return nil, err
		}

		ret = append(ret, BlockDeviceAttachment{
			SourcePath: blockDevice.SourcePath,
			TargetDev:  targetDev,
			Serial:     blockDevice.Serial,
			ReadOnly:   blockDevice.ReadOnly,
		})
	}

	return ret, nil
}

// ParseAdditionalBlockDeviceString parses a block device definition in the format
// '<path>[,serial=<serial>][,ro]', e.g. '/dev/mapper/vg0-backup,serial=backup,ro'.
//
// ',' is used as separator since block device paths can not contain it (see CheckBlockDeviceSourcePath),
// while ':' is part of e.g. '/dev/disk/by-path/pci-0000:00:17.0-ata-1'.
func ParseAdditionalBlockDeviceString(blockDeviceString string) (kvmutilsoptions.AdditionalBlockDevice, error) {
	if blockDeviceString == "" {
		return kvmutilsoptions.AdditionalBlockDevice{}, tracederrors.TracedErrorEmptyString("blockDeviceString")
	}

	parts := strings.Split(blockDeviceString, ",")

	ret := kvmutilsoptions.AdditionalBlockDevice{
		SourcePath: parts[0],
	}

	err := CheckBlockDeviceSourcePath(ret.SourcePath)
	if err != nil {
		return kvmutilsoptions.AdditionalBlockDevice{}, err
	}

	seenOptions := map[string]bool{}
	for _, option := range parts[1:] {
		key, value, hasValue := strings.Cut(option, "=")

		if seenOptions[key] {
			return kvmutilsoptions.AdditionalBlockDevice{}, tracederrors.TracedErrorf("Option '%s' is specified more than once in block device '%s'.", key, blockDeviceString)
		}
		seenOptions[key] = true

		switch key {
		case "ro":
			if hasValue {
				return kvmutilsoptions.AdditionalBlockDevice{}, tracederrors.TracedErrorf("Option 'ro' takes no value but got 'ro=%s' in block device '%s'.", value, blockDeviceString)
			}

			ret.ReadOnly = true
		case "serial":
			if value == "" {
				return kvmutilsoptions.AdditionalBlockDevice{}, tracederrors.TracedErrorf("Option 'serial' requires a value like 'serial=data01' in block device '%s'.", blockDeviceString)
			}

			err = CheckVirtioBlkSerial(value)
			if err != nil {
				return kvmutilsoptions.AdditionalBlockDevice{}, err
			}

			ret.Serial = value
		default:
			return kvmutilsoptions.AdditionalBlockDevice{}, tracederrors.TracedErrorf("Unknown option '%s' in block device '%s'. Supported: 'serial=<serial>', 'ro'.", option, blockDeviceString)
		}
	}

	return ret, nil
}
