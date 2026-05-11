package signalclirestapiutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/messengerutils/signalmessengerutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ReceiveCacheServerOptions struct {
	SignalResetClientApiUrl string
	Interval                string
	CacheSize               int
	AccountNumber           string
}

func (r *ReceiveCacheServerOptions) GetAccountNumber() (string, error) {
	if r.AccountNumber == "" {
		return "", tracederrors.TracedError("AccountNumber not set")
	}

	err := signalmessengerutils.CheckAccountNumber(r.AccountNumber)
	if err != nil {
		return "", err
	}

	return r.AccountNumber, nil
}
