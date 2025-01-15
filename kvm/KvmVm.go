package kvm

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
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

func (k *KvmVm) GetDomainXmlAsString(verbose bool) (domainXml string, err error) {
	hypervisor, err := k.GetHypervisor()
	if err != nil {
		return "", err
	}

	vmName, err := k.GetCachedName()
	if err != nil {
		return "", err
	}

	domainXml, err = hypervisor.RunKvmCommandAndGetStdout(
		[]string{"dumpxml", vmName},
		verbose,
	)
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

func (k *KvmVm) GetInfo(verbose bool) (vmInfo *KvmVmInfo, err error) {
	vmInfo = NewKvmVmInfo()

	vmName, err := k.GetCachedName()
	if err != nil {
		return nil, err
	}

	err = vmInfo.SetName(vmName)
	if err != nil {
		return nil, err
	}

	macAddress, err := k.GetMacAddress(verbose)
	if err != nil {
		return nil, err
	}

	err = vmInfo.SetMacAddress(macAddress)
	if err != nil {
		return nil, err
	}

	return vmInfo, nil
}

func (k *KvmVm) GetMacAddress(verbose bool) (macAddress string, err error) {
	domainXml, err := k.GetDomainXmlAsString(verbose)
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

func (k *KvmVm) MustGetCachedName() (cachedName string) {
	cachedName, err := k.GetCachedName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return cachedName
}

func (k *KvmVm) MustGetDomainXmlAsString(verbose bool) (domainXml string) {
	domainXml, err := k.GetDomainXmlAsString(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return domainXml
}

func (k *KvmVm) MustGetHypervisor() (hypervisor *KVMHypervisor) {
	hypervisor, err := k.GetHypervisor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hypervisor
}

func (k *KvmVm) MustGetId() (id int) {
	id, err := k.GetId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return id
}

func (k *KvmVm) MustGetInfo(verbose bool) (vmInfo *KvmVmInfo) {
	vmInfo, err := k.GetInfo(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vmInfo
}

func (k *KvmVm) MustGetMacAddress(verbose bool) (macAddress string) {
	macAddress, err := k.GetMacAddress(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return macAddress
}

func (k *KvmVm) MustGetName() (name string) {
	name, err := k.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (k *KvmVm) MustGetVmId() (vmId *int) {
	vmId, err := k.GetVmId()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vmId
}

func (k *KvmVm) MustSetCachedName(cachedName string) {
	err := k.SetCachedName(cachedName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmVm) MustSetHypervisor(hypervisor *KVMHypervisor) {
	err := k.SetHypervisor(hypervisor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmVm) MustSetId(id int) {
	err := k.SetId(id)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmVm) MustSetVmId(vmId *int) {
	err := k.SetVmId(vmId)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
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
