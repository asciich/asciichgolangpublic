package kvmutilsoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/filesutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	// Used when no RAM size is set in the create options.
	DefaultVmRamSizeGb = 2

	// Used when no CPU count is set in the create options.
	DefaultCpuCount = 2
)

type KvmCreateVmOptions struct {
	VmName          string
	DiskImage       filesinterfaces.File
	MacAddress      string
	BridgeInterface string
	VmRamSizeGb     int
	CpuCount        int
}

func NewKvmCreateVmOptions() (k *KvmCreateVmOptions) {
	return new(KvmCreateVmOptions)
}

func (k *KvmCreateVmOptions) GetDiskImage() (diskImage filesinterfaces.File, err error) {
	if k.DiskImage == nil {
		return nil, tracederrors.TracedErrorf("DiskImage not set")
	}

	return k.DiskImage, nil
}

func (k *KvmCreateVmOptions) GetDiskImagePath() (diskImagePath string, err error) {
	diskImage, err := k.GetDiskImage()
	if err != nil {
		return "", err
	}

	diskImagePath, err = diskImage.GetLocalPath()
	if err != nil {
		return "", err
	}

	return diskImagePath, nil
}

func (k *KvmCreateVmOptions) GetMacAddress() (macAddress string, err error) {
	if k.MacAddress == "" {
		return "", tracederrors.TracedErrorf("MacAddress not set")
	}

	return k.MacAddress, nil
}

func (k *KvmCreateVmOptions) GetBridgeInterface() (bridgeInterface string, err error) {
	if k.BridgeInterface == "" {
		return "", tracederrors.TracedErrorf("BridgeInterface not set")
	}

	return k.BridgeInterface, nil
}

func (k *KvmCreateVmOptions) GetBridgeInterfaceOrEmptyStringIfUnset() (string, error) {
	return k.BridgeInterface, nil
}

func (k *KvmCreateVmOptions) GetVmName() (vmName string, err error) {
	if k.VmName == "" {
		return "", tracederrors.TracedErrorf("VmName not set")
	}

	return k.VmName, nil
}

func (k *KvmCreateVmOptions) GetVmRamSizeGb() (vmRamSizeGb int, err error) {
	if k.VmRamSizeGb == 0 {
		return 0, tracederrors.TracedErrorf("VmRamSizeGb not set")
	}

	if k.VmRamSizeGb < 0 {
		return 0, tracederrors.TracedErrorf("Invalid VmRamSizeGb '%d'. Must be > 0.", k.VmRamSizeGb)
	}

	return k.VmRamSizeGb, nil
}

// Returns the RAM size in MiB (e.g. for 'virt-install --memory').
func (k *KvmCreateVmOptions) GetVmRamSizeMiB() (vmRamSizeMiB int, err error) {
	vmRamSizeGb, err := k.GetVmRamSizeGb()
	if err != nil {
		return 0, err
	}

	return vmRamSizeGb * 1024, nil
}

// Returns the RAM size in KiB (the default unit of <memory> in libvirt domain XML).
func (k *KvmCreateVmOptions) GetVmRamSizeKiB() (vmRamSizeKiB int, err error) {
	vmRamSizeMiB, err := k.GetVmRamSizeMiB()
	if err != nil {
		return 0, err
	}

	return vmRamSizeMiB * 1024, nil
}

func (k *KvmCreateVmOptions) GetCpuCount() (cpuCount int, err error) {
	if k.CpuCount == 0 {
		return 0, tracederrors.TracedErrorf("CpuCount not set")
	}

	if k.CpuCount < 0 {
		return 0, tracederrors.TracedErrorf("Invalid CpuCount '%d'. Must be > 0.", k.CpuCount)
	}

	return k.CpuCount, nil
}

func (k *KvmCreateVmOptions) SetDiskImageByPath(path string) error {
	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	file, err := filesutils.NewFileByPath(path)
	if err != nil {
		return err
	}

	return k.SetDiskImage(file)
}

func (k *KvmCreateVmOptions) SetDiskImage(diskImage filesinterfaces.File) (err error) {
	if diskImage == nil {
		return tracederrors.TracedErrorf("diskImage is nil")
	}

	k.DiskImage = diskImage

	return nil
}

func (k *KvmCreateVmOptions) SetMacAddress(macAddress string) (err error) {
	if macAddress == "" {
		return tracederrors.TracedErrorf("macAddress is empty string")
	}

	k.MacAddress = macAddress

	return nil
}

func (k *KvmCreateVmOptions) SetBridgeInterface(bridgeInterface string) (err error) {
	if bridgeInterface == "" {
		return tracederrors.TracedErrorf("bridgeInterface is empty string")
	}

	k.BridgeInterface = bridgeInterface

	return nil
}

func (k *KvmCreateVmOptions) SetVmName(vmName string) (err error) {
	if vmName == "" {
		return tracederrors.TracedErrorf("vmName is empty string")
	}

	k.VmName = vmName

	return nil
}

func (k *KvmCreateVmOptions) SetVmRamSizeGb(vmRamSizeGb int) (err error) {
	if vmRamSizeGb <= 0 {
		return tracederrors.TracedErrorf("Invalid vmRamSizeGb '%d'. Must be > 0.", vmRamSizeGb)
	}

	k.VmRamSizeGb = vmRamSizeGb

	return nil
}

func (k *KvmCreateVmOptions) SetCpuCount(cpuCount int) (err error) {
	if cpuCount <= 0 {
		return tracederrors.TracedErrorf("Invalid cpuCount '%d'. Must be > 0.", cpuCount)
	}

	k.CpuCount = cpuCount

	return nil
}

// Returns the RAM size in GB or DefaultVmRamSizeGb if unset.
// Invalid values (< 0) still result in an error.
func (k *KvmCreateVmOptions) GetVmRamSizeGbOrDefaultIfUnset() (vmRamSizeGb int, err error) {
	if k.VmRamSizeGb == 0 {
		return DefaultVmRamSizeGb, nil
	}

	return k.GetVmRamSizeGb()
}

// Returns the RAM size in MiB or the default if unset (e.g. for 'virt-install --memory').
func (k *KvmCreateVmOptions) GetVmRamSizeMiBOrDefaultIfUnset() (vmRamSizeMiB int, err error) {
	vmRamSizeGb, err := k.GetVmRamSizeGbOrDefaultIfUnset()
	if err != nil {
		return 0, err
	}

	return vmRamSizeGb * 1024, nil
}

// Returns the RAM size in KiB or the default if unset (the default unit of <memory> in libvirt domain XML).
func (k *KvmCreateVmOptions) GetVmRamSizeKiBOrDefaultIfUnset() (vmRamSizeKiB int, err error) {
	vmRamSizeMiB, err := k.GetVmRamSizeMiBOrDefaultIfUnset()
	if err != nil {
		return 0, err
	}

	return vmRamSizeMiB * 1024, nil
}

// Returns the CPU count or DefaultCpuCount if unset.
// Invalid values (< 0) still result in an error.
func (k *KvmCreateVmOptions) GetCpuCountOrDefaultIfUnset() (cpuCount int, err error) {
	if k.CpuCount == 0 {
		return DefaultCpuCount, nil
	}

	return k.GetCpuCount()
}
