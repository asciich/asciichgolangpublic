package kvm

import (
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic"
	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	astrings "github.com/asciich/asciichgolangpublic/datatypes/strings"
	"github.com/asciich/asciichgolangpublic/hosts"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type KVMHypervisor struct {
	host hosts.Host

	// Run kvm commands and connection directly on localhost instead of using SSH.
	useLocalhost bool
}

func GetKvmHypervisorByHost(host hosts.Host) (kvmHypervisor *KVMHypervisor, err error) {
	if host == nil {
		return nil, tracederrors.TracedError("host is nil")
	}

	kvmHypervisor = NewKVMHypervisor()
	err = kvmHypervisor.SetHost(host)
	if err != nil {
		return nil, err
	}

	return kvmHypervisor, nil
}

func GetKvmHypervisorOnLocalhost() (kvmHypervisor *KVMHypervisor, err error) {
	kvmHypervisor = NewKVMHypervisor()
	err = kvmHypervisor.SetUseLocalhost(true)
	if err != nil {
		return nil, err
	}

	return kvmHypervisor, nil
}

func MustGetKvmHypervisorByHost(host hosts.Host) (kvmHypervisor *KVMHypervisor) {
	kvmHypervisor, err := GetKvmHypervisorByHost(host)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return kvmHypervisor
}

func MustGetKvmHypervisorOnLocalhost() (kvmHypervisor *KVMHypervisor) {
	kvmHypervisor, err := GetKvmHypervisorOnLocalhost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return kvmHypervisor
}

func NewKVMHypervisor() (kvmHypervisor *KVMHypervisor) {
	return new(KVMHypervisor)
}

func (k *KVMHypervisor) CreateVm(createOptions *KvmCreateVmOptions) (createdVm *KvmVm, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
	}

	vmName, err := createOptions.GetVmName()
	if err != nil {
		return nil, err
	}

	exists, err := k.VmByNameExists(vmName)
	if err != nil {
		return nil, err
	}

	if exists {
		if createOptions.Verbose {
			logging.LogInfof("VM '%s' already exists", vmName)
		}

		createdVm, err = k.GetVmByName(vmName, createOptions.Verbose)
		if err != nil {
			return nil, err
		}

		return createdVm, nil
	}

	diskImage, err := createOptions.GetDiskImage()
	if err != nil {
		return nil, err
	}

	diskImagePath, err := diskImage.GetLocalPath()
	if err != nil {
		return nil, err
	}

	diskImageExists, err := diskImage.Exists(createOptions.Verbose)
	if err != nil {
		return nil, err
	}

	if !diskImageExists {
		return nil, tracederrors.TracedErrorf("Disk image '%s' does not exist to create VM.", diskImagePath)
	}

	vmXml, err := asciichgolangpublic.TemporaryFiles().CreateEmptyTemporaryFile(createOptions.Verbose)
	if err != nil {
		return nil, err
	}
	defer vmXml.Delete(createOptions.Verbose)

	err = LibvirtXmls().WriteXmlForVmOnLatopToFile(createOptions, vmXml)
	if err != nil {
		return nil, err
	}

	vmXmlPath, err := vmXml.GetLocalPath()
	if err != nil {
		return nil, err
	}

	createOutput, err := k.RunKvmCommandAndGetStdout(
		[]string{"create", vmXmlPath},
		createOptions.Verbose,
	)
	if err != nil {
		return nil, err
	}

	if createOptions.Verbose {
		logging.LogInfof("Output of VM '%s' creation:\n%s", vmName, createOutput)
	}

	createdVm, err = k.GetVmByName(vmName, createOptions.Verbose)
	if err != nil {
		return nil, err
	}

	if createOptions.Verbose {
		logging.LogChangedf("Vm '%s' created.", vmName)
	}

	return createdVm, nil
}

func (k *KVMHypervisor) GetHost() (host hosts.Host, err error) {
	if k.host == nil {
		return nil, tracederrors.TracedError("host not set")
	}

	return k.host, nil
}

func (k *KVMHypervisor) GetHostName() (hostname string, err error) {
	if k.useLocalhost {
		return "localhost_connection", nil
	}

	host, err := k.GetHost()
	if err != nil {
		return "", err
	}

	hostname, err = host.GetHostName()
	if err != nil {
		return "", err
	}

	return hostname, nil
}

func (k *KVMHypervisor) GetStoragePoolNames(verbose bool) (storagePoolNames []string, err error) {
	storagePools, err := k.GetStoragePools(verbose)
	if err != nil {
		return nil, err
	}

	storagePoolNames = []string{}
	for _, pool := range storagePools {
		nameToAdd, err := pool.GetName()
		if err != nil {
			return nil, err
		}

		storagePoolNames = append(storagePoolNames, nameToAdd)
	}

	return
}

func (k *KVMHypervisor) GetStoragePools(verbose bool) (storagePools []*KvmStoragePool, err error) {
	if verbose {
		logging.LogInfo("Get storage pools on kvm hypervisor started.")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	listPoolOutput, err := k.RunKvmCommandAndGetStdout([]string{"pool-list"}, verbose)
	if err != nil {
		return nil, err
	}

	firstLine, unparsedOutput := astrings.SplitFirstLineAndContent(listPoolOutput)
	firstLine = strings.TrimSpace(firstLine)
	if !strings.HasPrefix(firstLine, "Name") {
		return nil, tracederrors.TracedErrorf("Unexpected first line of list pool output: '%s'", firstLine)
	}

	secondLine, unparsedOutput := astrings.SplitFirstLineAndContent(unparsedOutput)
	secondLine = strings.TrimSpace(secondLine)
	if strings.Count(secondLine, "-") < 5 {
		return nil, tracederrors.TracedErrorf("Unexpected second line of list pool output: '%s'", secondLine)
	}

	storagePools = []*KvmStoragePool{}
	for _, line := range astrings.SplitLines(unparsedOutput, true) {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		splitted := astrings.SplitAtSpacesAndRemoveEmptyStrings(line)
		if len(splitted) != 3 {
			return nil, tracederrors.TracedErrorf("Unable to splitt list pool line '%v' : '%v'", line, splitted)
		}

		nameToAdd := splitted[0]
		poolToAdd := NewKvmStoragePool()
		err = poolToAdd.SetName(nameToAdd)
		if err != nil {
			return nil, err
		}

		err = poolToAdd.SetHypervisor(k)
		if err != nil {
			return nil, err
		}

		storagePools = append(storagePools, poolToAdd)
	}

	if verbose {
		logging.LogInfof("Collected '%d' storage pools on kvm host '%s'", len(storagePools), hostname)
	}

	if verbose {
		logging.LogInfo("Get storage pools on kvm hypervisor finished.")
	}

	return storagePools, nil
}

func (k *KVMHypervisor) GetUseLocalhost() (useLocalhost bool, err error) {

	return k.useLocalhost, nil
}

func (k *KVMHypervisor) GetVmById(vmId int) (vm *KvmVm, err error) {
	vm = NewKvmVm()

	err = vm.SetHypervisor(k)
	if err != nil {
		return nil, err
	}

	err = vm.SetId(vmId)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func (k *KVMHypervisor) GetVmByName(vmName string, verbose bool) (vm *KvmVm, err error) {
	if vmName == "" {
		return nil, tracederrors.TracedError("vmName")
	}

	vms, err := k.GetVmList(verbose)
	if err != nil {
		return nil, err
	}

	for _, vm := range vms {
		nameToCheck, err := vm.GetCachedName()
		if err != nil {
			return nil, err
		}

		if nameToCheck == vmName {
			return vm, nil
		}
	}

	return nil, tracederrors.TracedErrorf("No VM named '%s' found", vmName)
}

func (k *KVMHypervisor) GetVmInfoList(verbose bool) (vmInfos []*KvmVmInfo, err error) {
	vms, err := k.GetVmList(verbose)
	if err != nil {
		return nil, err
	}

	vmInfos = []*KvmVmInfo{}
	for _, vm := range vms {
		infoToAdd, err := vm.GetInfo(verbose)
		if err != nil {
			return nil, err
		}

		vmInfos = append(vmInfos, infoToAdd)
	}

	return vmInfos, nil
}

func (k *KVMHypervisor) GetVmList(verbose bool) (vms []*KvmVm, err error) {
	listOutput, err := k.RunKvmCommandAndGetStdout([]string{"list", "--all"}, verbose)
	if err != nil {
		return nil, err
	}

	firstLine, unparsedOutput := astrings.SplitFirstLineAndContent(listOutput)
	firstLine = strings.TrimSpace(firstLine)
	if !strings.HasPrefix(firstLine, "Id ") {
		return nil, tracederrors.TracedErrorf("Unexpected first line '%s'. Full output is '%s'.", firstLine, listOutput)
	}

	secondLine, unparsedOutput := astrings.SplitFirstLineAndContent(unparsedOutput)
	if !strings.Contains(secondLine, "-----") {
		return nil, tracederrors.TracedErrorf("Unexpected second line '%s'. Full output is '%s'.", secondLine, listOutput)
	}

	vms = []*KvmVm{}
	for _, line := range astrings.SplitLines(unparsedOutput, true) {
		if len(strings.TrimSpace(line)) <= 0 {
			continue
		}

		lineToProcess := strings.ReplaceAll(line, "shut off", "shut_off")

		splitted := astrings.SplitAtSpacesAndRemoveEmptyStrings(lineToProcess)
		if len(splitted) != 3 {
			return nil, tracederrors.TracedErrorf("Failed to split line '%s'", line)
		}

		vmToAdd := NewKvmVm()
		err = vmToAdd.SetHypervisor(k)
		if err != nil {
			return nil, err
		}

		vmName := splitted[1]

		vmIdString := splitted[0]
		if vmIdString != "-" {
			vmId, err := strconv.Atoi(vmIdString)
			if err != nil {
				return nil, tracederrors.TracedErrorf("Unable to extract Vm id: '%s'", err.Error())
			}

			vmToAdd, err = k.GetVmById(vmId)
			if err != nil {
				return nil, err
			}
		}

		err = vmToAdd.SetCachedName(vmName)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vmToAdd)
	}

	if verbose {
		logging.LogInfof("Collected '%d' KVM Vms", len(vms))
	}

	return vms, nil
}

func (k *KVMHypervisor) GetVmNames(verbose bool) (vmNames []string, err error) {
	vms, err := k.GetVmList(verbose)
	if err != nil {
		return nil, err
	}

	vmNames = []string{}
	for _, vm := range vms {
		nameToAdd, err := vm.GetCachedName()
		if err != nil {
			return nil, err
		}

		vmNames = append(vmNames, nameToAdd)
	}

	return vmNames, nil
}

func (k *KVMHypervisor) GetVolumeByName(volumeName string) (volume *KvmVolume, err error) {
	if len(volumeName) <= 0 {
		return nil, tracederrors.TracedError("volumeName is empty string")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	const verboseVolumeCollection = false
	volumes, err := k.GetVolumes(verboseVolumeCollection)
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		nameToCheck, err := volume.GetName()
		if err != nil {
			return nil, err
		}

		if nameToCheck == volumeName {
			return volume, nil
		}
	}

	return nil, tracederrors.TracedErrorf("No volume '%s' found on hypervisor '%s'.", volumeName, hostname)
}

func (k *KVMHypervisor) GetVolumeNames(verbose bool) (volumeNames []string, err error) {
	volumes, err := k.GetVolumes(verbose)
	if err != nil {
		return nil, err
	}

	volumeNames = []string{}
	for _, volume := range volumes {
		nameToAdd, err := volume.GetName()
		if err != nil {
			return nil, err
		}

		volumeNames = append(volumeNames, nameToAdd)
	}

	return volumeNames, nil
}

func (k *KVMHypervisor) GetVolumes(verbose bool) (volumes []*KvmVolume, err error) {
	if verbose {
		logging.LogInfo("Get storage pools on kvm hypervisor started.")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	volumes = []*KvmVolume{}
	storagePools, err := k.GetStoragePools(verbose)
	if err != nil {
		return nil, err
	}

	for _, storagePool := range storagePools {
		volumesToAdd, err := storagePool.GetVolumes(verbose)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, volumesToAdd...)
	}

	if verbose {
		logging.LogInfof("Collected '%d' volumes from '%d' pools on kvm host '%s'", len(volumes), len(storagePools), hostname)
	}

	if verbose {
		logging.LogInfo("Get storage pools on kvm hypervisor finished.")
	}

	return volumes, nil
}

func (k *KVMHypervisor) MustCreateVm(createOptions *KvmCreateVmOptions) (createdVm *KvmVm) {
	createdVm, err := k.CreateVm(createOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return createdVm
}

func (k *KVMHypervisor) MustGetHost() (host hosts.Host) {
	host, err := k.GetHost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return host
}

func (k *KVMHypervisor) MustGetHostName() (hostname string) {
	hostname, err := k.GetHostName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return hostname
}

func (k *KVMHypervisor) MustGetStoragePoolNames(verbose bool) (storagePoolNames []string) {
	storagePoolNames, err := k.GetStoragePoolNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return storagePoolNames
}

func (k *KVMHypervisor) MustGetStoragePools(verbose bool) (storagePools []*KvmStoragePool) {
	storagePools, err := k.GetStoragePools(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return storagePools
}

func (k *KVMHypervisor) MustGetUseLocalhost() (useLocalhost bool) {
	useLocalhost, err := k.GetUseLocalhost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return useLocalhost
}

func (k *KVMHypervisor) MustGetVmById(vmId int) (vm *KvmVm) {
	vm, err := k.GetVmById(vmId)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vm
}

func (k *KVMHypervisor) MustGetVmByName(vmName string, verbose bool) (vm *KvmVm) {
	vm, err := k.GetVmByName(vmName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vm
}

func (k *KVMHypervisor) MustGetVmInfoList(verbose bool) (vmInfos []*KvmVmInfo) {
	vmInfos, err := k.GetVmInfoList(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vmInfos
}

func (k *KVMHypervisor) MustGetVmList(verbose bool) (vms []*KvmVm) {
	vms, err := k.GetVmList(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vms
}

func (k *KVMHypervisor) MustGetVmNames(verbose bool) (vmNames []string) {
	vmNames, err := k.GetVmNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vmNames
}

func (k *KVMHypervisor) MustGetVolumeByName(volumeName string) (volume *KvmVolume) {
	volume, err := k.GetVolumeByName(volumeName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return volume
}

func (k *KVMHypervisor) MustGetVolumeNames(verbose bool) (volumeNames []string) {
	volumeNames, err := k.GetVolumeNames(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return volumeNames
}

func (k *KVMHypervisor) MustGetVolumes(verbose bool) (volumes []*KvmVolume) {
	volumes, err := k.GetVolumes(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return volumes
}

func (k *KVMHypervisor) MustRemoveVm(removeOptions *KvmRemoveVmOptions) {
	err := k.RemoveVm(removeOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KVMHypervisor) MustRemoveVolumeByName(volumeName string, verbose bool) {
	err := k.RemoveVolumeByName(volumeName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KVMHypervisor) MustRunKvmCommand(kvmCommand []string, verbose bool) (commandOutput *asciichgolangpublic.CommandOutput) {
	commandOutput, err := k.RunKvmCommand(kvmCommand, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return commandOutput
}

func (k *KVMHypervisor) MustRunKvmCommandAndGetStdout(kvmCommand []string, verbose bool) (stdout string) {
	stdout, err := k.RunKvmCommandAndGetStdout(kvmCommand, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return stdout
}

func (k *KVMHypervisor) MustSetHost(host hosts.Host) {
	err := k.SetHost(host)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KVMHypervisor) MustSetUseLocalhost(useLocalhost bool) {
	err := k.SetUseLocalhost(useLocalhost)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KVMHypervisor) MustVmByNameExists(vmName string) (vmExists bool) {
	vmExists, err := k.VmByNameExists(vmName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return vmExists
}

func (k *KVMHypervisor) MustVolumeByNameExists(volumeName string) (volumeExists bool) {
	volumeExists, err := k.VolumeByNameExists(volumeName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return volumeExists
}

func (k *KVMHypervisor) RemoveVm(removeOptions *KvmRemoveVmOptions) (err error) {
	if removeOptions == nil {
		return tracederrors.TracedError("removeOptions is nil")
	}

	if len(removeOptions.VmName) <= 0 {
		return tracederrors.TracedError("vmName is empty string")
	}

	hostName, err := k.GetHostName()
	if err != nil {
		return err
	}

	if removeOptions.Verbose {
		logging.LogInfof("Going to delete kvm VM '%s' on host '%s'.", removeOptions.VmName, hostName)
	}

	vmExists, err := k.VmByNameExists(removeOptions.VmName)
	if err != nil {
		return
	}

	if vmExists {
		_, err = k.RunKvmCommandAndGetStdout([]string{"destroy", removeOptions.VmName}, removeOptions.Verbose)
		if err != nil {
			return err
		}

		vmExists, err = k.VmByNameExists(removeOptions.VmName)
		if err != nil {
			return
		}
		if vmExists {
			_, err = k.RunKvmCommandAndGetStdout([]string{"undefine", removeOptions.VmName}, removeOptions.Verbose)
			if err != nil {
				return err
			}
		}
		if removeOptions.Verbose {
			logging.LogChangedf("Vm '%s' removed on host '%s'.", removeOptions.VmName, hostName)
		}
	} else {
		if removeOptions.Verbose {
			logging.LogInfof("Vm '%s' is already removed on host '%s'.", removeOptions.VmName, hostName)
		}
	}

	if removeOptions.RemoveVolumes {
		for _, volumeName := range removeOptions.VolumeNamesToRemove {
			k.RemoveVolumeByName(volumeName, removeOptions.Verbose)
		}
	}

	return nil
}

func (k *KVMHypervisor) RemoveVolumeByName(volumeName string, verbose bool) (err error) {
	if len(volumeName) <= 0 {
		return tracederrors.TracedError("voluemName is empty string")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return err
	}

	volumeExists, err := k.VolumeByNameExists(volumeName)
	if err != nil {
		return err
	}

	if volumeExists {
		volume, err := k.GetVolumeByName(volumeName)
		if err != nil {
			return err
		}

		err = volume.Remove(verbose)
		if err != nil {
			return nil
		}

		logging.LogChangedf("Volume '%s' on KVM hypervisor '%s' deleted.", volumeName, hostname)
	} else {
		logging.LogInfof("Volume '%s' on KVM hypervisor '%s' was already deleted.", volumeName, hostname)
	}

	return nil
}

func (k *KVMHypervisor) RunKvmCommand(kvmCommand []string, verbose bool) (commandOutput *asciichgolangpublic.CommandOutput, err error) {
	if kvmCommand == nil {
		return nil, tracederrors.TracedError("kvmCommand is nil")
	}

	command := []string{"virsh", "-c", "qemu:///system"}
	command = append(command, kvmCommand...)

	if k.useLocalhost {
		commandOutput, err = asciichgolangpublic.Bash().RunCommand(
			&asciichgolangpublic.RunCommandOptions{
				Command: command,
				Verbose: verbose,
			},
		)
		if err != nil {
			return nil, err
		}
	} else {
		host, err := k.GetHost()
		if err != nil {
			return nil, err
		}

		commandOutput, err = host.RunCommand(
			&asciichgolangpublic.RunCommandOptions{
				Command: command,
				Verbose: verbose,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	if commandOutput == nil {
		return nil, tracederrors.TracedError("commandOutput is nil")
	}

	return commandOutput, nil
}

func (k *KVMHypervisor) RunKvmCommandAndGetStdout(kvmCommand []string, verbose bool) (stdout string, err error) {
	commandOutput, err := k.RunKvmCommand(kvmCommand, verbose)
	if err != nil {
		return "", err
	}

	stdout, err = commandOutput.GetStdoutAsString()
	if err != nil {
		return "", err
	}

	return stdout, nil
}

func (k *KVMHypervisor) SetHost(host hosts.Host) (err error) {
	if host == nil {
		return tracederrors.TracedError("nost is nil")
	}

	k.host = host

	return nil
}

func (k *KVMHypervisor) SetUseLocalhost(useLocalhost bool) (err error) {
	k.useLocalhost = useLocalhost
	return nil
}

func (k *KVMHypervisor) VmByNameExists(vmName string) (vmExists bool, err error) {
	if len(vmName) <= 0 {
		return false, tracederrors.TracedError("vmName is empty string")
	}

	const verbose = false
	vmNameList, err := k.GetVmNames(verbose)
	if err != nil {
		return false, err
	}

	if aslices.ContainsString(vmNameList, vmName) {
		return true, nil
	} else {
		return false, nil
	}
}

func (k *KVMHypervisor) VolumeByNameExists(volumeName string) (volumeExists bool, err error) {
	if len(volumeName) <= 0 {
		return false, tracederrors.TracedError("volumeName is empty string")
	}

	const verboseVolumeCollection = false
	volumes, err := k.GetVolumes(verboseVolumeCollection)
	if err != nil {
		return false, err
	}

	for _, volume := range volumes {
		nameToCheck, err := volume.GetName()
		if err != nil {
			return false, err
		}

		if nameToCheck == volumeName {
			return true, nil
		}
	}

	return false, nil
}
