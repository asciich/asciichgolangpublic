package gopass

import (
	"context"
	"crypto"
	"crypto/x509"
	"fmt"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/binaryinfo"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/cryptoutils"
	"github.com/asciich/asciichgolangpublic/pkg/randomgenerator"
	"github.com/asciich/asciichgolangpublic/pkg/tlsutils/x509utils"
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

func Generate(ctx context.Context, credentialName string) (generatedCredential *GopassCredential, err error) {
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

	logging.LogInfoByCtxf(ctx, "Gopass credentail '%s' generated.", credentialName)

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
		contextutils.ContextSilent(),
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

func GetPrivateKey(ctx context.Context, getOptions *parameteroptions.GopassSecretOptions) (crypto.PrivateKey, error) {
	if getOptions == nil {
		return nil, tracederrors.TracedError("getOptions is nil")
	}

	credential, err := GetCredentialValueAsString(getOptions)
	if err != nil {
		return nil, err
	}

	key, err := cryptoutils.LoadPrivateKeyFromPEMString(credential)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func GetSslCertificate(ctx context.Context, getOptions *parameteroptions.GopassSecretOptions) (cert *x509.Certificate, err error) {
	if getOptions == nil {
		return nil, tracederrors.TracedError("getOptions is nil")
	}

	credential, err := GetCredentialValueAsString(getOptions)
	if err != nil {
		return nil, err
	}

	cert, err = x509utils.LoadCertificateFromPEMString(credential)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func InsertFileByString(ctx context.Context, fileContent string, gopassOptions *parameteroptions.GopassSecretOptions) (err error) {
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
		ctx,
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

	logging.LogChangedByCtxf(ctx, "Added file content to gopass as '%s'", gopassPath)

	return nil
}

func InsertFile(ctx context.Context, fileToInsert files.File, gopassOptions *parameteroptions.GopassSecretOptions) (err error) {
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

	fileExists, err := fileToInsert.Exists(contextutils.GetVerboseFromContext(ctx))
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
		ctx,
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

	logging.LogInfoByCtxf(ctx, "Added file '%v' to gopass as '%v'", fileToInsertPath, gopassPath)

	return nil
}

func InsertPrivateKey(ctx context.Context, privateKey crypto.PrivateKey, gopassOptions *parameteroptions.GopassSecretOptions) error {
	if privateKey == nil {
		return tracederrors.TracedErrorNil("privateKey")
	}

	if gopassOptions == nil {
		return tracederrors.TracedErrorNil("gopassOptions")
	}

	encodedKey, err := cryptoutils.EncodePrivateKeyAsPEMString(privateKey)
	if err != nil {
		return err
	}

	return InsertSecret(ctx, encodedKey, gopassOptions)
}

func InsertX509Certificate(ctx context.Context, cert *x509.Certificate, gopassOptions *parameteroptions.GopassSecretOptions) error {
	if cert == nil {
		return tracederrors.TracedErrorNil("cert")
	}

	if gopassOptions == nil {
		return tracederrors.TracedErrorNil("gopassOptions")
	}

	encodedCert, err := x509utils.EncodeCertificateAsPEMString(cert)
	if err != nil {
		return err
	}

	return InsertSecret(ctx, encodedCert, gopassOptions)
}

func InsertSecret(ctx context.Context, secretToInsert string, gopassOptions *parameteroptions.GopassSecretOptions) (err error) {
	if len(secretToInsert) <= 0 {
		return tracederrors.TracedErrorEmptyString("secretToInsert")
	}

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
		fmt.Sprintf("echo '%s' | gopass insert -f '%s'", secretToInsert, gopassPath),
	}

	_, err = commandexecutor.Bash().RunCommand(
		ctx,
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

	logging.LogInfoByCtxf(ctx, "Added credentail '%v' to gopass.", gopassPath)

	return nil
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

func Sync(ctx context.Context) (err error) {
	_, err = commandexecutor.Bash().RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"gopass", "sync"},
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
		contextutils.ContextSilent(),
		&parameteroptions.RunCommandOptions{
			Command: insertCommand,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func WriteSecretIntoTemporaryFile(ctx context.Context, getOptions *parameteroptions.GopassSecretOptions) (temporaryFile files.File, err error) {
	if getOptions == nil {
		return nil, tracederrors.TracedError("getOptions is nil")
	}

	credential, err := GetCredential(getOptions)
	if err != nil {
		return nil, err
	}

	temporaryFile, err = credential.WriteIntoTemporaryFile(ctx)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}
