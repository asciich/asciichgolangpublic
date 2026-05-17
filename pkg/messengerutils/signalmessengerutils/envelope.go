package signalmessengerutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Envelope struct {
	Source                   string          `json:"source"`
	SourceNumber             string          `json:"sourceNumber"`
	SourceUuid               string          `json:"sourceUuid"`
	SourceName               string          `json:"sourceName"`
	SourceDevice             int             `json:"sourceDevice"`
	Timestamp                int64           `json:"timestamp"`
	ServerReceivedTimestamp  int64           `json:"serverReceivedTimestamp"`
	ServerDeliveredTimestamp int64           `json:"serverDeliveredTimestamp"`
	DataMessage              *DataMessage    `json:"dataMessage,omitempty"`
	SyncMessage              *SyncMessage    `json:"syncMessage,omitempty"`
	ReceiptMessage           *ReceiptMessage `json:"receiptMessage,omitempty"`
}

func (e *Envelope) GetContentAsString() (string, error) {
	if e.DataMessage != nil {
		content, err := e.DataMessage.GetContentAsString()
		if err == nil {
			return content, nil
		}
	}

	if e.SyncMessage != nil {
		content, err := e.SyncMessage.GetContentAsString()
		if err == nil {
			return content, nil
		}
	}

	return "", tracederrors.TracedError("Unable to get signal message content as string")
}
