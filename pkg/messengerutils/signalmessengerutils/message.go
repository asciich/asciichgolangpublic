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

func (m *Message) GetSenderAccountAsString() (string, error) {
	envelope, err := m.GetEnvelope()
	if err != nil {
		return "", err
	}

	return envelope.GetSenderAccount()
}

func (m *Message) GetTimestampMilliseconds() (int64, error) {
	envelope, err := m.GetEnvelope()
	if err != nil {
		return 0, err
	}

	return envelope.GetTimestampMilliseconds()
}

func (m *Message) GetRecipientsAsStringSlice() ([]string, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (m *Message) IsSenderAccount(accountToMatch string) (bool, error) {
	if accountToMatch == "" {
		return false, tracederrors.TracedErrorEmptyString("accountToMatch")
	}

	sender, err := m.GetSenderAccountAsString()
	if err != nil {
		return false, err
	}

	return sender == accountToMatch, nil
}

func (m *Message) IsDataMessage() (bool, error) {
	envelope, err := m.GetEnvelope()
	if err != nil {
		return false, err
	}

	return envelope.DataMessage != nil, nil
}
