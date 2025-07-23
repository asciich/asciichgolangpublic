package asciichgolangpublic

import (
	"regexp"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func CheckStringIsAMacAddress(input string) (isMacAddress bool, err error) {
	isMacAddress = IsStringAMacAddress(input)
	if !isMacAddress {
		return false, tracederrors.TracedErrorf("'%s' is not a valid mac address", input)
	}

	return true, nil
}

func IsStringAMacAddress(input string) (isMacAddress bool) {
	r := regexp.MustCompile("^[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}$")
	return r.MatchString(input)
}

func MustCheckStringIsAMacAddress(input string) (isMacAddress bool) {
	isMacAddress, err := CheckStringIsAMacAddress(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isMacAddress
}
