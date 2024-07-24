package asciichgolangpublic

import (
	"regexp"
)

type MacAddressesService struct{}

func MacAddresses() (m *MacAddressesService) {
	return NewMacAddressesService()
}

func NewMacAddressesService() (m *MacAddressesService) {
	return new(MacAddressesService)
}

func (m *MacAddressesService) CheckStringIsAMacAddress(input string) (isMacAddress bool, err error) {
	isMacAddress = m.IsStringAMacAddress(input)
	if !isMacAddress {
		return false, TracedErrorf("'%s' is not a valid mac address", input)
	}

	return true, nil
}

func (m *MacAddressesService) IsStringAMacAddress(input string) (isMacAddress bool) {
	r := regexp.MustCompile("^[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}$")
	return r.MatchString(input)
}

func (m *MacAddressesService) MustCheckStringIsAMacAddress(input string) (isMacAddress bool) {
	isMacAddress, err := m.CheckStringIsAMacAddress(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return isMacAddress
}
