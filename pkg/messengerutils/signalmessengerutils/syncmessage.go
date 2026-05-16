package signalmessengerutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type SentMessage struct {
	Destination        string `json:"destination"`
	DestinationNumber  string `json:"destinationNumber"`
	DestinationUuid    string `json:"destinationUuid"`
	Timestamp          int64  `json:"timestamp"`
	Message            string `json:"message"`
	ExpiresInSeconds   int    `json:"expiresInSeconds"`
	IsExpirationUpdate bool   `json:"isExpirationUpdate"`
	ViewOnce           bool   `json:"viewOnce"`
}

func (s *SentMessage) GetContentAsString() (string, error) {
	if s.Message == "" {
		return "", tracederrors.TracedError("message not set in signals Sent message")
	}

	return s.Message, nil
}

type SyncMessage struct {
	SentMessage *SentMessage `json:"sentMessage,omitempty"`
}

func (s *SyncMessage) GetSentMessage() (*SentMessage, error) {
	if s.SentMessage == nil {
		return nil, tracederrors.TracedError("SentMessage is nil ")
	}

	return s.SentMessage, nil
}

func (s *SyncMessage) GetContentAsString() (string, error) {
	sentMessage, err := s.GetSentMessage()
	if err != nil {
		return "", err
	}

	return sentMessage.GetContentAsString()
}
