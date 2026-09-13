package kvmutilsgeneric

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsinterfaces"
)

type KvmStoragePool struct {
	name       string
	hypervisor kvmutilsinterfaces.Hypervisor
}

func NewKvmStoragePool() (kvmStoragePool *KvmStoragePool) {
	return new(KvmStoragePool)
}

func (k *KvmStoragePool) GetHostName() (hostname string, err error) {
	hypervisor, err := k.GetHypervisor()
	if err != nil {
		return "", err
	}

	hostname, err = hypervisor.GetHostName()
	if err != nil {
		return "", err
	}

	return hostname, nil
}

func (k *KvmStoragePool) GetHypervisor() (hypervisor kvmutilsinterfaces.Hypervisor, err error) {
	if k.hypervisor == nil {
		return nil, tracederrors.TracedError("hypervisor not set")
	}

	return k.hypervisor, nil
}

func (k *KvmStoragePool) GetName() (name string, err error) {
	if len(k.name) <= 0 {
		return "", tracederrors.TracedError("name not set")
	}

	return k.name, nil
}

func (k *KvmStoragePool) ListVolumes(ctx context.Context) (volumes []kvmutilsinterfaces.Volume, err error) {
	storagePoolName, err := k.GetName()
	if err != nil {
		return nil, err
	}

	hypervisor, err := k.GetHypervisor()
	if err != nil {
		return nil, err
	}

	return hypervisor.ListVolumes(ctx, storagePoolName)
}

func (k *KvmStoragePool) SetHypervisor(hypervisor kvmutilsinterfaces.Hypervisor) (err error) {
	if hypervisor == nil {
		return tracederrors.TracedError("hypervisor is nil")
	}

	k.hypervisor = hypervisor

	return nil
}

func (k *KvmStoragePool) SetName(name string) (err error) {
	if len(name) <= 0 {
		return tracederrors.TracedError("name is nil")
	}

	k.name = name

	return nil
}
