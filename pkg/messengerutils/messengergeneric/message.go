package messengergeneric

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type Message struct {
	Message       string
	SenderAccount string
	Recipients    []string
}

func (s *Message) GetContentAsString() (string, error) {
	if s.Message == "" {
		return "", tracederrors.TracedError("Message not set")
	}

	return s.Message, nil
}

func (s *Message) GetSenderAccountAsString() (string, error) {
	if s.SenderAccount == "" {
		return "", tracederrors.TracedError("SenderAccount not set")
	}

	return s.SenderAccount, nil
}

func (s *Message) GetRecipientsAsStringSlice() ([]string, error) {
	if s.Recipients == nil {
		return nil, tracederrors.TracedError("Recipients not set")
	}

	if len(s.Recipients) <= 0 {
		return nil, tracederrors.TracedError("Recipients not set: slice is empty")
	}

	return s.Recipients, nil
}

func (s *Message) GetTimestampMilliseconds() (int64, error) {
	return 0, tracederrors.TracedErrorNotImplemented()
}

func (s *Message) IsSenderAccount(accountToMatch string) (bool, error) {
	if accountToMatch == "" {
		return false, tracederrors.TracedErrorEmptyString("accountToMatch")
	}

	sender, err := s.GetSenderAccountAsString()
	if err != nil {
		return false, nil
	}

	return sender == accountToMatch, nil
}

func (s *Message) IsDataMessage() (bool, error) {
	return true, nil
}
