package gopass

import (
	"fmt"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/binaryinfo"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/randomgenerator"
	x509utils "github.com/asciich/asciichgolangpublic/tlsutils/x509utils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func CredentialExists(fullCredentialPath string) (credentialExists bool, err error) {
	fullCredentialPath = strings.TrimSpace(fullCredentialPath)

	if len(fullCredentialPath) <= 0 {
		return false, tracederrors.TracedError("fullCredentailPath is empty string")
	}

	credentailList, err := GetCredentialNameList()
	if err != nil {
		return false, err
	}

	return slices.Contains(credentailList, fullCredentialPath), nil
}

func Generate(credentialName string, verbose bool) (generatedCredential *GopassCredential, err error) {
	if credentialName == "" {
		return nil, tracederrors.TracedError("credentailName is empty string")
	}

	newPassword, err := randomgenerator.GetRandomString(16)
	if err != nil {
		return nil, err
	}

	credential, err := GetGopassCredentialByName(credentialName)
	if err != nil {
		return nil, err
	}

	err = credential.SetByString(newPassword)
	if err != nil {
		return nil, err
	}

	if verbose {
		logging.LogInfof("Gopass credentail '%s' generated.", credentialName)
	}

	return credential, nil
}

func GetCredential(getOptions *parameteroptions.GopassSecretOptions) (credential *GopassCredential, err error) {
	if getOptions == nil {
		return nil, tracederrors.TracedError("getOptions is nil")
	}

	name, err := getOptions.GetSecretPath()
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

func GetCredentialList() (credentials []*GopassCredential, err error) {
	outLines, err := commandexecutor.Bash().RunCommandAndGetStdoutAsLines(
		&parameteroptions.RunCommandOptions{
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

func GetCredentialNameList() (credentialNames []string, err error) {
	credentials, err := GetCredentialList()
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

func GetCredentialValueAsString(getOptions *parameteroptions.GopassSecretOptions) (credentialValue string, err error) {
	if getOptions == nil {
		return "", tracederrors.TracedError("getOptions is nil")
	}

	credential, err := GetCredential(getOptions)
	if err != nil {
		return
	}

	credentialValue, err = credential.GetAsString()
	if err != nil {
		return
	}

	return credentialValue, nil
}

func GetCredentialValueAsStringByPath(secretPath string) (secretValue string, err error) {
	if secretPath == "" {
		return "", tracederrors.TracedError("secretPath is empty string")
	}

	secretValue, err = GetCredentialValueAsString(
		&parameteroptions.GopassSecretOptions{
			SecretPath: secretPath,
		},
	)
	if err != nil {
		return "", err
	}

	return secretValue, nil
}

func GetCredentialValueOrEmptyIfUnsetAsStringByPath(secretPath string) (credentialValue string, err error) {
	if secretPath == "" {
		return "", tracederrors.TracedErrorEmptyString(secretPath)
	}

	credential, err := GetCredential(
		&parameteroptions.GopassSecretOptions{
			SecretPath: secretPath,
		},
	)
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

func GetSslCertificate(getOptions *parameteroptions.GopassSecretOptions) (cert *x509utils.X509Certificate, err error) {
	if getOptions == nil {
		return nil, tracederrors.TracedError("getOptions is nil")
	}

	credential, err := GetCredential(getOptions)
	if err != nil {
		return nil, err
	}

	cert, err = credential.GetSslCertificate()
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func InsertFileByString(fileContent string, gopassOptions *parameteroptions.GopassSecretOptions) (err error) {
	if gopassOptions == nil {
		return tracederrors.TracedErrorNil("gopassOptions")
	}

	gopassPath, err := gopassOptions.GetSecretPath()
	if err != nil {
		return err
	}

	if !gopassOptions.Overwrite {
		secretExists, err := SecretNameExist(gopassPath)
		if err != nil {
			return err
		}

		if secretExists {
			return tracederrors.TracedErrorf("Secret '%v' already exists in gopass.", gopassPath)
		}
	}

	insertCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("gpass cat '%s'", gopassPath),
	}

	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command:     insertCommand,
			StdinString: fileContent,
		},
	)
	if err != nil {
		return err
	}

	err = WriteInfoToGopass(gopassPath)
	if err != nil {
		return err
	}

	if gopassOptions.Verbose {
		logging.LogChangedf("Added file content to gopass as '%s'", gopassPath)
	}

	return nil
}

func InsertFile(fileToInsert files.File, gopassOptions *parameteroptions.GopassSecretOptions) (err error) {
	if fileToInsert == nil {
		return tracederrors.TracedError("fileToInsert is nil")
	}

	if gopassOptions == nil {
		return tracederrors.TracedError("gopassOptions is nil")
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
		return tracederrors.TracedError("fileToInsert does not exist in file system.")
	}

	gopassPath, err := gopassOptions.GetSecretPath()
	if err != nil {
		return err
	}

	if !gopassOptions.Overwrite {
		secretExists, err := SecretNameExist(gopassPath)
		if err != nil {
			return err
		}

		if secretExists {
			return tracederrors.TracedErrorf("Secret '%v' already exists in gopass.", gopassPath)
		}
	}

	insertCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("cat '%s' | gopass cat '%s'", fileToInsertPath, gopassPath),
	}

	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: insertCommand,
		},
	)
	if err != nil {
		return err
	}

	err = WriteInfoToGopass(gopassPath)
	if err != nil {
		return err
	}

	if gopassOptions.Verbose {
		logging.LogInfof("Added file '%v' to gopass as '%v'", fileToInsertPath, gopassPath)
	}

	return nil
}

func InsertSecret(secretToInsert string, gopassOptions *parameteroptions.GopassSecretOptions) (err error) {
	if len(secretToInsert) <= 0 {
		return tracederrors.TracedError("secretToInsert is empty string")
	}

	if gopassOptions == nil {
		return tracederrors.TracedError("gopassOptions is nil")
	}

	gopassPath, err := gopassOptions.GetSecretPath()
	if err != nil {
		return err
	}

	if !gopassOptions.Overwrite {
		secretExists, err := SecretNameExist(gopassPath)
		if err != nil {
			return err
		}

		if secretExists {
			return tracederrors.TracedErrorf("Secret '%v' already exists in gopass.", gopassPath)
		}
	}

	insertCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("echo '%s' | gopass insert -f '%s'", secretToInsert, gopassPath),
	}

	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: insertCommand,
		},
	)
	if err != nil {
		return err
	}

	err = WriteInfoToGopass(gopassPath)
	if err != nil {
		return err
	}

	if gopassOptions.Verbose {
		logging.LogInfof("Added credentail '%v' to gopass.", gopassPath)
	}

	return nil
}

func MustCredentialExists(fullCredentialPath string) (credentialExists bool) {
	credentialExists, err := CredentialExists(fullCredentialPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return credentialExists
}

func MustGenerate(credentialName string, verbose bool) (generatedCredential *GopassCredential) {
	generatedCredential, err := Generate(credentialName, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return generatedCredential
}

func MustGetCredential(getOptions *parameteroptions.GopassSecretOptions) (credential *GopassCredential) {
	credential, err := GetCredential(getOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return credential
}

func MustGetCredentialList() (credentials []*GopassCredential) {
	credentials, err := GetCredentialList()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return credentials
}

func MustGetCredentialNameList() (credentialNames []string) {
	credentialNames, err := GetCredentialNameList()
	if err != nil {
		logging.LogFatalf("gopass.GetCredentialNameList failed: '%v'", err)
	}

	return credentialNames
}

func MustGetCredentialValue(getOptions *parameteroptions.GopassSecretOptions) (credentialValue string) {
	credentialValue, err := GetCredentialValueAsString(getOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return credentialValue
}

func MustGetCredentialValueAsString(getOptions *parameteroptions.GopassSecretOptions) (credentialValue string) {
	credentialValue, err := GetCredentialValueAsString(getOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return credentialValue
}

func MustGetCredentialValueAsStringByPath(secretPath string) (secretValue string) {
	secretValue, err := GetCredentialValueAsStringByPath(secretPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return secretValue
}

func MustGetCredentialValueOrEmptyIfUnsetAsStringByPath(secretPath string) (credentialValue string) {
	credentialValue, err := GetCredentialValueOrEmptyIfUnsetAsStringByPath(secretPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return credentialValue
}

func MustGetSslCertificate(getOptions *parameteroptions.GopassSecretOptions) (cert *x509utils.X509Certificate) {
	cert, err := GetSslCertificate(getOptions)
	if err != nil {
		logging.LogFatalf("Gopass.GetSslCertificate: '%v'", err)
	}

	return cert
}

func MustInsertFile(fileToInsert files.File, gopassOptions *parameteroptions.GopassSecretOptions) {
	err := InsertFile(fileToInsert, gopassOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func MustInsertSecret(secretToInsert string, gopassOptions *parameteroptions.GopassSecretOptions) {
	err := InsertSecret(secretToInsert, gopassOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func MustSecretNameExist(secretName string) (secretExists bool) {
	secretExists, err := SecretNameExist(secretName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return secretExists
}

func MustSync(verbose bool) {
	err := Sync(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func MustWriteInfoToGopass(gopassPath string) {
	err := WriteInfoToGopass(gopassPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func MustWriteSecretIntoTemporaryFile(getOptions *parameteroptions.GopassSecretOptions) (temporaryFile files.File) {
	temporaryFile, err := WriteSecretIntoTemporaryFile(getOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return temporaryFile
}

func SecretNameExist(secretName string) (secretExists bool, err error) {
	secretName = strings.TrimSpace(secretName)
	if len(secretName) <= 0 {
		return false, tracederrors.TracedError("secretName is empty string")
	}

	secretNames, err := GetCredentialNameList()
	if err != nil {
		return false, err
	}

	return slicesutils.ContainsString(secretNames, secretName), nil
}

func Sync(verbose bool) (err error) {
	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
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

func WriteInfoToGopass(gopassPath string) (err error) {
	gopassPath = strings.TrimSpace(gopassPath)
	if len(gopassPath) <= 0 {
		return tracederrors.TracedError("gopassPath is empty string")
	}

	gopassPath += "_info"

	infoString := fmt.Sprintf("This secret was added by '%v'", binaryinfo.GetInfoString())
	infoString = strings.ReplaceAll(infoString, "'", "\"")

	insertCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("echo '%v' | gopass insert -f '%v'", infoString, gopassPath),
	}

	_, err = commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: insertCommand,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func WriteSecretIntoTemporaryFile(getOptions *parameteroptions.GopassSecretOptions) (temporaryFile files.File, err error) {
	if getOptions == nil {
		return nil, tracederrors.TracedError("getOptions is nil")
	}

	credential, err := GetCredential(getOptions)
	if err != nil {
		return nil, err
	}

	temporaryFile, err = credential.WriteIntoTemporaryFile(getOptions.Verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func MustCreateRootCaAndAddToGopass(createOptions *x509utils.X509CreateCertificateOptions, gopassOptions *parameteroptions.GopassSecretOptions) {
	err := CreateRootCaAndAddToGopass(createOptions, gopassOptions)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func CreateRootCaAndAddToGopass(createOptions *x509utils.X509CreateCertificateOptions, gopassOptions *parameteroptions.GopassSecretOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedError("createOptions is nil")
	}

	if gopassOptions == nil {
		return tracederrors.TracedError("gopassOptions is nil")
	}

	if createOptions.Verbose {
		logging.LogInfo("Create root CA and add to gopass started.")
	}

	certHandler := x509utils.GetNativeX509CertificateHandler()

	caCert, caKey, err := certHandler.CreateRootCaCertificate(createOptions)
	if err != nil {
		return err
	}

	ceCertPem, err := x509utils.EncodeCertificateAsPEMString(caCert)
	if err != nil {
		return err
	}

	caKeyPem, err := x509utils.EncodePrivateKeyAsPEMString(caKey)
	if err != nil {
		return err
	}

	certOptions := gopassOptions.GetDeepCopy()
	err = certOptions.SetBaseName("rootCa.crt")
	if err != nil {
		return err
	}

	err = InsertFileByString(ceCertPem, certOptions)
	if err != nil {
		return err
	}

	keyOptions := gopassOptions.GetDeepCopy()
	err = keyOptions.SetBaseName("rootCa.key")
	if err != nil {
		return err
	}

	err = InsertFileByString(caKeyPem, keyOptions)
	if err != nil {
		return err
	}

	if createOptions.Verbose {
		logging.LogInfo("Create root CA and add to gopass finished.")
	}

	return nil
}
