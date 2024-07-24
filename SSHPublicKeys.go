package asciichgolangpublic

import (
	"strings"

)

type SSHPublicKeysService struct{}

func NewSSHPublicKeysService() (sshPublicKeys *SSHPublicKeysService) {
	return new(SSHPublicKeysService)
}

func SSHPublicKeys() (sshPublicKeys *SSHPublicKeysService) {
	return NewSSHPublicKeysService()
}

func (s *SSHPublicKeysService) LoadKeysFromFile(sshKeysFile File, verbose bool) (sshKeys []*SSHPublicKey, err error) {
	if sshKeysFile == nil {
		return nil, TracedError("sshKeysFile is nil")
	}

	filePath, err := sshKeysFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	if verbose {
		LogInfof("Load SSH public keys from file '%s' started.", filePath)
	}

	lines, err := sshKeysFile.ReadAsLinesWithoutComments()
	if err != nil {
		return nil, err
	}

	sshKeys = []*SSHPublicKey{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		keyToAdd := NewSSHPublicKey()
		err = keyToAdd.SetFromString(line)
		if err != nil {
			return nil, err
		}

		sshKeys = append(sshKeys, keyToAdd)
	}

	if verbose {
		LogInfof("Load SSH public keys from file '%s' finished.", filePath)
	}

	return sshKeys, nil
}

func (s *SSHPublicKeysService) MustLoadKeysFromFile(sshKeysFile File, verbose bool) (sshKeys []*SSHPublicKey) {
	sshKeys, err := s.LoadKeysFromFile(sshKeysFile, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sshKeys
}
