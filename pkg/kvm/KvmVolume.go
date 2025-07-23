package kvm

import (
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

type KvmVolume struct {
	name        string
	storagePool *KvmStoragePool
}

func NewKvmVolume() (kvmVolume *KvmVolume) {
	return new(KvmVolume)
}

func (k *KvmVolume) MustGetHostName() (hostname string) {
	hostname, err := k.GetHostName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostname
}

func (k *KvmVolume) MustGetHypervisor() (hypervisor *KVMHypervisor) {
	hypervisor, err := k.GetHypervisor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hypervisor
}

func (k *KvmVolume) MustGetName() (name string) {
	name, err := k.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (k *KvmVolume) MustGetStoragePool() (storagePool *KvmStoragePool) {
	storagePool, err := k.GetStoragePool()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return storagePool
}

func (k *KvmVolume) MustGetStoragePoolName() (storagePoolName string) {
	storagePoolName, err := k.GetStoragePoolName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return storagePoolName
}

func (k *KvmVolume) MustRemove(verbose bool) {
	err := k.Remove(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmVolume) MustSetName(name string) {
	err := k.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmVolume) MustSetStoragePool(storagePool *KvmStoragePool) {
	err := k.SetStoragePool(storagePool)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (v *KvmVolume) GetHostName() (hostname string, err error) {
	pool, err := v.GetStoragePool()
	if err != nil {
		return "", err
	}

	hostname, err = pool.GetHostName()
	if err != nil {
		return "", err
	}

	return hostname, nil
}

func (v *KvmVolume) GetHypervisor() (hypervisor *KVMHypervisor, err error) {
	pool, err := v.GetStoragePool()
	if err != nil {
		return nil, err
	}

	hypervisor, err = pool.GetHypervisor()
	if err != nil {
		return nil, err
	}

	return hypervisor, nil
}

func (v *KvmVolume) GetName() (name string, err error) {
	if len(v.name) <= 0 {
		return "", tracederrors.TracedError("name not set")
	}

	return v.name, nil
}

func (v *KvmVolume) GetStoragePool() (storagePool *KvmStoragePool, err error) {
	if v.storagePool == nil {
		return nil, tracederrors.TracedError("storage pool not set")
	}

	return v.storagePool, nil
}

func (v *KvmVolume) GetStoragePoolName() (storagePoolName string, err error) {
	pool, err := v.GetStoragePool()
	if err != nil {
		return "", err
	}

	storagePoolName, err = pool.GetName()
	if err != nil {
		return "", err
	}

	return storagePoolName, nil
}

func (v *KvmVolume) Remove(verbose bool) (err error) {
	hostname, err := v.GetHostName()
	if err != nil {
		return err
	}

	hypervisor, err := v.GetHypervisor()
	if err != nil {
		return err
	}

	volumeName, err := v.GetName()
	if err != nil {
		return err
	}

	poolName, err := v.GetStoragePoolName()
	if err != nil {
		return err
	}

	_, err = hypervisor.RunKvmCommand([]string{"vol-delete", "--pool", poolName, volumeName}, verbose)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("KVM volume '%s' on storage pool '%s' on host '%v' deleted.", volumeName, poolName, hostname)
	}

	return nil
}

func (v *KvmVolume) SetName(name string) (err error) {
	if len(name) <= 0 {
		return tracederrors.TracedError("name is empty string")
	}

	v.name = name

	return nil
}

func (v *KvmVolume) SetStoragePool(storagePool *KvmStoragePool) (err error) {
	if storagePool == nil {
		return tracederrors.TracedError("storagePool is nil")
	}

	v.storagePool = storagePool

	return nil
}
