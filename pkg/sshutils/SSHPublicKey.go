package sshutils

import (
	"context"
	"path/filepath"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
)

type SSHPublicKey struct {
	// Type. E.g. "ssh-ras" or "ssh-ed25519"
	KeyType string

	// The effective key material
	KeyMaterial string

	// Name and host usually added in "user@hoast" form at the end of a line in the *.pup key file.
	KeyUserName string
	KeyUserHost string
}

func NewSSHPublicKey() (sshPublicKey *SSHPublicKey) {
	return new(SSHPublicKey)
}

func (k *SSHPublicKey) Equals(other *SSHPublicKey) (isEqual bool) {
	if other == nil {
		return false
	}

	if k.KeyType != other.KeyType {
		return false
	}

	if k.KeyMaterial != other.KeyMaterial {
		return false
	}

	if k.KeyUserName != other.KeyUserName {
		return false
	}

	if k.KeyUserHost != other.KeyUserHost {
		return false
	}

	return true
}

func (k *SSHPublicKey) GetAsPublicKeyLine() (publicKeyLine string, err error) {
	keyType, err := k.GetKeyType()
	if err != nil {
		return "", err
	}

	publicKeyLine += keyType

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
	if len(k.KeyUserHost) <= 0 {
		return "", err
	}

	return k.KeyUserHost, nil
}

func (k *SSHPublicKey) GetKeyMaterialAsString() (keyMaterial string, err error) {
	if len(k.KeyMaterial) <= 0 {
		return "", tracederrors.TracedError("key material not set")
	}

	return k.KeyMaterial, nil
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
	if len(k.KeyUserName) <= 0 {
		return "", tracederrors.TracedErrorf("keyUserName is empty string. Available data: '%v'", *k)
	}

	return k.KeyUserName, nil
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
		k.KeyMaterial = keyMaterial
	} else if numberOfSpacesInKeyMaterial > 0 && numberOfSpacesInKeyMaterial <= 3 {
		splittedAllElements := strings.Split(keyMaterial, " ")
		splitted := slicesutils.TrimSpace(splittedAllElements)

		for _, possibleKeyType := range []string{"ssh-rsa", "ssh-ed25519", "ecdsa-sha2-nistp256"} {
			if slices.Contains(splitted, possibleKeyType) {
				k.KeyType = possibleKeyType
				splitted = slicesutils.RemoveMatchingStrings(splitted, possibleKeyType)
				break
			}
		}

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
					"unable to extract key material. len(splitted) = '%v' != 1 as expected. key material is '%s' and splitted is '%v'",
					len(splitted),
					keyMaterial,
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

		k.KeyMaterial = keyMaterialToAdd

		for _, part := range splittedAllElements {
			if strings.Contains(part, "@") {
				splitted := strings.Split(part, "@")

				if len(splitted) >= 1 {
					k.KeyUserName = splitted[0]
				}

				if len(splitted) >= 2 {
					k.KeyUserHost = splitted[1]
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

func (k *SSHPublicKey) WriteToFile(ctx context.Context, outputFile files.File) (err error) {
	if outputFile == nil {
		return tracederrors.TracedError("outputFile is nil")
	}

	sshKeyLine, err := k.GetAsPublicKeyLine()
	if err != nil {
		return err
	}

	sshKeyLine = stringsutils.EnsureEndsWithExactlyOneLineBreak(sshKeyLine)
	err = outputFile.WriteString(sshKeyLine, contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHPublicKey) GetKeyMaterial() (keyMaterial string, err error) {
	if s.KeyMaterial == "" {
		return "", tracederrors.TracedErrorf("keyMaterial not set")
	}

	return s.KeyMaterial, nil
}

// Key type like "ssh-rsa" or "ssh-ed25519"
func (s *SSHPublicKey) GetKeyType() (keyType string, err error) {
	if s.KeyType == "" {
		return "", tracederrors.TracedError("keyType not set")
	}

	return s.KeyType, nil
}

func (s *SSHPublicKey) GetKeyUserHost() (keyUserHost string, err error) {
	if s.KeyUserHost == "" {
		return "", tracederrors.TracedErrorf("keyUserHost not set")
	}

	return s.KeyUserHost, nil
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

func (s *SSHPublicKey) MustWriteToFile(ctx context.Context, outputFile files.File) {
	err := s.WriteToFile(ctx, outputFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) SetKeyMaterial(keyMaterial string) (err error) {
	if keyMaterial == "" {
		return tracederrors.TracedErrorf("keyMaterial is empty string")
	}

	s.KeyMaterial = keyMaterial

	return nil
}

func (s *SSHPublicKey) SetKeyUserHost(keyUserHost string) (err error) {
	if keyUserHost == "" {
		return tracederrors.TracedErrorf("keyUserHost is empty string")
	}

	s.KeyUserHost = keyUserHost

	return nil
}

func (s *SSHPublicKey) SetKeyUserName(keyUserName string) (err error) {
	if keyUserName == "" {
		return tracederrors.TracedErrorf("keyUserName is empty string")
	}

	s.KeyUserName = keyUserName

	return nil
}

func LoadPublicKeysFromFile(ctx context.Context, sshKeysFile files.File) (sshKeys []*SSHPublicKey, err error) {
	if sshKeysFile == nil {
		return nil, tracederrors.TracedError("sshKeysFile is nil")
	}

	logging.LogInfoByCtxf(ctx, "Load SSH public keys from file '%s' started.", sshKeysFile)

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

	logging.LogInfoByCtxf(ctx, "Load SSH public keys from file '%s' finished.", sshKeysFile)

	return sshKeys, nil
}

func MustLoadPublicKeysFromFile(ctx context.Context, sshKeysFile files.File) (sshKeys []*SSHPublicKey) {
	sshKeys, err := LoadPublicKeysFromFile(ctx, sshKeysFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return sshKeys
}

func MustLoadPublicKeyFromString(keyMaterial string) (key *SSHPublicKey) {
	key, err := LoadPublicKeyFromString(keyMaterial)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return key
}

func LoadPublicKeyFromString(keyMaterial string) (key *SSHPublicKey, err error) {
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

func GetCurrentUsersSshDirectory() (sshDir files.Directory, err error) {
	homeDir, err := userutils.GetHomeDirectory()
	if err != nil {
		return nil, err
	}

	sshDir, err = homeDir.GetSubDirectory(".ssh")
	if err != nil {
		return nil, err
	}

	return sshDir, nil
}

func GetSshPublicKey(verbose bool) (sshPublicKey *SSHPublicKey, err error) {
	sshDirectory, err := GetCurrentUsersSshDirectory()
	if err != nil {
		return nil, err
	}

	sshPublicKey = new(SSHPublicKey)
	err = sshPublicKey.LoadFromSshDir(sshDirectory, verbose)
	if err != nil {
		return nil, err
	}

	return sshPublicKey, nil
}
