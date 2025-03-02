package sshutils

import (
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type SSHPublicKey struct {
	keyMaterial string
	keyUserName string
	keyUserHost string
}

func NewSSHPublicKey() (sshPublicKey *SSHPublicKey) {
	return new(SSHPublicKey)
}

func (k *SSHPublicKey) Equals(other *SSHPublicKey) (isEqual bool) {
	logging.LogFatalf("SSHPublicKey.Equals NOT implemented")
	return false
}

func (k *SSHPublicKey) GetAsPublicKeyLine() (publicKeyLine string, err error) {
	publicKeyLine = "ssh-rsa"

	keyMaterial, err := k.GetKeyMaterialAsString()
	if err != nil {
		return "", err
	}
	publicKeyLine += " " + keyMaterial

	userAtHost, err := k.GetKeyUserAtHost()
	if err != nil {
		return "", err
	}
	publicKeyLine += " " + userAtHost

	return publicKeyLine, nil
}

func (k *SSHPublicKey) GetKeyHostName() (hostName string, err error) {
	if len(k.keyUserHost) <= 0 {
		return "", err
	}

	return k.keyUserHost, nil
}

func (k *SSHPublicKey) GetKeyMaterialAsString() (keyMaterial string, err error) {
	if len(k.keyMaterial) <= 0 {
		return "", tracederrors.TracedError("key material not set")
	}

	return k.keyMaterial, nil
}

func (k *SSHPublicKey) GetKeyUserAtHost() (userAtHost string, err error) {
	username, err := k.GetKeyUserName()
	if err != nil {
		return "", err
	}

	hostname, err := k.GetKeyHostName()
	if err != nil {
		return "", err
	}

	userAtHost = username + "@" + hostname

	return userAtHost, nil
}

func (k *SSHPublicKey) GetKeyUserName() (keyUserName string, err error) {
	if len(k.keyUserName) <= 0 {
		return "", tracederrors.TracedErrorf("keyUserName is empty string. Available data: '%v'", *k)
	}

	return k.keyUserName, nil
}

func (k *SSHPublicKey) LoadFromSshDir(sshDirectory files.Directory, verbose bool) (err error) {
	if sshDirectory == nil {
		return tracederrors.TracedError("sshDirectory is nil")
	}

	sshDirPath, err := sshDirectory.GetLocalPath()
	if err != nil {
		return err
	}

	exists, err := sshDirectory.Exists(verbose)
	if err != nil {
		return err
	}

	if !exists {
		return tracederrors.TracedErrorf("ssh key directory '%v' does not exist", sshDirPath)
	}

	keyFilePath := filepath.Join(sshDirPath, "id_rsa.pub")
	keyFile, err := files.GetLocalFileByPath(keyFilePath)
	if err != nil {
		return err
	}

	keyMaterial, err := keyFile.ReadAsString()
	if err != nil {
		return err
	}

	err = k.SetFromString(keyMaterial)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogInfof("Loaded ssh public key from '%v'", keyFilePath)
	}

	return nil
}

func (k *SSHPublicKey) MustGetKeyHostName() (hostName string) {
	hostName, err := k.GetKeyHostName()
	if err != nil {
		logging.LogFatalf("sshPublicKey.GetKeyHostName failed: '%v'", err)
	}

	return hostName
}

func (k *SSHPublicKey) MustGetKeyMaterialAsString() (keyMaterial string) {
	keyMaterial, err := k.GetKeyMaterialAsString()
	if err != nil {
		logging.LogFatalf("sshPublicKey.GetKeyMaterialAsString failed: '%v'", err)
	}

	return keyMaterial
}

func (k *SSHPublicKey) MustGetKeyUserName() (keyUserName string) {
	keyUserName, err := k.GetKeyUserName()
	if err != nil {
		logging.LogFatalf("sshPublicKey.GetKeyUserName failed: '%v'", err)
	}

	return keyUserName
}

func (k *SSHPublicKey) MustSetFromString(keyMaterial string) {
	err := k.SetFromString(keyMaterial)
	if err != nil {
		logging.LogFatalf("sshPublicKey.SetFromString failed: '%v'", err)
	}
}

func (k *SSHPublicKey) SetFromString(keyMaterial string) (err error) {
	keyMaterial = strings.TrimSpace(keyMaterial)
	if len(keyMaterial) <= 0 {
		return tracederrors.TracedError("keyMaterial is empty string")
	}

	numberOfSpacesInKeyMaterial := strings.Count(keyMaterial, " ")
	if numberOfSpacesInKeyMaterial == 0 {
		k.keyMaterial = keyMaterial
	} else if slicesutils.ContainsInt([]int{1, 2, 3}, numberOfSpacesInKeyMaterial) {
		splittedAllElements := strings.Split(keyMaterial, " ")
		splitted := slicesutils.TrimSpace(splittedAllElements)
		splitted = slicesutils.RemoveMatchingStrings(splitted, "ssh-rsa")
		splitted, err = slicesutils.RemoveStringsWhichContains(splitted, "@")
		if err != nil {
			return err
		}

		var keyMaterialToAdd string = ""
		if len(splitted) == 1 {
			keyMaterialToAdd = splitted[0]
		} else {
			firstElement := splitted[0]
			if strings.HasPrefix(firstElement, "AAA") {
				keyMaterialToAdd = splitted[0]
			} else {
				return tracederrors.TracedErrorf(
					"unable to extract key material. len(splitted) = '%v' != 1 as expected. splitted is '%v'",
					len(splitted),
					splitted,
				)
			}
		}

		keyMaterialToAdd = strings.TrimSpace(keyMaterialToAdd)
		if len(keyMaterialToAdd) <= 0 {
			return tracederrors.TracedErrorf(
				"unable to extract key material. keyMaterialToAdd is empty string calculated from '%v'",
				keyMaterial,
			)
		}

		k.keyMaterial = keyMaterialToAdd

		for _, part := range splittedAllElements {
			if strings.Contains(part, "@") {
				splitted := strings.Split(part, "@")

				if len(splitted) >= 1 {
					k.keyUserName = splitted[0]
				}

				if len(splitted) >= 2 {
					k.keyUserHost = splitted[1]
				}

				break
			}
		}
	} else {
		return tracederrors.TracedErrorf(
			"unable to extract key material. numberOfSpacesInKeyMaterial is '%v' from '%s'",
			numberOfSpacesInKeyMaterial,
			keyMaterial,
		)
	}

	return nil
}

func (k *SSHPublicKey) WriteToFile(outputFile files.File, verbose bool) (err error) {
	if outputFile == nil {
		return tracederrors.TracedError("outputFile is nil")
	}

	sshKeyLine, err := k.GetAsPublicKeyLine()
	if err != nil {
		return err
	}

	sshKeyLine = stringsutils.EnsureEndsWithExactlyOneLineBreak(sshKeyLine)
	err = outputFile.WriteString(sshKeyLine, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHPublicKey) GetKeyMaterial() (keyMaterial string, err error) {
	if s.keyMaterial == "" {
		return "", tracederrors.TracedErrorf("keyMaterial not set")
	}

	return s.keyMaterial, nil
}

func (s *SSHPublicKey) GetKeyUserHost() (keyUserHost string, err error) {
	if s.keyUserHost == "" {
		return "", tracederrors.TracedErrorf("keyUserHost not set")
	}

	return s.keyUserHost, nil
}

func (s *SSHPublicKey) MustGetAsPublicKeyLine() (publicKeyLine string) {
	publicKeyLine, err := s.GetAsPublicKeyLine()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return publicKeyLine
}

func (s *SSHPublicKey) MustGetKeyMaterial() (keyMaterial string) {
	keyMaterial, err := s.GetKeyMaterial()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyMaterial
}

func (s *SSHPublicKey) MustGetKeyUserAtHost() (userAtHost string) {
	userAtHost, err := s.GetKeyUserAtHost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return userAtHost
}

func (s *SSHPublicKey) MustGetKeyUserHost() (keyUserHost string) {
	keyUserHost, err := s.GetKeyUserHost()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyUserHost
}

func (s *SSHPublicKey) MustLoadFromSshDir(sshDirectory files.Directory, verbose bool) {
	err := s.LoadFromSshDir(sshDirectory, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) MustSetKeyMaterial(keyMaterial string) {
	err := s.SetKeyMaterial(keyMaterial)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) MustSetKeyUserHost(keyUserHost string) {
	err := s.SetKeyUserHost(keyUserHost)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) MustSetKeyUserName(keyUserName string) {
	err := s.SetKeyUserName(keyUserName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) MustWriteToFile(outputFile files.File, verbose bool) {
	err := s.WriteToFile(outputFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) SetKeyMaterial(keyMaterial string) (err error) {
	if keyMaterial == "" {
		return tracederrors.TracedErrorf("keyMaterial is empty string")
	}

	s.keyMaterial = keyMaterial

	return nil
}

func (s *SSHPublicKey) SetKeyUserHost(keyUserHost string) (err error) {
	if keyUserHost == "" {
		return tracederrors.TracedErrorf("keyUserHost is empty string")
	}

	s.keyUserHost = keyUserHost

	return nil
}

func (s *SSHPublicKey) SetKeyUserName(keyUserName string) (err error) {
	if keyUserName == "" {
		return tracederrors.TracedErrorf("keyUserName is empty string")
	}

	s.keyUserName = keyUserName

	return nil
}

func LoadPublicKeysFromFile(sshKeysFile files.File, verbose bool) (sshKeys []*SSHPublicKey, err error) {
	if sshKeysFile == nil {
		return nil, tracederrors.TracedError("sshKeysFile is nil")
	}

	filePath, err := sshKeysFile.GetLocalPath()
	if err != nil {
		return nil, err
	}

	if verbose {
		logging.LogInfof("Load SSH public keys from file '%s' started.", filePath)
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
		logging.LogInfof("Load SSH public keys from file '%s' finished.", filePath)
	}

	return sshKeys, nil
}

func MustLoadPublicKeysFromFile(sshKeysFile files.File, verbose bool) (sshKeys []*SSHPublicKey) {
	sshKeys, err := LoadPublicKeysFromFile(sshKeysFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sshKeys
}

func MustLoadPublicKeyFromString(keyMaterial string) (key *SSHPublicKey) {
	key, err := LoadPublicKeyFromFile(keyMaterial)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return key
}

func LoadPublicKeyFromFile(keyMaterial string) (key *SSHPublicKey, err error) {
	if keyMaterial == "" {
		return nil, tracederrors.TracedErrorEmptyString("keyMaterial")
	}

	key = NewSSHPublicKey()
	err = key.SetFromString(keyMaterial)
	if err != nil {
		return nil, err
	}

	return key, nil
}
