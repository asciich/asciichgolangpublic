package kvmutils

import (
	"context"
	"slices"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/hosts"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KVMHypervisor struct {
	host hosts.Host

	// Run kvm commands and connection directly on localhost instead of using SSH.
	useLocalhost bool
}

func GetKvmHypervisorByHostName(hostname string) (kvmHypervisor *KVMHypervisor, err error) {
	if hostname == "" {
		return nil, tracederrors.TracedErrorEmptyString("hostname")
	}

	if hostname == "localhost" {
		return GetKvmHypervisorOnLocalhost()
	}

	host, err := hosts.GetHostByHostname(hostname)
	if err != nil {
		return nil, err
	}

	return GetKvmHypervisorByHost(host)
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

func NewKVMHypervisor() (kvmHypervisor *KVMHypervisor) {
	return new(KVMHypervisor)
}

func (k *KVMHypervisor) CreateVm(ctx context.Context, createOptions *KvmCreateVmOptions) (createdVm *KvmVm, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
	}

	vmName, err := createOptions.GetVmName()
	if err != nil {
		return nil, err
	}

	exists, err := k.VmByNameExists(ctx, vmName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "VM '%s' already exists", vmName)

		createdVm, err = k.GetVmByName(ctx, vmName)
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

	diskImageExists, err := diskImage.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !diskImageExists {
		return nil, tracederrors.TracedErrorf("Disk image '%s' does not exist to create VM.", diskImagePath)
	}

	vmXml, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
	if err != nil {
		return nil, err
	}
	defer vmXml.Delete(ctx, &filesoptions.DeleteOptions{})

	err = LibvirtXmls().WriteXmlForVmOnLatopToFile(ctx, createOptions, vmXml)
	if err != nil {
		return nil, err
	}

	vmXmlPath, err := vmXml.GetLocalPath()
	if err != nil {
		return nil, err
	}

	createOutput, err := k.RunKvmCommandAndGetStdout(
		ctx,
		[]string{"create", vmXmlPath},
	)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Output of VM '%s' creation:\n%s", vmName, createOutput)

	createdVm, err = k.GetVmByName(ctx, vmName)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Vm '%s' created.", vmName)

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

func (k *KVMHypervisor) ListStoragePoolNames(ctx context.Context) (storagePoolNames []string, err error) {
	storagePools, err := k.ListStoragePools(ctx)
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

func (k *KVMHypervisor) ListStoragePools(ctx context.Context) (storagePools []*KvmStoragePool, err error) {
	logging.LogInfoByCtxf(ctx, "Get storage pools on kvm hypervisor started.")

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	listPoolOutput, err := k.RunKvmCommandAndGetStdout(ctx, []string{"pool-list"})
	if err != nil {
		return nil, err
	}

	firstLine, unparsedOutput := stringsutils.SplitFirstLineAndContent(listPoolOutput)
	firstLine = strings.TrimSpace(firstLine)
	if !strings.HasPrefix(firstLine, "Name") {
		return nil, tracederrors.TracedErrorf("Unexpected first line of list pool output: '%s'", firstLine)
	}

	secondLine, unparsedOutput := stringsutils.SplitFirstLineAndContent(unparsedOutput)
	secondLine = strings.TrimSpace(secondLine)
	if strings.Count(secondLine, "-") < 5 {
		return nil, tracederrors.TracedErrorf("Unexpected second line of list pool output: '%s'", secondLine)
	}

	storagePools = []*KvmStoragePool{}
	for _, line := range stringsutils.SplitLines(unparsedOutput, true) {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		splitted := stringsutils.SplitAtSpacesAndRemoveEmptyStrings(line)
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

	logging.LogInfoByCtxf(ctx, "Collected '%d' storage pools on kvm host '%s'", len(storagePools), hostname)

	logging.LogInfoByCtxf(ctx, "Get storage pools on kvm hypervisor finished.")

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

func (k *KVMHypervisor) GetVmByName(ctx context.Context, vmName string) (vm *KvmVm, err error) {
	if vmName == "" {
		return nil, tracederrors.TracedError("vmName")
	}

	vms, err := k.ListVms(ctx)
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

func (k *KVMHypervisor) GetVmInfoList(ctx context.Context) (vmInfos []*KvmVmInfo, err error) {
	vms, err := k.ListVms(ctx)
	if err != nil {
		return nil, err
	}

	vmInfos = []*KvmVmInfo{}
	for _, vm := range vms {
		infoToAdd, err := vm.GetInfo(ctx)
		if err != nil {
			return nil, err
		}

		vmInfos = append(vmInfos, infoToAdd)
	}

	return vmInfos, nil
}

func (k *KVMHypervisor) ListVms(ctx context.Context) (vms []*KvmVm, err error) {
	listOutput, err := k.RunKvmCommandAndGetStdout(ctx, []string{"list", "--all"})
	if err != nil {
		return nil, err
	}

	firstLine, unparsedOutput := stringsutils.SplitFirstLineAndContent(listOutput)
	firstLine = strings.TrimSpace(firstLine)
	if !strings.HasPrefix(firstLine, "Id ") {
		return nil, tracederrors.TracedErrorf("Unexpected first line '%s'. Full output is '%s'.", firstLine, listOutput)
	}

	secondLine, unparsedOutput := stringsutils.SplitFirstLineAndContent(unparsedOutput)
	if !strings.Contains(secondLine, "-----") {
		return nil, tracederrors.TracedErrorf("Unexpected second line '%s'. Full output is '%s'.", secondLine, listOutput)
	}

	vms = []*KvmVm{}
	for _, line := range stringsutils.SplitLines(unparsedOutput, true) {
		if len(strings.TrimSpace(line)) <= 0 {
			continue
		}

		lineToProcess := strings.ReplaceAll(line, "shut off", "shut_off")

		splitted := stringsutils.SplitAtSpacesAndRemoveEmptyStrings(lineToProcess)
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

	logging.LogInfoByCtxf(ctx, "Collected '%d' KVM Vms", len(vms))

	return vms, nil
}

func (k *KVMHypervisor) ListVmNames(ctx context.Context) (vmNames []string, err error) {
	vms, err := k.ListVms(ctx)
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

func (k *KVMHypervisor) GetVolumeByName(ctx context.Context, volumeName string) (volume *KvmVolume, err error) {
	if len(volumeName) <= 0 {
		return nil, tracederrors.TracedError("volumeName is empty string")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	volumes, err := k.GetVolumes(contextutils.WithSilent(ctx))
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

func (k *KVMHypervisor) GetVolumeNames(ctx context.Context) (volumeNames []string, err error) {
	volumes, err := k.GetVolumes(ctx)
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

func (k *KVMHypervisor) GetVolumes(ctx context.Context) (volumes []*KvmVolume, err error) {
	logging.LogInfoByCtxf(ctx, "Get storage pools on kvm hypervisor started.")

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	volumes = []*KvmVolume{}
	storagePools, err := k.ListStoragePools(ctx)
	if err != nil {
		return nil, err
	}

	for _, storagePool := range storagePools {
		volumesToAdd, err := storagePool.GetVolumes(ctx)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, volumesToAdd...)
	}

	logging.LogInfoByCtxf(ctx, "Collected '%d' volumes from '%d' pools on kvm host '%s'", len(volumes), len(storagePools), hostname)

	logging.LogInfoByCtxf(ctx, "Get storage pools on kvm hypervisor finished.")

	return volumes, nil
}

func (k *KVMHypervisor) RemoveVm(ctx context.Context, removeOptions *KvmRemoveVmOptions) (err error) {
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

	logging.LogInfoByCtxf(ctx, "Going to delete kvm VM '%s' on host '%s'.", removeOptions.VmName, hostName)

	vmExists, err := k.VmByNameExists(ctx, removeOptions.VmName)
	if err != nil {
		return
	}

	if vmExists {
		_, err = k.RunKvmCommandAndGetStdout(ctx, []string{"destroy", removeOptions.VmName})
		if err != nil {
			return err
		}

		vmExists, err = k.VmByNameExists(ctx, removeOptions.VmName)
		if err != nil {
			return
		}
		if vmExists {
			_, err = k.RunKvmCommandAndGetStdout(ctx, []string{"undefine", removeOptions.VmName})
			if err != nil {
				return err
			}
		}
		logging.LogChangedByCtxf(ctx, "Vm '%s' removed on host '%s'.", removeOptions.VmName, hostName)
	} else {
		logging.LogInfoByCtxf(ctx, "Vm '%s' is already removed on host '%s'.", removeOptions.VmName, hostName)
	}

	if removeOptions.RemoveVolumes {
		for _, volumeName := range removeOptions.VolumeNamesToRemove {
			k.RemoveVolumeByName(ctx, volumeName)
		}
	}

	return nil
}

func (k *KVMHypervisor) RemoveVolumeByName(ctx context.Context, volumeName string) (err error) {
	if len(volumeName) <= 0 {
		return tracederrors.TracedError("voluemName is empty string")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return err
	}

	volumeExists, err := k.VolumeByNameExists(ctx, volumeName)
	if err != nil {
		return err
	}

	if volumeExists {
		volume, err := k.GetVolumeByName(ctx, volumeName)
		if err != nil {
			return err
		}

		err = volume.Remove(ctx)
		if err != nil {
			return nil
		}

		logging.LogChangedf("Volume '%s' on KVM hypervisor '%s' deleted.", volumeName, hostname)
	} else {
		logging.LogInfof("Volume '%s' on KVM hypervisor '%s' was already deleted.", volumeName, hostname)
	}

	return nil
}

func (k *KVMHypervisor) RunKvmCommand(ctx context.Context, kvmCommand []string) (commandOutput *commandoutput.CommandOutput, err error) {
	if kvmCommand == nil {
		return nil, tracederrors.TracedError("kvmCommand is nil")
	}

	command := []string{"virsh", "-c", "qemu:///system"}
	command = append(command, kvmCommand...)

	if k.useLocalhost {
		commandOutput, err = commandexecutorbashoo.Bash().RunCommand(
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
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
			ctx,
			&parameteroptions.RunCommandOptions{
				Command: command,
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

func (k *KVMHypervisor) RunKvmCommandAndGetStdout(ctx context.Context, kvmCommand []string) (stdout string, err error) {
	commandOutput, err := k.RunKvmCommand(ctx, kvmCommand)
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

func (k *KVMHypervisor) VmByNameExists(ctx context.Context, vmName string) (vmExists bool, err error) {
	if len(vmName) <= 0 {
		return false, tracederrors.TracedError("vmName is empty string")
	}

	const verbose = false
	vmNameList, err := k.ListVmNames(ctx)
	if err != nil {
		return false, err
	}

	if slices.Contains(vmNameList, vmName) {
		return true, nil
	} else {
		return false, nil
	}
}

func (k *KVMHypervisor) VolumeByNameExists(ctx context.Context, volumeName string) (volumeExists bool, err error) {
	if len(volumeName) <= 0 {
		return false, tracederrors.TracedError("volumeName is empty string")
	}

	volumes, err := k.GetVolumes(contextutils.WithSilent(ctx))
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
