package signalmessengerutils

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

// Message represents a signal message.
type Message struct {
	Envelope *Envelope `json:"envelope"`
	Account  string    `json:"account"`
}

func (m *Message) GetEnvelope() (*Envelope, error) {
	if m.Envelope == nil {
		return nil, tracederrors.TracedError("envelope not set")
	}

	return m.Envelope, nil
}

func (m *Message) GetContentAsString() (string, error) {
	envelope, err := m.GetEnvelope()
	if err != nil {
		return "", err
	}

	return envelope.GetContentAsString()
}
