package kvmutilsinterfaces

import "context"

type Network interface {
	GetAutostart() (string, error)

	GetName() (string, error)

	GetDhcpRange(ctx context.Context) (string, string, error)
	GetForwardMode(ctx context.Context) (string, error)
	GetGatewayIp(ctx context.Context) (string, error)
	GetNetmask(ctx context.Context) (string, error)
	GetPersistent() (string, error)
	GetState() (string, error)

	IsActive() (bool, error)
	IsNatted(ctx context.Context) (bool, error)

	SetDhcpStartIp(ctx context.Context, startIp string) error
}
