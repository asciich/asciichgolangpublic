package vmoptions

import (
	"slices"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
)

type CreateVmOptions struct {
	VmName          string
	DiskImageName   string
	MacAddress      string
	BridgeInterface string
	VmRamSizeGb     int
	CpuCount        int

	// Optional host block devices (e.g. '/dev/sda', '/dev/sda1' or '/dev/mapper/abc')
	// passed to the VM as additional raw virtio disks (vdb, vdc, ... in the given order).
	AdditionalBlockDevices []kvmutilsoptions.AdditionalBlockDevice
}

func NewCreateVmOptions() (c *CreateVmOptions) {
	return new(CreateVmOptions)
}

func (c *CreateVmOptions) GetVmName() (vmName string, err error) {
	if c.VmName == "" {
		return "", tracederrors.TracedErrorf("VmName not set")
	}

	return c.VmName, nil
}

func (c *CreateVmOptions) GetDiskImageName() (diskImageName string, err error) {
	if c.DiskImageName == "" {
		return "", tracederrors.TracedErrorf("DiskImageName not set")
	}

	return c.DiskImageName, nil
}

func (c *CreateVmOptions) GetMacAddress() (macAddress string, err error) {
	if c.MacAddress == "" {
		return "", tracederrors.TracedErrorf("MacAddress not set")
	}

	return c.MacAddress, nil
}

func (c *CreateVmOptions) GetBridgeInterface() (bridgeInterface string, err error) {
	if c.BridgeInterface == "" {
		return "", tracederrors.TracedErrorf("BridgeInterface not set")
	}

	return c.BridgeInterface, nil
}

func (c *CreateVmOptions) GetVmRamSizeGb() (vmRamSizeGb int, err error) {
	if c.VmRamSizeGb == 0 {
		return 0, tracederrors.TracedErrorf("VmRamSizeGb not set")
	}

	if c.VmRamSizeGb < 0 {
		return 0, tracederrors.TracedErrorf("Invalid VmRamSizeGb '%d'. Must be > 0.", c.VmRamSizeGb)
	}

	return c.VmRamSizeGb, nil
}

func (c *CreateVmOptions) GetCpuCount() (cpuCount int, err error) {
	if c.CpuCount == 0 {
		return 0, tracederrors.TracedErrorf("CpuCount not set")
	}

	if c.CpuCount < 0 {
		return 0, tracederrors.TracedErrorf("Invalid CpuCount '%d'. Must be > 0.", c.CpuCount)
	}

	return c.CpuCount, nil
}

// GetAdditionalBlockDevices returns a copy of the requested block devices.
// Block devices are optional, so an empty slice is returned when none are set.
func (c *CreateVmOptions) GetAdditionalBlockDevices() []kvmutilsoptions.AdditionalBlockDevice {
	ret := slices.Clone(c.AdditionalBlockDevices)
	if ret == nil {
		ret = []kvmutilsoptions.AdditionalBlockDevice{}
	}

	return ret
}

func (c *CreateVmOptions) SetVmName(vmName string) (err error) {
	if vmName == "" {
		return tracederrors.TracedErrorf("vmName is empty string")
	}

	c.VmName = vmName

	return nil
}

func (c *CreateVmOptions) SetDiskImageName(diskImageName string) (err error) {
	if diskImageName == "" {
		return tracederrors.TracedErrorf("diskImageName is empty string")
	}

	c.DiskImageName = diskImageName

	return nil
}

func (c *CreateVmOptions) SetMacAddress(macAddress string) (err error) {
	if macAddress == "" {
		return tracederrors.TracedErrorf("macAddress is empty string")
	}

	c.MacAddress = macAddress

	return nil
}

func (c *CreateVmOptions) SetBridgeInterface(bridgeInterface string) (err error) {
	if bridgeInterface == "" {
		return tracederrors.TracedErrorf("bridgeInterface is empty string")
	}

	c.BridgeInterface = bridgeInterface

	return nil
}

func (c *CreateVmOptions) SetVmRamSizeGb(vmRamSizeGb int) (err error) {
	if vmRamSizeGb <= 0 {
		return tracederrors.TracedErrorf("Invalid vmRamSizeGb '%d'. Must be > 0.", vmRamSizeGb)
	}

	c.VmRamSizeGb = vmRamSizeGb

	return nil
}

func (c *CreateVmOptions) SetCpuCount(cpuCount int) (err error) {
	if cpuCount <= 0 {
		return tracederrors.TracedErrorf("Invalid cpuCount '%d'. Must be > 0.", cpuCount)
	}

	c.CpuCount = cpuCount

	return nil
}

// SetAdditionalBlockDevices replaces all requested block devices.
// The format is validated so invalid block devices fail early. nil/ empty removes all block devices.
func (c *CreateVmOptions) SetAdditionalBlockDevices(blockDevices []kvmutilsoptions.AdditionalBlockDevice) (err error) {
	_, err = kvmutilsgeneric.GetBlockDeviceAttachments(blockDevices)
	if err != nil {
		return err
	}

	// Clone so later modifications of the caller's slice do not affect the options.
	c.AdditionalBlockDevices = slices.Clone(blockDevices)

	return nil
}

// AddAdditionalBlockDevice appends a block device (e.g. for a repeatable CLI flag).
// The whole list is validated so duplicates of source paths or serials are rejected.
func (c *CreateVmOptions) AddAdditionalBlockDevice(blockDevice kvmutilsoptions.AdditionalBlockDevice) (err error) {
	newBlockDevices := append(slices.Clone(c.AdditionalBlockDevices), blockDevice)

	return c.SetAdditionalBlockDevices(newBlockDevices)
}

func (c *CreateVmOptions) GetBridgeInterfaceOrEmptyStringIfUnset() (string, error) {
	return c.BridgeInterface, nil
}

// Returns the RAM size in GB or kvmutilsoptions.DefaultVmRamSizeGb if unset.
func (c *CreateVmOptions) GetVmRamSizeGbOrDefaultIfUnset() (vmRamSizeGb int, err error) {
	if c.VmRamSizeGb == 0 {
		return kvmutilsoptions.DefaultVmRamSizeGb, nil
	}

	return c.GetVmRamSizeGb()
}

// Returns the CPU count or kvmutilsoptions.DefaultCpuCount if unset.
func (c *CreateVmOptions) GetCpuCountOrDefaultIfUnset() (cpuCount int, err error) {
	if c.CpuCount == 0 {
		return kvmutilsoptions.DefaultCpuCount, nil
	}

	return c.GetCpuCount()
}
