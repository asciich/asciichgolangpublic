package asciichgolangpublic

import (
	"path/filepath"
	"strings"

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
	LogFatalf("SSHPublicKey.Equals NOT implemented")
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
		return "", TracedError("key material not set")
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
		return "", TracedErrorf("keyUserName is empty string. Available data: '%v'", *k)
	}

	return k.keyUserName, nil
}

func (k *SSHPublicKey) LoadFromSshDir(sshDirectory Directory, verbose bool) (err error) {
	if sshDirectory == nil {
		return TracedError("sshDirectory is nil")
	}

	sshDirPath, err := sshDirectory.GetLocalPath()
	if err != nil {
		return err
	}

	exists, err := sshDirectory.Exists()
	if err != nil {
		return err
	}

	if !exists {
		return TracedErrorf("ssh key directory '%v' does not exist", sshDirPath)
	}

	keyFilePath := filepath.Join(sshDirPath, "id_rsa.pub")
	keyFile, err := GetLocalFileByPath(keyFilePath)
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
		LogInfof("Loaded ssh public key from '%v'", keyFilePath)
	}

	return nil
}

func (k *SSHPublicKey) MustGetKeyHostName() (hostName string) {
	hostName, err := k.GetKeyHostName()
	if err != nil {
		LogFatalf("sshPublicKey.GetKeyHostName failed: '%v'", err)
	}

	return hostName
}

func (k *SSHPublicKey) MustGetKeyMaterialAsString() (keyMaterial string) {
	keyMaterial, err := k.GetKeyMaterialAsString()
	if err != nil {
		LogFatalf("sshPublicKey.GetKeyMaterialAsString failed: '%v'", err)
	}

	return keyMaterial
}

func (k *SSHPublicKey) MustGetKeyUserName() (keyUserName string) {
	keyUserName, err := k.GetKeyUserName()
	if err != nil {
		LogFatalf("sshPublicKey.GetKeyUserName failed: '%v'", err)
	}

	return keyUserName
}

func (k *SSHPublicKey) MustSetFromString(keyMaterial string) {
	err := k.SetFromString(keyMaterial)
	if err != nil {
		LogFatalf("sshPublicKey.SetFromString failed: '%v'", err)
	}
}

func (k *SSHPublicKey) SetFromString(keyMaterial string) (err error) {
	keyMaterial = strings.TrimSpace(keyMaterial)
	if len(keyMaterial) <= 0 {
		return TracedError("keyMaterial is empty string")
	}

	numberOfSpacesInKeyMaterial := strings.Count(keyMaterial, " ")
	if numberOfSpacesInKeyMaterial == 0 {
		k.keyMaterial = keyMaterial
	} else if Slices().ContainsInt([]int{1, 2, 3}, numberOfSpacesInKeyMaterial) {
		splittedAllElements := strings.Split(keyMaterial, " ")
		splitted := Slices().TrimSpace(splittedAllElements)
		splitted = Slices().RemoveMatchingStrings(splitted, "ssh-rsa")
		splitted, err = Slices().RemoveStringsWhichContains(splitted, "@")
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
				return TracedErrorf(
					"unable to extract key material. len(splitted) = '%v' != 1 as expected. splitted is '%v'",
					len(splitted),
					splitted,
				)
			}
		}

		keyMaterialToAdd = strings.TrimSpace(keyMaterialToAdd)
		if len(keyMaterialToAdd) <= 0 {
			return TracedErrorf(
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
		return TracedErrorf(
			"unable to extract key material. numberOfSpacesInKeyMaterial is '%v' from '%s'",
			numberOfSpacesInKeyMaterial,
			keyMaterial,
		)
	}

	return nil
}

func (k *SSHPublicKey) WriteToFile(outputFile File, verbose bool) (err error) {
	if outputFile == nil {
		return TracedError("outputFile is nil")
	}

	sshKeyLine, err := k.GetAsPublicKeyLine()
	if err != nil {
		return err
	}

	sshKeyLine = Strings().EnsureEndsWithExactlyOneLineBreak(sshKeyLine)
	err = outputFile.WriteString(sshKeyLine, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHPublicKey) GetKeyMaterial() (keyMaterial string, err error) {
	if s.keyMaterial == "" {
		return "", TracedErrorf("keyMaterial not set")
	}

	return s.keyMaterial, nil
}

func (s *SSHPublicKey) GetKeyUserHost() (keyUserHost string, err error) {
	if s.keyUserHost == "" {
		return "", TracedErrorf("keyUserHost not set")
	}

	return s.keyUserHost, nil
}

func (s *SSHPublicKey) MustGetAsPublicKeyLine() (publicKeyLine string) {
	publicKeyLine, err := s.GetAsPublicKeyLine()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return publicKeyLine
}

func (s *SSHPublicKey) MustGetKeyMaterial() (keyMaterial string) {
	keyMaterial, err := s.GetKeyMaterial()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keyMaterial
}

func (s *SSHPublicKey) MustGetKeyUserAtHost() (userAtHost string) {
	userAtHost, err := s.GetKeyUserAtHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return userAtHost
}

func (s *SSHPublicKey) MustGetKeyUserHost() (keyUserHost string) {
	keyUserHost, err := s.GetKeyUserHost()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keyUserHost
}

func (s *SSHPublicKey) MustLoadFromSshDir(sshDirectory Directory, verbose bool) {
	err := s.LoadFromSshDir(sshDirectory, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) MustSetKeyMaterial(keyMaterial string) {
	err := s.SetKeyMaterial(keyMaterial)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) MustSetKeyUserHost(keyUserHost string) {
	err := s.SetKeyUserHost(keyUserHost)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) MustSetKeyUserName(keyUserName string) {
	err := s.SetKeyUserName(keyUserName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) MustWriteToFile(outputFile File, verbose bool) {
	err := s.WriteToFile(outputFile, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (s *SSHPublicKey) SetKeyMaterial(keyMaterial string) (err error) {
	if keyMaterial == "" {
		return TracedErrorf("keyMaterial is empty string")
	}

	s.keyMaterial = keyMaterial

	return nil
}

func (s *SSHPublicKey) SetKeyUserHost(keyUserHost string) (err error) {
	if keyUserHost == "" {
		return TracedErrorf("keyUserHost is empty string")
	}

	s.keyUserHost = keyUserHost

	return nil
}

func (s *SSHPublicKey) SetKeyUserName(keyUserName string) (err error) {
	if keyUserName == "" {
		return TracedErrorf("keyUserName is empty string")
	}

	s.keyUserName = keyUserName

	return nil
}
