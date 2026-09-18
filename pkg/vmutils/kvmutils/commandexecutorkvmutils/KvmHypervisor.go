package commandexecutorkvmutils

import (
	"context"
	"encoding/xml"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandoutput"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
)

type CommandExecutrKvmHypervisor struct {
	host hostsutilsinterfaces.Host

	// Run kvm commands and connection directly on localhost instead of using SSH.
	useLocalhost bool
}

func GetKvmHypervisorOnLocalhost() (kvmHypervisor *CommandExecutrKvmHypervisor, err error) {
	kvmHypervisor = NewKVMHypervisor()
	err = kvmHypervisor.SetUseLocalhost(true)
	if err != nil {
		return nil, err
	}

	return kvmHypervisor, nil
}

func NewKVMHypervisor() (kvmHypervisor *CommandExecutrKvmHypervisor) {
	return new(CommandExecutrKvmHypervisor)
}

func (k *CommandExecutrKvmHypervisor) GetHost() (host hostsutilsinterfaces.Host, err error) {
	if k.host == nil {
		return nil, tracederrors.TracedError("host not set")
	}

	return k.host, nil
}

func (k *CommandExecutrKvmHypervisor) GetHostName() (hostname string, err error) {
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

func (k *CommandExecutrKvmHypervisor) ListStoragePoolNames(ctx context.Context) (storagePoolNames []string, err error) {
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

func (k *CommandExecutrKvmHypervisor) GetStoragePoolByName(ctx context.Context, storagePoolName string) (kvmutilsinterfaces.StoragePool, error) {
	if storagePoolName == "" {
		return nil, tracederrors.TracedErrorEmptyString("storagePoolName")
	}

	list, err := k.ListStoragePools(ctx)
	if err != nil {
		return nil, err
	}

	var got kvmutilsinterfaces.StoragePool
	for _, p := range list {
		name, err := p.GetName()
		if err != nil {
			return nil, err
		}

		if name == storagePoolName {
			got = p
			break
		}
	}

	if got == nil {
		return nil, tracederrors.TracedErrorf("KVM storage pool '%s' not found.", storagePoolName)
	}

	return got, nil
}

func (k *CommandExecutrKvmHypervisor) ListStoragePools(ctx context.Context) (storagePools []*kvmutilsgeneric.KvmStoragePool, err error) {
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

	storagePools = []*kvmutilsgeneric.KvmStoragePool{}
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
		poolToAdd := kvmutilsgeneric.NewKvmStoragePool()
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

func (k *CommandExecutrKvmHypervisor) GetUseLocalhost() (useLocalhost bool, err error) {

	return k.useLocalhost, nil
}

func (k *CommandExecutrKvmHypervisor) GetVmById(vmId int) (vm kvmutilsinterfaces.VM, err error) {
	vm = kvmutilsgeneric.NewKvmVm()

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

func (k *CommandExecutrKvmHypervisor) GetVmByName(ctx context.Context, vmName string) (vm kvmutilsinterfaces.VM, err error) {
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

func (k *CommandExecutrKvmHypervisor) ListVmInfos(ctx context.Context) (vmInfos []kvmutilsinterfaces.VmInfo, err error) {
	vms, err := k.ListVms(ctx)
	if err != nil {
		return nil, err
	}

	vmInfos = []kvmutilsinterfaces.VmInfo{}
	for _, vm := range vms {
		infoToAdd, err := vm.GetInfo(ctx)
		if err != nil {
			return nil, err
		}

		vmInfos = append(vmInfos, infoToAdd)
	}

	return vmInfos, nil
}

func (k *CommandExecutrKvmHypervisor) ListVms(ctx context.Context) (vms []kvmutilsinterfaces.VM, err error) {
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

	vms = []kvmutilsinterfaces.VM{}
	for _, line := range stringsutils.SplitLines(unparsedOutput, true) {
		if len(strings.TrimSpace(line)) <= 0 {
			continue
		}

		lineToProcess := strings.ReplaceAll(line, "shut off", "shut_off")

		splitted := stringsutils.SplitAtSpacesAndRemoveEmptyStrings(lineToProcess)
		if len(splitted) != 3 {
			return nil, tracederrors.TracedErrorf("Failed to split line '%s'", line)
		}

		var vmToAdd kvmutilsinterfaces.VM = kvmutilsgeneric.NewKvmVm()
		err = vmToAdd.SetHypervisor(k)
		if err != nil {
			return nil, err
		}

		vmName := splitted[1]

		// Undo the "shut off" -> "shut_off" normalization to store the real libvirt state.
		vmState := strings.ReplaceAll(splitted[2], "shut_off", "shut off")

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

		err = vmToAdd.SetCachedState(vmState)
		if err != nil {
			return nil, err
		}

		vms = append(vms, vmToAdd)
	}

	logging.LogInfoByCtxf(ctx, "Collected '%d' KVM Vms", len(vms))

	return vms, nil
}

func (k *CommandExecutrKvmHypervisor) ListVmNames(ctx context.Context) (vmNames []string, err error) {
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

func (k *CommandExecutrKvmHypervisor) GetVolumeByName(ctx context.Context, volumeName string) (volume kvmutilsinterfaces.Volume, err error) {
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

func (k *CommandExecutrKvmHypervisor) ListVolumeNames(ctx context.Context) (volumeNames []string, err error) {
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

func (k *CommandExecutrKvmHypervisor) GetVolumes(ctx context.Context) (volumes []kvmutilsinterfaces.Volume, err error) {
	logging.LogInfoByCtxf(ctx, "Get storage pools on kvm hypervisor started.")

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	volumes = []kvmutilsinterfaces.Volume{}
	storagePools, err := k.ListStoragePools(ctx)
	if err != nil {
		return nil, err
	}

	for _, storagePool := range storagePools {
		volumesToAdd, err := storagePool.ListVolumes(ctx)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, volumesToAdd...)
	}

	logging.LogInfoByCtxf(ctx, "Collected '%d' volumes from '%d' pools on kvm host '%s'", len(volumes), len(storagePools), hostname)

	logging.LogInfoByCtxf(ctx, "Get storage pools on kvm hypervisor finished.")

	return volumes, nil
}

func (k *CommandExecutrKvmHypervisor) DeleteVm(ctx context.Context, removeOptions *kvmutilsoptions.KvmRemoveVmOptions) (err error) {
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

func (k *CommandExecutrKvmHypervisor) RemoveVolumeByName(ctx context.Context, volumeName string) (err error) {
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

		err = volume.Delete(ctx)
		if err != nil {
			return err
		}

		logging.LogChangedf("Volume '%s' on KVM hypervisor '%s' deleted.", volumeName, hostname)
	} else {
		logging.LogInfof("Volume '%s' on KVM hypervisor '%s' was already deleted.", volumeName, hostname)
	}

	return nil
}

func (k *CommandExecutrKvmHypervisor) RunKvmCommand(ctx context.Context, kvmCommand []string) (commandOutput *commandoutput.CommandOutput, err error) {
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

func (k *CommandExecutrKvmHypervisor) RunKvmCommandAndGetStdout(ctx context.Context, kvmCommand []string) (stdout string, err error) {
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

func (k *CommandExecutrKvmHypervisor) SetHost(host hostsutilsinterfaces.Host) (err error) {
	if host == nil {
		return tracederrors.TracedError("nost is nil")
	}

	k.host = host

	return nil
}

func (k *CommandExecutrKvmHypervisor) SetUseLocalhost(useLocalhost bool) (err error) {
	k.useLocalhost = useLocalhost
	return nil
}

func (k *CommandExecutrKvmHypervisor) VmByNameExists(ctx context.Context, vmName string) (vmExists bool, err error) {
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

func (k *CommandExecutrKvmHypervisor) VolumeByNameExists(ctx context.Context, volumeName string) (volumeExists bool, err error) {
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

func (k *CommandExecutrKvmHypervisor) ListNetworkNames(ctx context.Context) (networkNames []string, err error) {
	networks, err := k.ListNetworks(ctx)
	if err != nil {
		return nil, err
	}

	networkNames = []string{}
	for _, network := range networks {
		nameToAdd, err := network.GetName()
		if err != nil {
			return nil, err
		}

		networkNames = append(networkNames, nameToAdd)
	}

	return networkNames, nil
}

func (k *CommandExecutrKvmHypervisor) ListNetworks(ctx context.Context) (networks []kvmutilsinterfaces.Network, err error) {
	logging.LogInfoByCtxf(ctx, "List networks on kvm hypervisor started.")

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	// --all also lists inactive networks (e.g. an inactive 'default' network).
	listNetworkOutput, err := k.RunKvmCommandAndGetStdout(ctx, []string{"net-list", "--all"})
	if err != nil {
		return nil, err
	}

	firstLine, unparsedOutput := stringsutils.SplitFirstLineAndContent(listNetworkOutput)
	firstLine = strings.TrimSpace(firstLine)
	if !strings.HasPrefix(firstLine, "Name") {
		return nil, tracederrors.TracedErrorf("Unexpected first line of list network output: '%s'", firstLine)
	}

	secondLine, unparsedOutput := stringsutils.SplitFirstLineAndContent(unparsedOutput)
	secondLine = strings.TrimSpace(secondLine)
	if strings.Count(secondLine, "-") < 5 {
		return nil, tracederrors.TracedErrorf("Unexpected second line of list network output: '%s'", secondLine)
	}

	networks = []kvmutilsinterfaces.Network{}
	for _, line := range stringsutils.SplitLines(unparsedOutput, true) {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		splitted := stringsutils.SplitAtSpacesAndRemoveEmptyStrings(line)
		if len(splitted) != 4 {
			return nil, tracederrors.TracedErrorf("Unable to split list network line '%v' : '%v'", line, splitted)
		}

		networkToAdd := kvmutilsgeneric.NewNetwork()

		err = networkToAdd.SetHypervisor(k)
		if err != nil {
			return nil, err
		}

		err = networkToAdd.SetName(splitted[0])
		if err != nil {
			return nil, err
		}

		err = networkToAdd.SetState(splitted[1])
		if err != nil {
			return nil, err
		}

		err = networkToAdd.SetAutostart(splitted[2])
		if err != nil {
			return nil, err
		}

		err = networkToAdd.SetPersistent(splitted[3])
		if err != nil {
			return nil, err
		}

		networks = append(networks, networkToAdd)
	}

	logging.LogInfoByCtxf(ctx, "Collected '%d' networks on kvm host '%s'", len(networks), hostname)

	logging.LogInfoByCtxf(ctx, "List networks on kvm hypervisor finished.")

	return networks, nil
}

func (k *CommandExecutrKvmHypervisor) NetworkByNameExists(ctx context.Context, networkName string) (networkExists bool, err error) {
	if networkName == "" {
		return false, tracederrors.TracedErrorEmptyString("networkName")
	}

	networkNames, err := k.ListNetworkNames(contextutils.WithSilent(ctx))
	if err != nil {
		return false, err
	}

	return slices.Contains(networkNames, networkName), nil
}

func (k *CommandExecutrKvmHypervisor) StartNetworkByName(ctx context.Context, networkName string) (err error) {
	if networkName == "" {
		return tracederrors.TracedErrorEmptyString("networkName")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return err
	}

	networks, err := k.ListNetworks(contextutils.WithSilent(ctx))
	if err != nil {
		return err
	}

	for _, network := range networks {
		nameToCheck, err := network.GetName()
		if err != nil {
			return err
		}

		if nameToCheck != networkName {
			continue
		}

		isActive, err := network.IsActive()
		if err != nil {
			return err
		}

		if isActive {
			logging.LogInfoByCtxf(ctx, "Network '%s' is already active on kvm host '%s'.", networkName, hostname)
			return nil
		}

		_, err = k.RunKvmCommandAndGetStdout(ctx, []string{"net-start", networkName})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Network '%s' started on kvm host '%s'.", networkName, hostname)
		return nil
	}

	return tracederrors.TracedErrorf("No network named '%s' found on kvm host '%s'.", networkName, hostname)
}

func (k *CommandExecutrKvmHypervisor) ResetVm(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	hostName, err := k.GetHostName()
	if err != nil {
		return err
	}

	vm, err := k.GetVmByName(ctx, name)
	if err != nil {
		return err
	}

	isRunning, err := vm.IsRunning(ctx)
	if err != nil {
		return err
	}

	if !isRunning {
		logging.LogInfoByCtxf(ctx, "Vm '%s' is not running on host '%s'. Skip reset.", name, hostName)
		return nil
	}

	_, err = k.RunKvmCommandAndGetStdout(ctx, []string{"reset", name})
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Vm '%s' reset on host '%s'.", name, hostName)

	return nil
}

func (k *CommandExecutrKvmHypervisor) GetNetworkByName(networkName string) (kvmutilsinterfaces.Network, error) {
	if networkName == "" {
		return nil, tracederrors.TracedErrorEmptyString("networkName")
	}

	network := kvmutilsgeneric.NewNetwork()
	err := network.SetName(networkName)
	if err != nil {
		return nil, err
	}

	err = network.SetHypervisor(k)
	if err != nil {
		return nil, err
	}

	return network, nil
}

func (k *CommandExecutrKvmHypervisor) GetParsedNetworkXml(ctx context.Context, networkName string) (kvmutilsinterfaces.KvmNetworkXml, error) {
	if networkName == "" {
		return nil, tracederrors.TracedErrorEmptyString("networkName")
	}

	stdout, err := k.RunKvmCommandAndGetStdout(ctx, []string{"net-dumpxml", networkName})
	if err != nil {
		return nil, err
	}

	parsed := &kvmutilsgeneric.KvmNetworkXml{}
	err = xml.Unmarshal([]byte(stdout), parsed)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse net-dumpxml output for network '%s': %w", networkName, err)
	}

	return parsed, nil
}

func (k *CommandExecutrKvmHypervisor) DeleteIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string) error {
	if networkName == "" {
		return tracederrors.TracedErrorEmptyString("networkName")
	}

	if startIp == "" {
		return tracederrors.TracedErrorEmptyString("startIp")
	}

	if endIp == "" {
		return tracederrors.TracedErrorEmptyString("endIp")
	}

	xml := fmt.Sprintf("<range start='%s' end='%s'/>", startIp, endIp)
	_, err := k.RunKvmCommandAndGetStdout(
		ctx,
		[]string{"net-update", networkName, "delete", "ip-dhcp-range", xml, "--live", "--config"},
	)
	if err != nil {
		return err
	}

	return nil
}

func (k *CommandExecutrKvmHypervisor) AddIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string) error {
	if networkName == "" {
		return tracederrors.TracedErrorEmptyString("networkName")
	}

	if startIp == "" {
		return tracederrors.TracedErrorEmptyString("startIp")
	}

	if endIp == "" {
		return tracederrors.TracedErrorEmptyString("endIp")
	}

	xml := fmt.Sprintf("<range start='%s' end='%s'/>", startIp, endIp)
	_, err := k.RunKvmCommandAndGetStdout(
		ctx,
		[]string{"net-update", networkName, "add", "ip-dhcp-range", xml, "--live", "--config"},
	)
	if err != nil {
		return err
	}

	return nil
}

func (k *CommandExecutrKvmHypervisor) GetHostDescription() (string, error) {
	host, err := k.GetHost()
	if err != nil {
		return "", err
	}

	return host.GetHostDescription()
}

func (k *CommandExecutrKvmHypervisor) DeleteVolumeByName(ctx context.Context, storagePoolName string, volumeName string) error {
	if storagePoolName == "" {
		return tracederrors.TracedErrorEmptyString(storagePoolName)
	}

	if volumeName == "" {
		return tracederrors.TracedErrorEmptyString(volumeName)
	}

	hostDescription, err := k.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Remove KVM volume '%s' of storage pool '%s' on hypervisor '%s' started.", volumeName, storagePoolName, hostDescription)

	_, err = k.RunKvmCommand(ctx, []string{"vol-delete", "--pool", storagePoolName, volumeName})
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Remove KVM volume '%s' of storage pool '%s' on hypervisor '%s' deleted.", volumeName, storagePoolName, hostDescription)

	logging.LogInfoByCtxf(ctx, "Remove KVM volume '%s' of storage pool '%s' on hypervisor '%s' finished.", volumeName, storagePoolName, hostDescription)

	return nil
}

func (k *CommandExecutrKvmHypervisor) ListVolumes(ctx context.Context, storagePoolName string) ([]kvmutilsinterfaces.Volume, error) {
	if storagePoolName == "" {
		return nil, tracederrors.TracedErrorEmptyString("storagePoolName")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Get volumes in storage pool '%s' on kvm hypervisor '%s' started.", storagePoolName, hostname)

	listPoolOutput, err := k.RunKvmCommandAndGetStdout(ctx, []string{"vol-list", "--pool", storagePoolName})
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

	storagePool, err := k.GetStoragePoolByName(ctx, storagePoolName)
	if err != nil {
		return nil, err
	}

	volumes := []kvmutilsinterfaces.Volume{}
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
		volumeToAdd := kvmutilsgeneric.NewKvmVolume()
		err = volumeToAdd.SetName(nameToAdd)
		if err != nil {
			return nil, err
		}

		err = volumeToAdd.SetStoragePool(storagePool)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, volumeToAdd)
	}

	logging.LogInfoByCtxf(ctx, "Collected '%d' storage pools on kvm host '%s'", len(volumes), hostname)

	logging.LogInfoByCtxf(ctx, "Get volumes in storage pool '%s' on kvm hypervisor '%s' finished.", storagePoolName, hostname)

	return volumes, nil
}

func (k *CommandExecutrKvmHypervisor) GetDomainXmlAsString(ctx context.Context, vmName string) (domainXml string, err error) {
	if vmName == "" {
		return "", tracederrors.TracedErrorEmptyString("vmName")
	}

	domainXml, err = k.RunKvmCommandAndGetStdout(ctx, []string{"dumpxml", vmName})
	if err != nil {
		return "", err
	}

	return domainXml, nil
}

func (k *CommandExecutrKvmHypervisor) GetIpAddress(ctx context.Context, vmName string) (string, error) {
	if vmName == "" {
		return "", tracederrors.TracedErrorEmptyString("vmName")
	}

	// Try multiple sources in order so both NAT ('default') and bridged ('br0') VMs work:
	//   - agent: queries the qemu-guest-agent inside the VM (works for any network if the agent runs).
	//   - lease: reads libvirt's dnsmasq DHCP leases (works for the NAT 'default' network).
	//   - arp:   reads the host's ARP table (works for bridged setups if there is an ARP entry).
	for _, source := range []string{"agent", "lease", "arp"} {
		// Use silent context so the individual (expected to sometimes fail) lookups do not spam the log.
		ipAddress, err := k.getIpAddressBySource(contextutils.WithSilent(ctx), vmName, source)
		if err == nil && ipAddress != "" {
			return ipAddress, nil
		}
	}

	return "", tracederrors.TracedErrorf("No IPv4 address found for VM '%s' (tried sources agent, lease, arp).", vmName)
}

func (k *CommandExecutrKvmHypervisor) getIpAddressBySource(ctx context.Context, vmName string, source string) (ipAddress string, err error) {
	if vmName == "" {
		return "", tracederrors.TracedErrorEmptyString("vmName")
	}

	if source == "" {
		return "", tracederrors.TracedErrorEmptyString("source")
	}

	output, err := k.RunKvmCommandAndGetStdout(ctx, []string{"domifaddr", vmName, "--source", source})
	if err != nil {
		return "", err
	}

	for _, line := range stringsutils.SplitLines(output, true) {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip header and separator lines.
		if strings.HasPrefix(line, "Name") {
			continue
		}
		if strings.Count(line, "-") > 5 {
			continue
		}

		splitted := stringsutils.SplitAtSpacesAndRemoveEmptyStrings(line)
		if len(splitted) != 4 {
			continue
		}

		if splitted[2] != "ipv4" {
			continue
		}

		// splitted[3] is like "192.168.122.94/24" -> strip the CIDR suffix.
		ipAddress = strings.SplitN(splitted[3], "/", 2)[0]

		return ipAddress, nil
	}

	return "", tracederrors.TracedErrorf("No IPv4 address found for VM '%s' via source '%s'.", vmName, source)
}
