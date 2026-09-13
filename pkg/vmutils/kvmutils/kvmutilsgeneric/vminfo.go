package kvmutilsgeneric

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type VmInfo struct {
	Name       string `json:"name"`
	MacAddress string `json:"mac_address"`
}

func NewKvmVmInfo() (k *VmInfo) {
	return new(VmInfo)
}

func (k *VmInfo) GetMacAddress() (macAddress string, err error) {
	if k.MacAddress == "" {
		return "", tracederrors.TracedErrorf("MacAddress not set")
	}

	return k.MacAddress, nil
}

func (k *VmInfo) GetName() (name string, err error) {
	if k.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return k.Name, nil
}

func (k *VmInfo) GetNameAndMacAddress() (name string, macAddress string, err error) {
	name, err = k.GetName()
	if err != nil {
		return "", "", err
	}

	macAddress, err = k.GetMacAddress()
	if err != nil {
		return "", "", err
	}

	return name, macAddress, nil
}

func (k *VmInfo) SetMacAddress(macAddress string) (err error) {
	if macAddress == "" {
		return tracederrors.TracedErrorf("macAddress is empty string")
	}

	k.MacAddress = macAddress

	return nil
}

func (k *VmInfo) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	k.Name = name

	return nil
}
