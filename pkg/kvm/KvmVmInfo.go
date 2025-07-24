package kvm

import (
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type KvmVmInfo struct {
	Name       string `json:"name"`
	MacAddress string `json:"mac_address"`
}

func NewKvmVmInfo() (k *KvmVmInfo) {
	return new(KvmVmInfo)
}

func (k *KvmVmInfo) GetMacAddress() (macAddress string, err error) {
	if k.MacAddress == "" {
		return "", tracederrors.TracedErrorf("MacAddress not set")
	}

	return k.MacAddress, nil
}

func (k *KvmVmInfo) GetName() (name string, err error) {
	if k.Name == "" {
		return "", tracederrors.TracedErrorf("Name not set")
	}

	return k.Name, nil
}

func (k *KvmVmInfo) GetNameAndMacAddress() (name string, macAddress string, err error) {
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

func (k *KvmVmInfo) MustGetMacAddress() (macAddress string) {
	macAddress, err := k.GetMacAddress()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return macAddress
}

func (k *KvmVmInfo) MustGetName() (name string) {
	name, err := k.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (k *KvmVmInfo) MustGetNameAndMacAddress() (name string, macAddress string) {
	name, macAddress, err := k.GetNameAndMacAddress()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name, macAddress
}

func (k *KvmVmInfo) MustSetMacAddress(macAddress string) {
	err := k.SetMacAddress(macAddress)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmVmInfo) MustSetName(name string) {
	err := k.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (k *KvmVmInfo) SetMacAddress(macAddress string) (err error) {
	if macAddress == "" {
		return tracederrors.TracedErrorf("macAddress is empty string")
	}

	k.MacAddress = macAddress

	return nil
}

func (k *KvmVmInfo) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	k.Name = name

	return nil
}
