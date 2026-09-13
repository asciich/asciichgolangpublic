package kvmutilsinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
)

type Hypervisor interface {
	AddIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string) error

	DeleteIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string) error
	DeleteVolumeByName(ctx context.Context, storagePoolName string, volumeName string) error
	DeleteVm(ctx context.Context, options *kvmutilsoptions.KvmRemoveVmOptions) error

	GetDomainXmlAsString(ctx context.Context, vmName string) (domainXml string, err error)

	GetHostName() (string, error)
	GetHostDescription() (string, error)

	GetIpAddress(ctx context.Context, vmName string) (string, error)
	GetNetworkByName(networkName string) (Network, error)
	GetParsedNetworkXml(ctx context.Context, networkName string) (KvmNetworkXml, error)
	GetVmByName(ctx context.Context, name string) (VM, error)

	ListNetworks(ctx context.Context) ([]Network, error)
	ListNetworkNames(ctx context.Context) ([]string, error)
	ListStoragePoolNames(ctx context.Context) ([]string, error)
	ListVolumes(ctx context.Context, storagePoolName string) ([]Volume, error)
	ListVolumeNames(ctx context.Context) ([]string, error)
	ListVms(ctx context.Context) ([]VM, error)

	ResetVm(ctx context.Context, vmName string) error

	StartNetworkByName(ctx context.Context, networkName string) error
}
