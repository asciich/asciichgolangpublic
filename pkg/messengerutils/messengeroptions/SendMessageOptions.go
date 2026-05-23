package messengeroptions

import "github.com/asciich/asciichgolangpublic/pkg/tracederrors"

type SendMessageOptions struct {
	Message       string
	SenderAccount string
	Recipients    []string
}

func (s *SendMessageOptions) GetMessage() (string, error) {
	if s.Message == "" {
		return "", tracederrors.TracedError("Message not set")
	}

	return s.Message, nil
}

func (s *SendMessageOptions) GetSenderAccount() (string, error) {
	if s.SenderAccount == "" {
		return "", tracederrors.TracedError("SenderAccount not set")
	}

	return s.SenderAccount, nil
}

func (s *SendMessageOptions) GetRecipients() ([]string, error) {
	if s.Recipients == nil {
		return nil, tracederrors.TracedError("Recipients not set")
	}

	return s.Recipients, nil
}
