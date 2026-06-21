package sshutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type SSHKeyPair struct {
	PublicKey  *SSHPublicKey
	PrivateKey *SSHPrivateKey
}

func (s *SSHKeyPair) Validate(ctx context.Context) error {
	_, err := s.GetPublicKey()
	if err != nil {
		return err
	}

	_, err = s.GetPrivateKey()
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHKeyPair) GetPublicKey() (*SSHPublicKey, error) {
	if s.PublicKey == nil {
		return nil, tracederrors.TracedError("PublicKey not set")
	}

	return s.PublicKey, nil
}

func (s *SSHKeyPair) GetPrivateKey() (*SSHPrivateKey, error) {
	if s.PrivateKey == nil {
		return nil, tracederrors.TracedError("SSHPrivateKey not set")
	}

	return s.PrivateKey, nil
}
