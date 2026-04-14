package kvmutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmRemoveVmOptions struct {
	VmName              string
	RemoveVolumes       bool
	VolumeNamesToRemove []string
}

func NewKvmRemoveVmOptions() (k *KvmRemoveVmOptions) {
	return new(KvmRemoveVmOptions)
}

func (k *KvmRemoveVmOptions) GetRemoveVolumes() (removeVolumes bool, err error) {

	return k.RemoveVolumes, nil
}

func (k *KvmRemoveVmOptions) GetVmName() (vmName string, err error) {
	if k.VmName == "" {
		return "", tracederrors.TracedErrorf("VmName not set")
	}

	return k.VmName, nil
}

func (k *KvmRemoveVmOptions) GetVolumeNamesToRemove() (volumeNamesToRemove []string, err error) {
	if k.VolumeNamesToRemove == nil {
		return nil, tracederrors.TracedErrorf("VolumeNamesToRemove not set")
	}

	if len(k.VolumeNamesToRemove) <= 0 {
		return nil, tracederrors.TracedErrorf("VolumeNamesToRemove has no elements")
	}

	return k.VolumeNamesToRemove, nil
}

func (k *KvmRemoveVmOptions) SetRemoveVolumes(removeVolumes bool) (err error) {
	k.RemoveVolumes = removeVolumes

	return nil
}

func (k *KvmRemoveVmOptions) SetVmName(vmName string) (err error) {
	if vmName == "" {
		return tracederrors.TracedErrorf("vmName is empty string")
	}

	k.VmName = vmName

	return nil
}

func (k *KvmRemoveVmOptions) SetVolumeNamesToRemove(volumeNamesToRemove []string) (err error) {
	if volumeNamesToRemove == nil {
		return tracederrors.TracedErrorf("volumeNamesToRemove is nil")
	}

	if len(volumeNamesToRemove) <= 0 {
		return tracederrors.TracedErrorf("volumeNamesToRemove has no elements")
	}

	k.VolumeNamesToRemove = volumeNamesToRemove

	return nil
}
