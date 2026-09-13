package kvmutilsinterfaces

type StoragePool interface {
	GetName() (string, error)

	GetHostName() (string, error)

	GetHypervisor() (Hypervisor, error)
}
