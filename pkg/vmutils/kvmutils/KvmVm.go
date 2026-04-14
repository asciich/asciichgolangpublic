package kvmutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmVm struct {
	vmId       *int
	cachedName string
	hypervisor *KVMHypervisor
}

func NewKvmVm() (kvmVm *KvmVm) {
	return new(KvmVm)
}

func (k *KvmVm) GetCachedName() (cachedName string, err error) {
	if len(k.cachedName) <= 0 {
		k.cachedName, err = k.GetName()
		if err != nil {
			return "", err
		}
	}

	if len(k.cachedName) <= 0 {
		return "", tracederrors.TracedError("Unable to load cached name")
	}

	return k.cachedName, nil
}

func (k *KvmVm) GetDomainXmlAsString(ctx context.Context) (domainXml string, err error) {
	hypervisor, err := k.GetHypervisor()
	if err != nil {
		return "", err
	}

	vmName, err := k.GetCachedName()
	if err != nil {
		return "", err
	}

	domainXml, err = hypervisor.RunKvmCommandAndGetStdout(ctx, []string{"dumpxml", vmName})
	if err != nil {
		return "", err
	}

	return domainXml, nil
}

func (k *KvmVm) GetHypervisor() (hypervisor *KVMHypervisor, err error) {
	if k.hypervisor == nil {
		return nil, tracederrors.TracedErrorf("hypervisor not set")
	}

	return k.hypervisor, nil
}

func (k *KvmVm) GetId() (id int, err error) {
	if k.vmId == nil {
		return -1, tracederrors.TracedError("name is not set")
	}

	return *(k.vmId), nil
}

func (k *KvmVm) GetInfo(ctx context.Context) (vmInfo *KvmVmInfo, err error) {
	vmInfo = NewKvmVmInfo()

	vmName, err := k.GetCachedName()
	if err != nil {
		return nil, err
	}

	err = vmInfo.SetName(vmName)
	if err != nil {
		return nil, err
	}

	macAddress, err := k.GetMacAddress(ctx)
	if err != nil {
		return nil, err
	}

	err = vmInfo.SetMacAddress(macAddress)
	if err != nil {
		return nil, err
	}

	return vmInfo, nil
}

func (k *KvmVm) GetMacAddress(ctx context.Context) (macAddress string, err error) {
	domainXml, err := k.GetDomainXmlAsString(ctx)
	if err != nil {
		return "", err
	}

	macAddress, err = LibvirtXmls().GetMacAddressFromXmlString(domainXml)
	if err != nil {
		return "", err
	}

	return macAddress, nil
}

func (k *KvmVm) GetName() (name string, err error) {
	return "", tracederrors.TracedError("Not implemented")
}

func (k *KvmVm) GetVmId() (vmId *int, err error) {
	if k.vmId == nil {
		return nil, tracederrors.TracedErrorf("vmId not set")
	}

	return k.vmId, nil
}

func (k *KvmVm) SetCachedName(cachedName string) (err error) {
	if len(cachedName) <= 0 {
		return tracederrors.TracedError("cachedName is empty string")
	}

	k.cachedName = cachedName

	return nil
}

func (k *KvmVm) SetHypervisor(hypervisor *KVMHypervisor) (err error) {
	if hypervisor == nil {
		return tracederrors.TracedErrorf("hypervisor is nil")
	}

	k.hypervisor = hypervisor

	return nil
}

func (k *KvmVm) SetId(id int) (err error) {
	if id < 0 {
		return tracederrors.TracedErrorf("invalid id '%d'", id)
	}

	idToAdd := id

	k.vmId = &idToAdd

	return nil
}

func (k *KvmVm) SetVmId(vmId *int) (err error) {
	if vmId == nil {
		return tracederrors.TracedErrorf("vmId is nil")
	}

	k.vmId = vmId

	return nil
}
