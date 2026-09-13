package kvmutilsinterfaces

type VmInfo interface {
	SetName(string) error
	SetMacAddress(string) error
}
