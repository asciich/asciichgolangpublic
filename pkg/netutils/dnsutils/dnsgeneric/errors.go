package dnsgeneric

import "errors"

var ErrDnsDomainNotFound = errors.New("dns domain not found")
var ErrDnsDomainRecordNotFound = errors.New("dns domain record not found")
var ErrDnsDomainRecordAlreadyExists = errors.New("dns domain record already exists")

func IsErrDnsDomainRecordNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrDnsDomainRecordNotFound)
}

func IsErrDnsDomainNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrDnsDomainNotFound)
}

// Returns true if the error is not found in both cases for domain and domain record.
func IsErrNotFound(err error) bool {
	if err == nil {
		return false
	}

	if IsErrDnsDomainNotFound(err) {
		return true
	}

	return IsErrDnsDomainRecordNotFound(err)
}

func IsErrDnsDomainRecordAlreadyExists(err error) bool {
	if err != nil {
		return false
	}

	return errors.Is(err, ErrDnsDomainRecordAlreadyExists)
}
