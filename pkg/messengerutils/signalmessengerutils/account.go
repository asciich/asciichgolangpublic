package signalmessengerutils

import (
	"regexp"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var regexpAccountNumber = regexp.MustCompile(`^\+[1-9]\d{6,14}$`)

func IsAccountNumber(accountNumber string) bool {
	return regexpAccountNumber.Match([]byte(accountNumber))
}

func CheckAccountNumber(accountNumber string) error {
	if !IsAccountNumber(accountNumber) {
		return tracederrors.TracedErrorf("'%s' is not a valid account number for signal.", accountNumber)
	}

	return nil
}
