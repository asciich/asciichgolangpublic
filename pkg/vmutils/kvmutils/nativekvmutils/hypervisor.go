package nativekvmutils

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/digitalocean/go-libvirt"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/iputils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
)

// NativeKvmHypervisor implements a KVM/libvirt hypervisor using the pure-Go
// libvirt RPC client (github.com/digitalocean/go-libvirt). It does NOT shell
// out to the 'virsh' binary.
//
// As per the repository constitution the native implementation is intended for
// local execution (localhost). Remote execution over SSH is nevertheless
// supported by libvirt's own 'qemu+ssh://' transport when useLocalhost is set
// to false and a host is configured.
type NativeKvmHypervisor struct {
	host hostsutilsinterfaces.Host

	// Connect directly to the local libvirt daemon instead of connecting via a
	// remote (qemu+ssh) transport.
	useLocalhost bool

	// Cached libvirt connection. Access is guarded by connectionMutex.
	connectionMutex   sync.Mutex
	libvirtConnection *libvirt.Libvirt
}

// GetKvmHypervisorOnLocalhost returns a NativeKvmHypervisor connecting to the
// local libvirt daemon (qemu:///system).
func GetKvmHypervisorOnLocalhost() (kvmHypervisor *NativeKvmHypervisor, err error) {
	kvmHypervisor = NewKvmHypervisor()

	err = kvmHypervisor.SetUseLocalhost(true)
	if err != nil {
		return nil, err
	}

	return kvmHypervisor, nil
}

// NewKvmHypervisor returns an empty NativeKvmHypervisor.
func NewKvmHypervisor() (kvmHypervisor *NativeKvmHypervisor) {
	return new(NativeKvmHypervisor)
}

// ---------------------------------------------------------------------------
// Connection handling
// ---------------------------------------------------------------------------

func (k *NativeKvmHypervisor) getConnectionURI() (uri *url.URL, err error) {
	if k.useLocalhost {
		uri, err = url.Parse(string(libvirt.QEMUSystem))
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to parse local libvirt URI: %w", err)
		}

		return uri, nil
	}

	host, err := k.GetHost()
	if err != nil {
		return nil, err
	}

	hostname, err := host.GetHostName()
	if err != nil {
		return nil, err
	}

	if hostname == "" {
		return nil, tracederrors.TracedError("hostname is empty string")
	}

	// Use libvirt's own remote transport so the native implementation still
	// works against remote hypervisors without shelling out to virsh.
	uri, err = url.Parse(fmt.Sprintf("qemu+ssh://%s/system", hostname))
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to parse remote libvirt URI for host '%s': %w", hostname, err)
	}

	return uri, nil
}

// getLibvirtConnection returns a connected libvirt client, establishing the
// connection lazily on first use.
func (k *NativeKvmHypervisor) getLibvirtConnection(ctx context.Context) (connection *libvirt.Libvirt, err error) {
	k.connectionMutex.Lock()
	defer k.connectionMutex.Unlock()

	if k.libvirtConnection != nil && k.libvirtConnection.IsConnected() {
		return k.libvirtConnection, nil
	}

	uri, err := k.getConnectionURI()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Connecting to libvirt using URI '%s'.", uri.String())

	connection, err = libvirt.ConnectToURI(uri)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to connect to libvirt using URI '%s': %w", uri.String(), err)
	}

	k.libvirtConnection = connection

	return connection, nil
}

// Close closes the underlying libvirt connection if it is open. It is safe to
// call multiple times and is therefore suitable to be used with 'defer'.
func (k *NativeKvmHypervisor) Close(ctx context.Context) (err error) {
	k.connectionMutex.Lock()
	defer k.connectionMutex.Unlock()

	if k.libvirtConnection == nil {
		return nil
	}

	if k.libvirtConnection.IsConnected() {
		err = k.libvirtConnection.Disconnect()
		if err != nil {
			return tracederrors.TracedErrorf("failed to disconnect from libvirt: %w", err)
		}

		logging.LogInfoByCtxf(ctx, "Disconnected from libvirt.")
	}

	k.libvirtConnection = nil

	return nil
}

// ---------------------------------------------------------------------------
// Host / connection meta data
// ---------------------------------------------------------------------------

func (k *NativeKvmHypervisor) GetHost() (host hostsutilsinterfaces.Host, err error) {
	if k.host == nil {
		return nil, tracederrors.TracedError("host not set")
	}

	return k.host, nil
}

func (k *NativeKvmHypervisor) SetHost(host hostsutilsinterfaces.Host) (err error) {
	if host == nil {
		return tracederrors.TracedError("host is nil")
	}

	k.host = host

	return nil
}

func (k *NativeKvmHypervisor) GetHostName() (hostname string, err error) {
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

func (k *NativeKvmHypervisor) GetHostDescription() (hostDescription string, err error) {
	if k.useLocalhost {
		return "localhost_connection", nil
	}

	host, err := k.GetHost()
	if err != nil {
		return "", err
	}

	return host.GetHostDescription()
}

func (k *NativeKvmHypervisor) GetUseLocalhost() (useLocalhost bool, err error) {
	return k.useLocalhost, nil
}

func (k *NativeKvmHypervisor) SetUseLocalhost(useLocalhost bool) (err error) {
	k.useLocalhost = useLocalhost

	return nil
}

// ---------------------------------------------------------------------------
// VMs / domains
// ---------------------------------------------------------------------------

func (k *NativeKvmHypervisor) CreateVm(ctx context.Context, createOptions *kvmutilsoptions.KvmCreateVmOptions) (createdVm kvmutilsinterfaces.VM, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedError("createOptions is nil")
	}

	vmName, err := createOptions.GetVmName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Create KVM VM '%s' started.", vmName)

	exists, err := k.VmByNameExists(ctx, vmName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "VM '%s' already exists.", vmName)

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

	if diskImageExists {
		logging.LogInfoByCtxf(ctx, "Going to use existing disk image '%s' to create VM '%s'.", diskImagePath, vmName)
	} else {
		return nil, tracederrors.TracedErrorf("Disk image '%s' does not exist to create VM '%s'.", diskImagePath, vmName)
	}

	domainXml, err := kvmutilsgeneric.CreateXmlForVmAsString(createOptions)
	if err != nil {
		return nil, err
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return nil, err
	}

	// DomainCreateXML creates and starts a transient domain from the given XML,
	// which is the native equivalent of 'virsh create <xml>'.
	_, err = connection.DomainCreateXML(domainXml, 0)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to create VM '%s' from XML: %w", vmName, err)
	}

	createdVm, err = k.GetVmByName(ctx, vmName)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "VM '%s' created.", vmName)

	logging.LogInfoByCtxf(ctx, "Create KVM VM '%s' finished.", vmName)

	return createdVm, nil
}

func (k *NativeKvmHypervisor) GetVmById(vmId int) (vm kvmutilsinterfaces.VM, err error) {
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

func (k *NativeKvmHypervisor) GetVmByName(ctx context.Context, vmName string) (vm kvmutilsinterfaces.VM, err error) {
	if vmName == "" {
		return nil, tracederrors.TracedErrorEmptyString("vmName")
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

func (k *NativeKvmHypervisor) ListVms(ctx context.Context) (vms []kvmutilsinterfaces.VM, err error) {
	logging.LogInfoByCtxf(ctx, "List KVM VMs started.")

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return nil, err
	}

	// Include active and inactive domains, matching 'virsh list --all'.
	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive

	domains, _, err := connection.ConnectListAllDomains(1, flags)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to list domains: %w", err)
	}

	vms = []kvmutilsinterfaces.VM{}
	for _, domain := range domains {
		var vmToAdd kvmutilsinterfaces.VM = kvmutilsgeneric.NewKvmVm()

		err = vmToAdd.SetHypervisor(k)
		if err != nil {
			return nil, err
		}

		// An inactive domain has an ID of -1. Only set the ID for running VMs
		// to mirror the '-' shown by 'virsh list --all'.
		if domain.ID > 0 {
			vmToAdd, err = k.GetVmById(int(domain.ID))
			if err != nil {
				return nil, err
			}
		}

		err = vmToAdd.SetCachedName(domain.Name)
		if err != nil {
			return nil, err
		}

		stateValue, _, err := connection.DomainGetState(domain, 0)
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to get state of domain '%s': %w", domain.Name, err)
		}

		err = vmToAdd.SetCachedState(libvirtDomainStateToVirshString(libvirt.DomainState(stateValue)))
		if err != nil {
			return nil, err
		}

		vms = append(vms, vmToAdd)
	}

	logging.LogInfoByCtxf(ctx, "List KVM VMs finished. Collected '%d' VMs.", len(vms))

	return vms, nil
}

func (k *NativeKvmHypervisor) ListVmNames(ctx context.Context) (vmNames []string, err error) {
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

func (k *NativeKvmHypervisor) ListVmInfos(ctx context.Context) (vmInfos []kvmutilsinterfaces.VmInfo, err error) {
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

func (k *NativeKvmHypervisor) VmByNameExists(ctx context.Context, vmName string) (vmExists bool, err error) {
	if vmName == "" {
		return false, tracederrors.TracedErrorEmptyString("vmName")
	}

	logging.LogInfoByCtxf(ctx, "KVM Vm by name '%s' exists started.", vmName)

	vmNameList, err := k.ListVmNames(contextutils.WithSilent(ctx))
	if err != nil {
		return false, err
	}

	var exists bool
	for _, name := range vmNameList {
		if name == vmName {
			exists = true
			break
		}
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "KVM Vm by name '%s' exists finished. VM exists.", vmName)
	} else {
		logging.LogInfoByCtxf(ctx, "KVM Vm by name '%s' exists finished. VM does not exist.", vmName)
	}

	return exists, nil
}

func (k *NativeKvmHypervisor) DeleteVm(ctx context.Context, removeOptions *kvmutilsoptions.KvmRemoveVmOptions) (err error) {
	if removeOptions == nil {
		return tracederrors.TracedError("removeOptions is nil")
	}

	if removeOptions.VmName == "" {
		return tracederrors.TracedErrorEmptyString("removeOptions.VmName")
	}

	hostName, err := k.GetHostName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete KVM VM '%s' on host '%s' started.", removeOptions.VmName, hostName)

	vmExists, err := k.VmByNameExists(ctx, removeOptions.VmName)
	if err != nil {
		return err
	}

	if vmExists {
		connection, err := k.getLibvirtConnection(ctx)
		if err != nil {
			return err
		}

		domain, err := connection.DomainLookupByName(removeOptions.VmName)
		if err != nil {
			return tracederrors.TracedErrorf("failed to look up domain '%s': %w", removeOptions.VmName, err)
		}

		// Only destroy running domains. For inactive domains DomainDestroy
		// would fail, so we check the state first to stay idempotent.
		stateValue, _, err := connection.DomainGetState(domain, 0)
		if err != nil {
			return tracederrors.TracedErrorf("failed to get state of domain '%s': %w", removeOptions.VmName, err)
		}

		if libvirt.DomainState(stateValue) != libvirt.DomainShutoff {
			err = connection.DomainDestroy(domain)
			if err != nil {
				return tracederrors.TracedErrorf("failed to destroy domain '%s': %w", removeOptions.VmName, err)
			}
		}

		// Undefine only if the domain is persistent. Transient domains
		// disappear on destroy and would make DomainUndefine fail.
		isPersistent, err := connection.DomainIsPersistent(domain)
		if err != nil {
			return tracederrors.TracedErrorf("failed to check if domain '%s' is persistent: %w", removeOptions.VmName, err)
		}

		if isPersistent == 1 {
			err = connection.DomainUndefine(domain)
			if err != nil {
				return tracederrors.TracedErrorf("failed to undefine domain '%s': %w", removeOptions.VmName, err)
			}
		}

		logging.LogChangedByCtxf(ctx, "VM '%s' removed on host '%s'.", removeOptions.VmName, hostName)
	} else {
		logging.LogInfoByCtxf(ctx, "VM '%s' is already removed on host '%s'.", removeOptions.VmName, hostName)
	}

	if removeOptions.RemoveVolumes {
		for _, volumeName := range removeOptions.VolumeNamesToRemove {
			err = k.RemoveVolumeByName(ctx, volumeName)
			if err != nil {
				return err
			}
		}
	}

	logging.LogInfoByCtxf(ctx, "Delete KVM VM '%s' on host '%s' finished.", removeOptions.VmName, hostName)

	return nil
}

func (k *NativeKvmHypervisor) ResetVm(ctx context.Context, name string) (err error) {
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
		logging.LogInfoByCtxf(ctx, "VM '%s' is not running on host '%s'. Skip reset.", name, hostName)
		return nil
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return err
	}

	domain, err := connection.DomainLookupByName(name)
	if err != nil {
		return tracederrors.TracedErrorf("failed to look up domain '%s': %w", name, err)
	}

	err = connection.DomainReset(domain, 0)
	if err != nil {
		return tracederrors.TracedErrorf("failed to reset domain '%s': %w", name, err)
	}

	logging.LogChangedByCtxf(ctx, "VM '%s' reset on host '%s'.", name, hostName)

	return nil
}

func (k *NativeKvmHypervisor) GetDomainXmlAsString(ctx context.Context, vmName string) (domainXml string, err error) {
	if vmName == "" {
		return "", tracederrors.TracedErrorEmptyString("vmName")
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return "", err
	}

	domain, err := connection.DomainLookupByName(vmName)
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to look up domain '%s': %w", vmName, err)
	}

	domainXml, err = connection.DomainGetXMLDesc(domain, 0)
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to get XML of domain '%s': %w", vmName, err)
	}

	return domainXml, nil
}

// ---------------------------------------------------------------------------
// IP address discovery
// ---------------------------------------------------------------------------

func (k *NativeKvmHypervisor) GetIpAddress(ctx context.Context, vmName string) (ipAddress string, err error) {
	if vmName == "" {
		return "", tracederrors.TracedErrorEmptyString("vmName")
	}

	// Try multiple sources in order so both NAT ('default') and bridged ('br0')
	// VMs work:
	//   - agent: queries the qemu-guest-agent inside the VM.
	//   - lease: reads libvirt's dnsmasq DHCP leases (NAT 'default' network).
	//   - arp:   reads the host's ARP table (bridged setups).
	for _, source := range []uint32{
		uint32(libvirt.DomainInterfaceAddressesSrcAgent),
		uint32(libvirt.DomainInterfaceAddressesSrcLease),
		uint32(libvirt.DomainInterfaceAddressesSrcArp),
	} {
		// Use a silent context so the individual (expected to sometimes fail)
		// lookups do not spam the log.
		ipAddress, err = k.getIpAddressBySource(contextutils.WithSilent(ctx), vmName, source)
		if err == nil && ipAddress != "" {
			return ipAddress, nil
		}
	}

	return "", tracederrors.TracedErrorf("No IPv4 address found for VM '%s' (tried sources agent, lease, arp).", vmName)
}

func (k *NativeKvmHypervisor) getIpAddressBySource(ctx context.Context, vmName string, source uint32) (ipAddress string, err error) {
	if vmName == "" {
		return "", tracederrors.TracedErrorEmptyString("vmName")
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return "", err
	}

	domain, err := connection.DomainLookupByName(vmName)
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to look up domain '%s': %w", vmName, err)
	}

	interfaces, err := connection.DomainInterfaceAddresses(domain, source, 0)
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to get interface addresses of domain '%s': %w", vmName, err)
	}

	for _, iface := range interfaces {
		for _, addr := range iface.Addrs {
			// Only consider IPv4 addresses.
			if addr.Type != int32(libvirt.IPAddrTypeIpv4) {
				continue
			}

			if addr.Addr == "" {
				continue
			}

			// Validate the discovered IP before returning it (constitution).
			err = iputils.CheckValidIPv4(ctx, addr.Addr)
			if err != nil {
				return "", err
			}

			return addr.Addr, nil
		}
	}

	return "", tracederrors.TracedErrorf("No IPv4 address found for VM '%s' via the requested source.", vmName)
}

// ---------------------------------------------------------------------------
// Storage pools
// ---------------------------------------------------------------------------

func (k *NativeKvmHypervisor) ListStoragePools(ctx context.Context) (storagePools []*kvmutilsgeneric.KvmStoragePool, err error) {
	logging.LogInfoByCtxf(ctx, "List storage pools on KVM hypervisor started.")

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return nil, err
	}

	flags := libvirt.ConnectListStoragePoolsActive | libvirt.ConnectListStoragePoolsInactive

	pools, _, err := connection.ConnectListAllStoragePools(1, flags)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to list storage pools: %w", err)
	}

	storagePools = []*kvmutilsgeneric.KvmStoragePool{}
	for _, pool := range pools {
		poolToAdd := kvmutilsgeneric.NewKvmStoragePool()

		err = poolToAdd.SetName(pool.Name)
		if err != nil {
			return nil, err
		}

		err = poolToAdd.SetHypervisor(k)
		if err != nil {
			return nil, err
		}

		storagePools = append(storagePools, poolToAdd)
	}

	logging.LogInfoByCtxf(ctx, "List storage pools on KVM host '%s' finished. Collected '%d' storage pools.", hostname, len(storagePools))

	return storagePools, nil
}

func (k *NativeKvmHypervisor) ListStoragePoolNames(ctx context.Context) (storagePoolNames []string, err error) {
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

	return storagePoolNames, nil
}

func (k *NativeKvmHypervisor) GetStoragePoolByName(ctx context.Context, storagePoolName string) (storagePool kvmutilsinterfaces.StoragePool, err error) {
	if storagePoolName == "" {
		return nil, tracederrors.TracedErrorEmptyString("storagePoolName")
	}

	list, err := k.ListStoragePools(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range list {
		name, err := p.GetName()
		if err != nil {
			return nil, err
		}

		if name == storagePoolName {
			return p, nil
		}
	}

	return nil, tracederrors.TracedErrorf("KVM storage pool '%s' not found.", storagePoolName)
}

// ---------------------------------------------------------------------------
// Volumes
// ---------------------------------------------------------------------------

func (k *NativeKvmHypervisor) ListVolumes(ctx context.Context, storagePoolName string) (volumes []kvmutilsinterfaces.Volume, err error) {
	if storagePoolName == "" {
		return nil, tracederrors.TracedErrorEmptyString("storagePoolName")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List volumes in storage pool '%s' on KVM hypervisor '%s' started.", storagePoolName, hostname)

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return nil, err
	}

	pool, err := connection.StoragePoolLookupByName(storagePoolName)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to look up storage pool '%s': %w", storagePoolName, err)
	}

	vols, _, err := connection.StoragePoolListAllVolumes(pool, 1, 0)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to list volumes of storage pool '%s': %w", storagePoolName, err)
	}

	storagePool, err := k.GetStoragePoolByName(ctx, storagePoolName)
	if err != nil {
		return nil, err
	}

	volumes = []kvmutilsinterfaces.Volume{}
	for _, vol := range vols {
		volumeToAdd := kvmutilsgeneric.NewKvmVolume()

		err = volumeToAdd.SetName(vol.Name)
		if err != nil {
			return nil, err
		}

		err = volumeToAdd.SetStoragePool(storagePool)
		if err != nil {
			return nil, err
		}

		volumes = append(volumes, volumeToAdd)
	}

	logging.LogInfoByCtxf(ctx, "List volumes in storage pool '%s' on KVM hypervisor '%s' finished. Collected '%d' volumes.", storagePoolName, hostname, len(volumes))

	return volumes, nil
}

func (k *NativeKvmHypervisor) GetVolumes(ctx context.Context) (volumes []kvmutilsinterfaces.Volume, err error) {
	logging.LogInfoByCtxf(ctx, "Get volumes on KVM hypervisor started.")

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

	logging.LogInfoByCtxf(ctx, "Get volumes on KVM hypervisor finished. Collected '%d' volumes from '%d' pools on KVM host '%s'.", len(volumes), len(storagePools), hostname)

	return volumes, nil
}

func (k *NativeKvmHypervisor) ListVolumeNames(ctx context.Context) (volumeNames []string, err error) {
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

func (k *NativeKvmHypervisor) GetVolumeByName(ctx context.Context, volumeName string) (volume kvmutilsinterfaces.Volume, err error) {
	if volumeName == "" {
		return nil, tracederrors.TracedErrorEmptyString("volumeName")
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

func (k *NativeKvmHypervisor) VolumeByNameExists(ctx context.Context, volumeName string) (volumeExists bool, err error) {
	if volumeName == "" {
		return false, tracederrors.TracedErrorEmptyString("volumeName")
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

func (k *NativeKvmHypervisor) RemoveVolumeByName(ctx context.Context, volumeName string) (err error) {
	if volumeName == "" {
		return tracederrors.TracedErrorEmptyString("volumeName")
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

		logging.LogChangedByCtxf(ctx, "Volume '%s' on KVM hypervisor '%s' deleted.", volumeName, hostname)
	} else {
		logging.LogInfoByCtxf(ctx, "Volume '%s' on KVM hypervisor '%s' was already deleted.", volumeName, hostname)
	}

	return nil
}

func (k *NativeKvmHypervisor) DeleteVolumeByName(ctx context.Context, storagePoolName string, volumeName string) (err error) {
	if storagePoolName == "" {
		return tracederrors.TracedErrorEmptyString("storagePoolName")
	}

	if volumeName == "" {
		return tracederrors.TracedErrorEmptyString("volumeName")
	}

	hostDescription, err := k.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Remove KVM volume '%s' of storage pool '%s' on hypervisor '%s' started.", volumeName, storagePoolName, hostDescription)

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return err
	}

	pool, err := connection.StoragePoolLookupByName(storagePoolName)
	if err != nil {
		return tracederrors.TracedErrorf("failed to look up storage pool '%s': %w", storagePoolName, err)
	}

	vol, err := connection.StorageVolLookupByName(pool, volumeName)
	if err != nil {
		return tracederrors.TracedErrorf("failed to look up volume '%s' in storage pool '%s': %w", volumeName, storagePoolName, err)
	}

	err = connection.StorageVolDelete(vol, 0)
	if err != nil {
		return tracederrors.TracedErrorf("failed to delete volume '%s' in storage pool '%s': %w", volumeName, storagePoolName, err)
	}

	logging.LogChangedByCtxf(ctx, "Remove KVM volume '%s' of storage pool '%s' on hypervisor '%s' deleted.", volumeName, storagePoolName, hostDescription)

	logging.LogInfoByCtxf(ctx, "Remove KVM volume '%s' of storage pool '%s' on hypervisor '%s' finished.", volumeName, storagePoolName, hostDescription)

	return nil
}

// ---------------------------------------------------------------------------
// Networks
// ---------------------------------------------------------------------------

func (k *NativeKvmHypervisor) ListNetworks(ctx context.Context) (networks []kvmutilsinterfaces.Network, err error) {
	logging.LogInfoByCtxf(ctx, "List networks on KVM hypervisor started.")

	hostname, err := k.GetHostName()
	if err != nil {
		return nil, err
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return nil, err
	}

	// Include active and inactive networks (e.g. an inactive 'default' network),
	// matching 'virsh net-list --all'.
	flags := libvirt.ConnectListNetworksActive | libvirt.ConnectListNetworksInactive

	libvirtNetworks, _, err := connection.ConnectListAllNetworks(1, flags)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to list networks: %w", err)
	}

	networks = []kvmutilsinterfaces.Network{}
	for _, net := range libvirtNetworks {
		networkToAdd := kvmutilsgeneric.NewNetwork()

		err = networkToAdd.SetHypervisor(k)
		if err != nil {
			return nil, err
		}

		err = networkToAdd.SetName(net.Name)
		if err != nil {
			return nil, err
		}

		active, err := connection.NetworkIsActive(net)
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to get active state of network '%s': %w", net.Name, err)
		}
		err = networkToAdd.SetState(boolIntToVirshActiveString(active))
		if err != nil {
			return nil, err
		}

		autostart, err := connection.NetworkGetAutostart(net)
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to get autostart of network '%s': %w", net.Name, err)
		}
		err = networkToAdd.SetAutostart(boolIntToVirshYesNoString(autostart))
		if err != nil {
			return nil, err
		}

		persistent, err := connection.NetworkIsPersistent(net)
		if err != nil {
			return nil, tracederrors.TracedErrorf("failed to get persistent state of network '%s': %w", net.Name, err)
		}
		err = networkToAdd.SetPersistent(boolIntToVirshYesNoString(persistent))
		if err != nil {
			return nil, err
		}

		networks = append(networks, networkToAdd)
	}

	logging.LogInfoByCtxf(ctx, "List networks on KVM host '%s' finished. Collected '%d' networks.", hostname, len(networks))

	return networks, nil
}

func (k *NativeKvmHypervisor) ListNetworkNames(ctx context.Context) (networkNames []string, err error) {
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

func (k *NativeKvmHypervisor) NetworkByNameExists(ctx context.Context, networkName string) (networkExists bool, err error) {
	if networkName == "" {
		return false, tracederrors.TracedErrorEmptyString("networkName")
	}

	networkNames, err := k.ListNetworkNames(contextutils.WithSilent(ctx))
	if err != nil {
		return false, err
	}

	for _, name := range networkNames {
		if name == networkName {
			return true, nil
		}
	}

	return false, nil
}

func (k *NativeKvmHypervisor) GetNetworkByName(networkName string) (network kvmutilsinterfaces.Network, err error) {
	if networkName == "" {
		return nil, tracederrors.TracedErrorEmptyString("networkName")
	}

	networkToReturn := kvmutilsgeneric.NewNetwork()

	err = networkToReturn.SetName(networkName)
	if err != nil {
		return nil, err
	}

	err = networkToReturn.SetHypervisor(k)
	if err != nil {
		return nil, err
	}

	return networkToReturn, nil
}

func (k *NativeKvmHypervisor) StartNetworkByName(ctx context.Context, networkName string) (err error) {
	if networkName == "" {
		return tracederrors.TracedErrorEmptyString("networkName")
	}

	hostname, err := k.GetHostName()
	if err != nil {
		return err
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return err
	}

	net, err := connection.NetworkLookupByName(networkName)
	if err != nil {
		return tracederrors.TracedErrorf("No network named '%s' found on KVM host '%s': %w", networkName, hostname, err)
	}

	active, err := connection.NetworkIsActive(net)
	if err != nil {
		return tracederrors.TracedErrorf("failed to get active state of network '%s': %w", networkName, err)
	}

	if active == 1 {
		logging.LogInfoByCtxf(ctx, "Network '%s' is already active on KVM host '%s'.", networkName, hostname)
		return nil
	}

	err = connection.NetworkCreate(net)
	if err != nil {
		return tracederrors.TracedErrorf("failed to start network '%s': %w", networkName, err)
	}

	logging.LogChangedByCtxf(ctx, "Network '%s' started on KVM host '%s'.", networkName, hostname)

	return nil
}

func (k *NativeKvmHypervisor) GetParsedNetworkXml(ctx context.Context, networkName string) (parsedXml kvmutilsinterfaces.KvmNetworkXml, err error) {
	if networkName == "" {
		return nil, tracederrors.TracedErrorEmptyString("networkName")
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return nil, err
	}

	net, err := connection.NetworkLookupByName(networkName)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to look up network '%s': %w", networkName, err)
	}

	stdout, err := connection.NetworkGetXMLDesc(net, 0)
	if err != nil {
		return nil, tracederrors.TracedErrorf("failed to get XML of network '%s': %w", networkName, err)
	}

	parsed := &kvmutilsgeneric.KvmNetworkXml{}
	err = xml.Unmarshal([]byte(stdout), parsed)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse network XML for network '%s': %w", networkName, err)
	}

	return parsed, nil
}

func (k *NativeKvmHypervisor) AddIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string) (err error) {
	return k.updateIpDhcpRangeInNetwork(ctx, networkName, startIp, endIp, uint32(libvirt.NetworkUpdateCommandAddLast))
}

func (k *NativeKvmHypervisor) DeleteIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string) (err error) {
	return k.updateIpDhcpRangeInNetwork(ctx, networkName, startIp, endIp, uint32(libvirt.NetworkUpdateCommandDelete))
}

func (k *NativeKvmHypervisor) updateIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string, command uint32) (err error) {
	if networkName == "" {
		return tracederrors.TracedErrorEmptyString("networkName")
	}

	// Validate the IP boundaries via iputils (constitution).
	err = iputils.CheckValidIP(ctx, startIp)
	if err != nil {
		return err
	}

	err = iputils.CheckValidIP(ctx, endIp)
	if err != nil {
		return err
	}

	connection, err := k.getLibvirtConnection(ctx)
	if err != nil {
		return err
	}

	net, err := connection.NetworkLookupByName(networkName)
	if err != nil {
		return tracederrors.TracedErrorf("failed to look up network '%s': %w", networkName, err)
	}

	rangeXml := fmt.Sprintf("<range start='%s' end='%s'/>", startIp, endIp)

	flags := libvirt.NetworkUpdateAffectLive | libvirt.NetworkUpdateAffectConfig

	err = connection.NetworkUpdate(
		net,
		command,
		uint32(libvirt.NetworkSectionIPDhcpRange),
		-1, // parent index: append/search across all <ip> sections.
		rangeXml,
		flags,
	)
	if err != nil {
		return tracederrors.TracedErrorf("failed to update DHCP range on network '%s': %w", networkName, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// libvirtDomainStateToVirshString maps a libvirt domain state to the human
// readable state string as printed by 'virsh list', so cached states stay
// compatible with the command-executor based implementation.
func libvirtDomainStateToVirshString(state libvirt.DomainState) string {
	switch state {
	case libvirt.DomainNostate:
		return "no state"
	case libvirt.DomainRunning:
		return "running"
	case libvirt.DomainBlocked:
		return "idle"
	case libvirt.DomainPaused:
		return "paused"
	case libvirt.DomainShutdown:
		return "in shutdown"
	case libvirt.DomainShutoff:
		return "shut off"
	case libvirt.DomainCrashed:
		return "crashed"
	case libvirt.DomainPmsuspended:
		return "pmsuspended"
	default:
		return "no state"
	}
}

func boolIntToVirshActiveString(value int32) string {
	if value == 1 {
		return "active"
	}

	return "inactive"
}

func boolIntToVirshYesNoString(value int32) string {
	if value == 1 {
		return "yes"
	}

	return "no"
}

var _ = strings.TrimSpace // keep strings imported if trimming helpers are added later
