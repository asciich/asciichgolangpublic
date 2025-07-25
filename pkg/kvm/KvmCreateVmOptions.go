package kvm

import (
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmCreateVmOptions struct {
	VmName     string
	DiskImage  filesinterfaces.File
	Verbose    bool
	MacAddress string
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

func (k *KvmCreateVmOptions) GetVerbose() (verbose bool, err error) {

	return k.Verbose, nil
}

func (k *KvmCreateVmOptions) GetVmName() (vmName string, err error) {
	if k.VmName == "" {
		return "", tracederrors.TracedErrorf("VmName not set")
	}

	return k.VmName, nil
}

func (k *KvmCreateVmOptions) MustGetDiskImage() (diskImage filesinterfaces.File) {
	diskImage, err := k.GetDiskImage()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return diskImage
}

func (k *KvmCreateVmOptions) MustGetDiskImagePath() (diskImagePath string) {
	diskImagePath, err := k.GetDiskImagePath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return diskImagePath
}

func (k *KvmCreateVmOptions) MustGetMacAddress() (macAddress string) {
	macAddress, err := k.GetMacAddress()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return macAddress
}

func (k *KvmCreateVmOptions) MustGetVerbose() (verbose bool) {
	verbose, err := k.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (k *KvmCreateVmOptions) MustGetVmName() (vmName string) {
	vmName, err := k.GetVmName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vmName
}

func (k *KvmCreateVmOptions) MustSetDiskImage(diskImage filesinterfaces.File) {
	err := k.SetDiskImage(diskImage)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmCreateVmOptions) MustSetMacAddress(macAddress string) {
	err := k.SetMacAddress(macAddress)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmCreateVmOptions) MustSetVerbose(verbose bool) {
	err := k.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmCreateVmOptions) MustSetVmName(vmName string) {
	err := k.SetVmName(vmName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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

func (k *KvmCreateVmOptions) SetVerbose(verbose bool) (err error) {
	k.Verbose = verbose

	return nil
}

func (k *KvmCreateVmOptions) SetVmName(vmName string) (err error) {
	if vmName == "" {
		return tracederrors.TracedErrorf("vmName is empty string")
	}

	k.VmName = vmName

	return nil
}
