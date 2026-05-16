package signalmessengerutils

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type DataMessage struct {
	Timestamp          int64  `json:"timestamp"`
	Message            string `json:"message"`
	ExpiresInSeconds   int    `json:"expiresInSeconds"`
	IsExpirationUpdate bool   `json:"isExpirationUpdate"`
	ViewOnce           bool   `json:"viewOnce"`
}

func (d *DataMessage) GetContentAsString() (string, error) {
	if d.Message == "" {
		return "", tracederrors.TracedErrorEmptyString("Message not set")
	}

	return d.Message, nil
}
