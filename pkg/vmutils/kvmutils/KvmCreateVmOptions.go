package kvmutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmCreateVmOptions struct {
	VmName     string
	DiskImage  filesinterfaces.File
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

func (k *KvmCreateVmOptions) GetVmName() (vmName string, err error) {
	if k.VmName == "" {
		return "", tracederrors.TracedErrorf("VmName not set")
	}

	return k.VmName, nil
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

func (k *KvmCreateVmOptions) SetVmName(vmName string) (err error) {
	if vmName == "" {
		return tracederrors.TracedErrorf("vmName is empty string")
	}

	k.VmName = vmName

	return nil
}
