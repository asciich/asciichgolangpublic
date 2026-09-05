package kvmutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmVm struct {
	vmId        *int
	cachedName  string
	cachedState string
	hypervisor  *KVMHypervisor
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

func (k *KvmVm) Delete(ctx context.Context) error {
	hypervisor, err := k.GetHypervisor()
	if err != nil {
		return err
	}

	name, err := k.GetCachedName()
	if err != nil {
		return err
	}

	return hypervisor.RemoveVm(
		ctx,
		&KvmRemoveVmOptions{
			VmName: name,
		},
	)
}

func (k *KvmVm) GetCachedState() (cachedState string, err error) {
	if len(k.cachedState) <= 0 {
		return "", tracederrors.TracedError("cachedState not set")
	}

	return k.cachedState, nil
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

	macAddress, err = GetMacAddressFromXmlString(domainXml)
	if err != nil {
		return "", err
	}

	return macAddress, nil
}

func (k *KvmVm) GetNetworkName(ctx context.Context) (networkName string, err error) {
	domainXml, err := k.GetDomainXmlAsString(ctx)
	if err != nil {
		return "", err
	}

	networkName, err = GetNetworkNameFromXmlString(domainXml)
	if err != nil {
		return "", err
	}

	return networkName, nil
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

func (k *KvmVm) IsRunning() (isRunning bool, err error) {
	cachedState, err := k.GetCachedState()
	if err != nil {
		return false, err
	}

	return cachedState == "running", nil
}

func (k *KvmVm) SetCachedName(cachedName string) (err error) {
	if len(cachedName) <= 0 {
		return tracederrors.TracedError("cachedName is empty string")
	}

	k.cachedName = cachedName

	return nil
}

func (k *KvmVm) SetCachedState(cachedState string) (err error) {
	if len(cachedState) <= 0 {
		return tracederrors.TracedError("cachedState is empty string")
	}

	k.cachedState = cachedState

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

func (k *KvmVm) GetIpAddress(ctx context.Context) (ipAddress string, err error) {
	vmName, err := k.GetCachedName()
	if err != nil {
		return "", err
	}

	// Try multiple sources in order so both NAT ('default') and bridged ('br0') VMs work:
	//   - agent: queries the qemu-guest-agent inside the VM (works for any network if the agent runs).
	//   - lease: reads libvirt's dnsmasq DHCP leases (works for the NAT 'default' network).
	//   - arp:   reads the host's ARP table (works for bridged setups if there is an ARP entry).
	for _, source := range []string{"agent", "lease", "arp"} {
		// Use silent context so the individual (expected to sometimes fail) lookups do not spam the log.
		ipAddress, err = k.getIpAddressBySource(contextutils.WithSilent(ctx), source)
		if err == nil && ipAddress != "" {
			return ipAddress, nil
		}
	}

	return "", tracederrors.TracedErrorf("No IPv4 address found for VM '%s' (tried sources agent, lease, arp).", vmName)
}

func (k *KvmVm) getIpAddressBySource(ctx context.Context, source string) (ipAddress string, err error) {
	if source == "" {
		return "", tracederrors.TracedErrorEmptyString("source")
	}

	hypervisor, err := k.GetHypervisor()
	if err != nil {
		return "", err
	}

	vmName, err := k.GetCachedName()
	if err != nil {
		return "", err
	}

	output, err := hypervisor.RunKvmCommandAndGetStdout(ctx, []string{"domifaddr", vmName, "--source", source})
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

func (k *KvmVm) GetVncPort(ctx context.Context) (vncPort int, err error) {
	domainXml, err := k.GetDomainXmlAsString(ctx)
	if err != nil {
		return -1, err
	}

	vncPort, err = GetVncPortFromXmlString(domainXml)
	if err != nil {
		return -1, err
	}

	return vncPort, nil
}
