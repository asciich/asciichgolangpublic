package vmoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type CreateVmOptions struct {
	VmName          string
	DiskImageName   string
	MacAddress      string
	BridgeInterface string
}

func NewCreateVmOptions() (c *CreateVmOptions) {
	return new(CreateVmOptions)
}

func (c *CreateVmOptions) GetVmName() (vmName string, err error) {
	if c.VmName == "" {
		return "", tracederrors.TracedErrorf("VmName not set")
	}

	return c.VmName, nil
}

func (c *CreateVmOptions) GetDiskImageName() (diskImageName string, err error) {
	if c.DiskImageName == "" {
		return "", tracederrors.TracedErrorf("DiskImageName not set")
	}

	return c.DiskImageName, nil
}

func (c *CreateVmOptions) GetMacAddress() (macAddress string, err error) {
	if c.MacAddress == "" {
		return "", tracederrors.TracedErrorf("MacAddress not set")
	}

	return c.MacAddress, nil
}

func (c *CreateVmOptions) GetBridgeInterface() (bridgeInterface string, err error) {
	if c.BridgeInterface == "" {
		return "", tracederrors.TracedErrorf("BridgeInterface not set")
	}

	return c.BridgeInterface, nil
}

func (c *CreateVmOptions) SetVmName(vmName string) (err error) {
	if vmName == "" {
		return tracederrors.TracedErrorf("vmName is empty string")
	}

	c.VmName = vmName

	return nil
}

func (c *CreateVmOptions) SetDiskImageName(diskImageName string) (err error) {
	if diskImageName == "" {
		return tracederrors.TracedErrorf("diskImageName is empty string")
	}

	c.DiskImageName = diskImageName

	return nil
}

func (c *CreateVmOptions) SetMacAddress(macAddress string) (err error) {
	if macAddress == "" {
		return tracederrors.TracedErrorf("macAddress is empty string")
	}

	c.MacAddress = macAddress

	return nil
}

func (c *CreateVmOptions) SetBridgeInterface(bridgeInterface string) (err error) {
	if bridgeInterface == "" {
		return tracederrors.TracedErrorf("bridgeInterface is empty string")
	}

	c.BridgeInterface = bridgeInterface

	return nil
}

func (c *CreateVmOptions) GetBridgeInterfaceOrEmptyStringIfUnset() (string, error) {
	return c.BridgeInterface, nil
}
