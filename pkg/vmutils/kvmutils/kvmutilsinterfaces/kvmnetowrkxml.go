package kvmutilsinterfaces

type KvmNetworkXml interface {
	GetForwardMode() (string, error)
	GetIpAddress() (string, error)
	GetIpNetmask() (string, error)

	GetIpDhcpRangeStart() (string, error)
	GetIpDhcpRangeEnd() (string, error)
}
