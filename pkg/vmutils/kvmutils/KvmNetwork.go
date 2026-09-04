package kvmutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmNetwork struct {
	hypervisor *KVMHypervisor

	name       string
	state      string
	autostart  string
	persistent string
}

func NewKvmNetwork() (ret *KvmNetwork) {
	return new(KvmNetwork)
}

func (n *KvmNetwork) SetHypervisor(hypervisor *KVMHypervisor) (err error) {
	if hypervisor == nil {
		return tracederrors.TracedError("hypervisor is nil")
	}

	n.hypervisor = hypervisor

	return nil
}

func (n *KvmNetwork) GetHypervisor() (hypervisor *KVMHypervisor, err error) {
	if n.hypervisor == nil {
		return nil, tracederrors.TracedError("hypervisor not set")
	}

	return n.hypervisor, nil
}

func (n *KvmNetwork) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	n.name = name

	return nil
}

func (n *KvmNetwork) GetName() (name string, err error) {
	if n.name == "" {
		return "", tracederrors.TracedError("name not set")
	}

	return n.name, nil
}

func (n *KvmNetwork) SetState(state string) (err error) {
	if state == "" {
		return tracederrors.TracedErrorEmptyString("state")
	}

	n.state = state

	return nil
}

func (n *KvmNetwork) GetState() (state string, err error) {
	if n.state == "" {
		return "", tracederrors.TracedError("state not set")
	}

	return n.state, nil
}

func (n *KvmNetwork) IsActive() (isActive bool, err error) {
	state, err := n.GetState()
	if err != nil {
		return false, err
	}

	return state == "active", nil
}

func (n *KvmNetwork) SetAutostart(autostart string) (err error) {
	if autostart == "" {
		return tracederrors.TracedErrorEmptyString("autostart")
	}

	n.autostart = autostart

	return nil
}

func (n *KvmNetwork) GetAutostart() (autostart string, err error) {
	if n.autostart == "" {
		return "", tracederrors.TracedError("autostart not set")
	}

	return n.autostart, nil
}

func (n *KvmNetwork) SetPersistent(persistent string) (err error) {
	if persistent == "" {
		return tracederrors.TracedErrorEmptyString("persistent")
	}

	n.persistent = persistent

	return nil
}

func (n *KvmNetwork) GetPersistent() (persistent string, err error) {
	if n.persistent == "" {
		return "", tracederrors.TracedError("persistent not set")
	}

	return n.persistent, nil
}
