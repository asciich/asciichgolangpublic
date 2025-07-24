package kvm

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmStoragePool struct {
	name       string
	hypervisor *KVMHypervisor
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

func (k *KvmStoragePool) GetHypervisor() (hypervisor *KVMHypervisor, err error) {
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

func (k *KvmStoragePool) GetVolumes(verbose bool) (volumes []*KvmVolume, err error) {
	poolName, err := k.GetName()
	if err != nil {
		return nil, err
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	if verbose {
		logging.LogInfof("Get volumes in storage pool '%s' on kvm hypervisor '%s' started.", poolName, hostname)
	}

	hypervisor, err := k.GetHypervisor()
	if err != nil {
		return nil, err
	}

	listPoolOutput, err := hypervisor.RunKvmCommandAndGetStdout([]string{"vol-list", "--pool", poolName}, verbose)
	if err != nil {
		return nil, err
	}

	firstLine, unparsedOutput := stringsutils.SplitFirstLineAndContent(listPoolOutput)
	firstLine = strings.TrimSpace(firstLine)
	if !strings.HasPrefix(firstLine, "Name") {
		return nil, tracederrors.TracedErrorf("Unexpected first line of list volumes output: '%s'", firstLine)
	}

	secondLine, unparsedOutput := stringsutils.SplitFirstLineAndContent(unparsedOutput)
	secondLine = strings.TrimSpace(secondLine)
	if strings.Count(secondLine, "-") < 5 {
		return nil, tracederrors.TracedErrorf("Unexpected second line of list volumes output: '%s'", secondLine)
	}

	volumes = []*KvmVolume{}
	for _, line := range stringsutils.SplitLines(unparsedOutput, true) {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		splitted := stringsutils.SplitAtSpacesAndRemoveEmptyStrings(line)
		if len(splitted) != 2 {
			return nil, tracederrors.TracedErrorf("Unable to splitt list volume line '%v' : '%v'", line, splitted)
		}

		nameToAdd := splitted[0]
		volumeToAdd := NewKvmVolume()
		err = volumeToAdd.SetName(nameToAdd)
		if err != nil {
			return nil, err
		}

		err = volumeToAdd.SetStoragePool(k)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, volumeToAdd)
	}

	if verbose {
		logging.LogInfof("Collected '%d' storage pools on kvm host '%s'", len(volumes), hostname)
	}

	if verbose {
		logging.LogInfof("Get volumes in storage pool '%s' on kvm hypervisor '%s' finished.", poolName, hostname)
	}

	return volumes, nil
}

func (k *KvmStoragePool) MustGetHostName() (hostname string) {
	hostname, err := k.GetHostName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostname
}

func (k *KvmStoragePool) MustGetHypervisor() (hypervisor *KVMHypervisor) {
	hypervisor, err := k.GetHypervisor()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hypervisor
}

func (k *KvmStoragePool) MustGetName() (name string) {
	name, err := k.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (k *KvmStoragePool) MustGetVolumes(verbose bool) (volumes []*KvmVolume) {
	volumes, err := k.GetVolumes(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return volumes
}

func (k *KvmStoragePool) MustSetHypervisor(hypervisor *KVMHypervisor) {
	err := k.SetHypervisor(hypervisor)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmStoragePool) MustSetName(name string) {
	err := k.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmStoragePool) SetHypervisor(hypervisor *KVMHypervisor) (err error) {
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
