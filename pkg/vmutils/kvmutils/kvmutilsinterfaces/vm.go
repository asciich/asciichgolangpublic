package kvmutilsinterfaces

import "context"

type VM interface {
	Delete(ctx context.Context) error

	GetCachedName() (string, error)
	GetCachedState() (string, error)

	GetIpAddress(ctx context.Context) (string, error)
	GetMacAddress(ctx context.Context) (string, error)
	GetNetworkName(ctx context.Context) (string, error)

	GetDomainXmlAsString(ctx context.Context) (string, error)
	GetInfo(ctx context.Context) (VmInfo, error)

	GetVncPort(ctx context.Context) (int, error)

	IsRunning(ctx context.Context) (bool, error)

	SetHypervisor(Hypervisor) error
	SetId(id int) error

	SetCachedName(string) error
	SetCachedState(string) error
}
