package asciichgolangpublic


type KvmVmInfo struct {
	Name       string `json:"name"`
	MacAddress string `json:"mac_address"`
}

func NewKvmVmInfo() (k *KvmVmInfo) {
	return new(KvmVmInfo)
}

func (k *KvmVmInfo) GetMacAddress() (macAddress string, err error) {
	if k.MacAddress == "" {
		return "", TracedErrorf("MacAddress not set")
	}

	return k.MacAddress, nil
}

func (k *KvmVmInfo) GetName() (name string, err error) {
	if k.Name == "" {
		return "", TracedErrorf("Name not set")
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
		LogGoErrorFatal(err)
	}

	return macAddress
}

func (k *KvmVmInfo) MustGetName() (name string) {
	name, err := k.GetName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name
}

func (k *KvmVmInfo) MustGetNameAndMacAddress() (name string, macAddress string) {
	name, macAddress, err := k.GetNameAndMacAddress()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return name, macAddress
}

func (k *KvmVmInfo) MustSetMacAddress(macAddress string) {
	err := k.SetMacAddress(macAddress)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (k *KvmVmInfo) MustSetName(name string) {
	err := k.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (k *KvmVmInfo) SetMacAddress(macAddress string) (err error) {
	if macAddress == "" {
		return TracedErrorf("macAddress is empty string")
	}

	k.MacAddress = macAddress

	return nil
}

func (k *KvmVmInfo) SetName(name string) (err error) {
	if name == "" {
		return TracedErrorf("name is empty string")
	}

	k.Name = name

	return nil
}
