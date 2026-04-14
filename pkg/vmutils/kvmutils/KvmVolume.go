package kvmutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmVolume struct {
	name        string
	storagePool *KvmStoragePool
}

func NewKvmVolume() (kvmVolume *KvmVolume) {
	return new(KvmVolume)
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

func (v *KvmVolume) Remove(ctx context.Context) (err error) {
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

	_, err = hypervisor.RunKvmCommand(ctx, []string{"vol-delete", "--pool", poolName, volumeName})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "KVM volume '%s' on storage pool '%s' on host '%v' deleted.", volumeName, poolName, hostname)

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
