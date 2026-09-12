package kvmutilsinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsgeneric"
)

type Hypervisor interface {
	GetHostName() (string, error)

	GetParsedNetworkXml(ctx context.Context, networkName string) (*kvmutilsgeneric.KvmNetworkXml, error)

	AddIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string) error
	DeleteIpDhcpRangeInNetwork(ctx context.Context, networkName string, startIp string, endIp string) error
}
