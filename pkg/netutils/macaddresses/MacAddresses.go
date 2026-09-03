package macaddresses

import (
	"regexp"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CheckStringIsAMacAddress(input string) (err error) {
	isMacAddress := IsStringAMacAddress(input)
	if !isMacAddress {
		return tracederrors.TracedErrorf("'%s' is not a valid mac address", input)
	}

	return nil
}

func IsStringAMacAddress(input string) (isMacAddress bool) {
	r := regexp.MustCompile("^[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}$")
	return r.MatchString(input)
}
