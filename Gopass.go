package asciichgolangpublic

import (
	"fmt"
	"path/filepath"
	"strings"

)

type GopassService struct {}

func Gopass() (gopass *GopassService) {
	return new(GopassService)
}

func NewGopassService() (g *GopassService) {
	return new(GopassService)
}

func (g *GopassService) CredentialExists(fullCredentialPath string) (credentialExists bool, err error) {
	fullCredentialPath = strings.TrimSpace(fullCredentialPath)

	if len(fullCredentialPath) <= 0 {
		return false, TracedError("fullCredentailPath is empty string")
	}

	credentailList, err := g.GetCredentialNameList()
	if err != nil {
		return false, err
	}

	return Slices().ContainsString(credentailList, fullCredentialPath), nil
}

func (g *GopassService) Generate(credentialName string, verbose bool) (generatedCredential *GopassCredential, err error) {
	if credentialName == "" {
		return nil, TracedError("credentailName is empty string")
	}

	newPassword, err := RandomGenerator().GetRandomString(16)
	if err != nil {
		return nil, err
	}

	credential, err := g.GetGopassCredentialByName(credentialName)
	if err != nil {
		return nil, err
	}

	err = credential.SetByString(newPassword)
	if err != nil {
		return nil, err
	}

	if verbose {
		LogInfof("Gopass credentail '%s' generated.", credentialName)
	}

	return credential, nil
}

func (g *GopassService) GetCredential(getOptions *GopassSecretOptions) (credential *GopassCredential, err error) {
	if getOptions == nil {
		return nil, TracedError("getOptions is nil")
	}

	name, err := getOptions.GetGopassPath()
	if err != nil {
		return nil, err
	}

	credential = NewGopassCredential()
	err = credential.SetName(name)
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (g *GopassService) GetCredentialList() (credentials []*GopassCredential, err error) {
	outLines, err := Bash().RunCommandAndGetStdoutAsLines(
		&RunCommandOptions{
			Command: []string{"gopass", "list", "-f"},
		},
	)
	if err != nil {
		return nil, err
	}

	credentials = []*GopassCredential{}
	for _, line := range outLines {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		credentialToAdd := NewGopassCredential()
		err = credentialToAdd.SetName(line)
		if err != nil {
			return nil, err
		}

		credentials = append(credentials, credentialToAdd)
	}

	return credentials, nil
}

func (g *GopassService) GetCredentialNameList() (credentialNames []string, err error) {
	credentials, err := g.GetCredentialList()
	if err != nil {
		return nil, err
	}

	credentialNames = []string{}
	for _, c := range credentials {
		nameToAdd, err := c.GetName()
		if err != nil {
			return nil, err
		}

		credentialNames = append(credentialNames, nameToAdd)
	}

	return credentialNames, nil
}

func (g *GopassService) GetCredentialValueAsString(getOptions *GopassSecretOptions) (credentialValue string, err error) {
	if getOptions == nil {
		return "", TracedError("getOptions is nil")
	}

	credential, err := g.GetCredential(getOptions)
	if err != nil {
		return
	}

	credentialValue, err = credential.GetAsString()
	if err != nil {
		return
	}

	return credentialValue, nil
}

func (g *GopassService) GetCredentialValueAsStringByPath(secretPath string) (secretValue string, err error) {
	if secretPath == "" {
		return "", TracedError("secretPath is empty string")
	}

	secretValue, err = g.GetCredentialValueAsString(&GopassSecretOptions{
		SecretRootDirectoryPath: filepath.Dir(secretPath),
		SecretBasename:          filepath.Base(secretPath),
	})
	if err != nil {
		return "", err
	}

	return secretValue, nil
}

func (g *GopassService) GetCredentialValueOrEmptyIfUnsetAsStringByPath(secretPath string) (credentialValue string, err error) {
	if secretPath == "" {
		return "", TracedErrorEmptyString(secretPath)
	}

	credential, err := g.GetCredential(&GopassSecretOptions{
		SecretRootDirectoryPath: filepath.Dir(secretPath),
		SecretBasename:          filepath.Base(secretPath),
	})
	if err != nil {
		return "", err
	}

	credentialExists, err := credential.Exists()
	if err != nil {
		return "", err
	}

	if !credentialExists {
		return "", nil
	}

	credentialValue, err = credential.GetAsString()
	if err != nil {
		return "", err
	}

	return credentialValue, nil
}

func (g *GopassService) GetGopassCredentialByName(name string) (credential *GopassCredential, err error) {
	if name == "" {
		return nil, TracedError("name is empty string")
	}

	credential, err = GetGopassCredentialByName(name)
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (g *GopassService) GetSslCertificate(getOptions *GopassSecretOptions) (cert *X509Certificate, err error) {
	if getOptions == nil {
		return nil, TracedError("getOptions is nil")
	}

	credential, err := g.GetCredential(getOptions)
	if err != nil {
		return nil, err
	}

	cert, err = credential.GetSslCertificate()
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (g *GopassService) InsertFile(fileToInsert File, gopassOptions *GopassSecretOptions) (err error) {
	if fileToInsert == nil {
		return TracedError("fileToInsert is nil")
	}

	if gopassOptions == nil {
		return TracedError("gopassOptions is nil")
	}

	fileToInsertPath, err := fileToInsert.GetLocalPath()
	if err != nil {
		return err
	}

	fileExists, err := fileToInsert.Exists(gopassOptions.Verbose)
	if err != nil {
		return err
	}

	if !fileExists {
		return TracedError("fileToInsert does not exist in file system.")
	}

	gopassPath, err := gopassOptions.GetGopassPath()
	if err != nil {
		return err
	}

	if !gopassOptions.Overwrite {
		secretExists, err := g.SecretNameExist(gopassPath)
		if err != nil {
			return err
		}

		if secretExists {
			return TracedErrorf("Secret '%v' already exists in gopass.", gopassPath)
		}
	}

	insertCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("cat '%s' | gopass cat '%s'", fileToInsertPath, gopassPath),
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command: insertCommand,
		},
	)
	if err != nil {
		return err
	}

	err = g.WriteInfoToGopass(gopassPath)
	if err != nil {
		return err
	}

	if gopassOptions.Verbose {
		LogInfof("Added file '%v' to gopass as '%v'", fileToInsertPath, gopassPath)
	}

	return nil
}

func (g *GopassService) InsertSecret(secretToInsert string, gopassOptions *GopassSecretOptions) (err error) {
	if len(secretToInsert) <= 0 {
		return TracedError("secretToInsert is empty string")
	}

	if gopassOptions == nil {
		return TracedError("gopassOptions is nil")
	}

	gopassPath, err := gopassOptions.GetGopassPath()
	if err != nil {
		return err
	}

	if !gopassOptions.Overwrite {
		secretExists, err := g.SecretNameExist(gopassPath)
		if err != nil {
			return err
		}

		if secretExists {
			return TracedErrorf("Secret '%v' already exists in gopass.", gopassPath)
		}
	}

	insertCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("echo '%s' | gopass insert -f '%s'", secretToInsert, gopassPath),
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command: insertCommand,
		},
	)
	if err != nil {
		return err
	}

	err = g.WriteInfoToGopass(gopassPath)
	if err != nil {
		return err
	}

	if gopassOptions.Verbose {
		LogInfof("Added credentail '%v' to gopass.", gopassPath)
	}

	return nil
}

func (g *GopassService) MustCredentialExists(fullCredentialPath string) (credentialExists bool) {
	credentialExists, err := g.CredentialExists(fullCredentialPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credentialExists
}

func (g *GopassService) MustGenerate(credentialName string, verbose bool) (generatedCredential *GopassCredential) {
	generatedCredential, err := g.Generate(credentialName, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return generatedCredential
}

func (g *GopassService) MustGetCredential(getOptions *GopassSecretOptions) (credential *GopassCredential) {
	credential, err := g.GetCredential(getOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credential
}

func (g *GopassService) MustGetCredentialList() (credentials []*GopassCredential) {
	credentials, err := g.GetCredentialList()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credentials
}

func (g *GopassService) MustGetCredentialNameList() (credentialNames []string) {
	credentialNames, err := g.GetCredentialNameList()
	if err != nil {
		LogFatalf("gopass.GetCredentialNameList failed: '%v'", err)
	}

	return credentialNames
}

func (g *GopassService) MustGetCredentialValue(getOptions *GopassSecretOptions) (credentialValue string) {
	credentialValue, err := g.GetCredentialValueAsString(getOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credentialValue
}

func (g *GopassService) MustGetCredentialValueAsString(getOptions *GopassSecretOptions) (credentialValue string) {
	credentialValue, err := g.GetCredentialValueAsString(getOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credentialValue
}

func (g *GopassService) MustGetCredentialValueAsStringByPath(secretPath string) (secretValue string) {
	secretValue, err := g.GetCredentialValueAsStringByPath(secretPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return secretValue
}

func (g *GopassService) MustGetCredentialValueOrEmptyIfUnsetAsStringByPath(secretPath string) (credentialValue string) {
	credentialValue, err := g.GetCredentialValueOrEmptyIfUnsetAsStringByPath(secretPath)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credentialValue
}

func (g *GopassService) MustGetGopassCredentialByName(name string) (credential *GopassCredential) {
	credential, err := g.GetGopassCredentialByName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credential
}

func (g *GopassService) MustGetSslCertificate(getOptions *GopassSecretOptions) (cert *X509Certificate) {
	cert, err := g.GetSslCertificate(getOptions)
	if err != nil {
		LogFatalf("Gopass.GetSslCertificate: '%v'", err)
	}

	return cert
}

func (g *GopassService) MustInsertFile(fileToInsert File, gopassOptions *GopassSecretOptions) {
	err := g.InsertFile(fileToInsert, gopassOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassService) MustInsertSecret(secretToInsert string, gopassOptions *GopassSecretOptions) {
	err := g.InsertSecret(secretToInsert, gopassOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassService) MustSecretNameExist(secretName string) (secretExists bool) {
	secretExists, err := g.SecretNameExist(secretName)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return secretExists
}

func (g *GopassService) MustSync(verbose bool) {
	err := g.Sync(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassService) MustWriteInfoToGopass(gopassPath string) {
	err := g.WriteInfoToGopass(gopassPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassService) MustWriteSecretIntoTemporaryFile(getOptions *GopassSecretOptions) (temporaryFile File) {
	temporaryFile, err := g.WriteSecretIntoTemporaryFile(getOptions)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return temporaryFile
}

func (g *GopassService) SecretNameExist(secretName string) (secretExists bool, err error) {
	secretName = strings.TrimSpace(secretName)
	if len(secretName) <= 0 {
		return false, TracedError("secretName is empty string")
	}

	secretNames, err := g.GetCredentialNameList()
	if err != nil {
		return false, err
	}

	secretExists = Slices().ContainsString(secretNames, secretName)
	return secretExists, nil
}

func (g *GopassService) Sync(verbose bool) (err error) {
	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command:            []string{"gopass", "sync"},
			Verbose:            verbose,
			LiveOutputOnStdout: verbose,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (g *GopassService) WriteInfoToGopass(gopassPath string) (err error) {
	gopassPath = strings.TrimSpace(gopassPath)
	if len(gopassPath) <= 0 {
		return TracedError("gopassPath is empty string")
	}

	gopassPath += "_info"

	infoString := fmt.Sprintf("This secret was added by '%v'", GetBinaryInfo().GetInfoString())
	infoString = strings.ReplaceAll(infoString, "'", "\"")

	insertCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("echo '%v' | gopass insert -f '%v'", infoString, gopassPath),
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command: insertCommand,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (g *GopassService) WriteSecretIntoTemporaryFile(getOptions *GopassSecretOptions) (temporaryFile File, err error) {
	if getOptions == nil {
		return nil, TracedError("getOptions is nil")
	}

	credential, err := g.GetCredential(getOptions)
	if err != nil {
		return nil, err
	}

	temporaryFile, err = credential.WriteIntoTemporaryFile(getOptions.Verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}
