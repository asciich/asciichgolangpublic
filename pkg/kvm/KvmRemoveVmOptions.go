package kvm

import (
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type KvmRemoveVmOptions struct {
	VmName              string
	RemoveVolumes       bool
	Verbose             bool
	VolumeNamesToRemove []string
}

func NewKvmRemoveVmOptions() (k *KvmRemoveVmOptions) {
	return new(KvmRemoveVmOptions)
}

func (k *KvmRemoveVmOptions) GetRemoveVolumes() (removeVolumes bool, err error) {

	return k.RemoveVolumes, nil
}

func (k *KvmRemoveVmOptions) GetVerbose() (verbose bool, err error) {

	return k.Verbose, nil
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

func (k *KvmRemoveVmOptions) MustGetRemoveVolumes() (removeVolumes bool) {
	removeVolumes, err := k.GetRemoveVolumes()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return removeVolumes
}

func (k *KvmRemoveVmOptions) MustGetVerbose() (verbose bool) {
	verbose, err := k.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (k *KvmRemoveVmOptions) MustGetVmName() (vmName string) {
	vmName, err := k.GetVmName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vmName
}

func (k *KvmRemoveVmOptions) MustGetVolumeNamesToRemove() (volumeNamesToRemove []string) {
	volumeNamesToRemove, err := k.GetVolumeNamesToRemove()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return volumeNamesToRemove
}

func (k *KvmRemoveVmOptions) MustSetRemoveVolumes(removeVolumes bool) {
	err := k.SetRemoveVolumes(removeVolumes)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmRemoveVmOptions) MustSetVerbose(verbose bool) {
	err := k.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmRemoveVmOptions) MustSetVmName(vmName string) {
	err := k.SetVmName(vmName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmRemoveVmOptions) MustSetVolumeNamesToRemove(volumeNamesToRemove []string) {
	err := k.SetVolumeNamesToRemove(volumeNamesToRemove)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmRemoveVmOptions) SetRemoveVolumes(removeVolumes bool) (err error) {
	k.RemoveVolumes = removeVolumes

	return nil
}

func (k *KvmRemoveVmOptions) SetVerbose(verbose bool) (err error) {
	k.Verbose = verbose

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
