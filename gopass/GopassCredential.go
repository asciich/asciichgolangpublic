package gopass

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/tlsutils/x509utils"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GopassCredential struct {
	name string
}

func GetGopassCredentialByName(name string) (credential *GopassCredential, err error) {
	name = strings.TrimSpace(name)
	if len(name) <= 0 {
		return nil, tracederrors.TracedError("name is empty string")
	}

	credential = NewGopassCredential()
	err = credential.SetName(name)
	if len(name) <= 0 {
		return nil, err
	}

	return credential, nil
}

func MustGetGopassCredentialByName(name string) (credential *GopassCredential) {
	credential, err := GetGopassCredentialByName(name)
	if err != nil {
		logging.LogFatalf("GetGopassCredentialByName failed: '%v'", err)
	}

	return credential
}

func NewGopassCredential() (gopassCredential *GopassCredential) {
	return new(GopassCredential)
}

func (c *GopassCredential) Exists() (exists bool, err error) {
	path, err := c.GetName()
	if err != nil {
		return false, err
	}

	output, err := commandexecutor.Bash().RunCommand(
		&parameteroptions.RunCommandOptions{
			Command: []string{
				"bash",
				"-c",
				fmt.Sprintf(
					"gopass cat '%s' > /dev/null || true",
					path,
				),
			},
		},
	)
	if err != nil {
		return false, err
	}

	stderr, err := output.GetStderrAsString()
	if err != nil {
		return false, err
	}

	return stderr == "", nil
}

func (c *GopassCredential) GetAsBytes() (credential []byte, err error) {
	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	credential, err = commandexecutor.Bash().RunCommandAndGetStdoutAsBytes(
		&parameteroptions.RunCommandOptions{
			Command: []string{"gopass", "cat", name},
		})
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (c *GopassCredential) GetAsInt() (value int, err error) {
	valueString, err := c.GetAsString()
	if err != nil {
		return -1, err
	}

	valueString = strings.TrimSpace(valueString)

	value, err = strconv.Atoi(valueString)
	if err != nil {
		return -1, err
	}

	return value, nil
}

func (c *GopassCredential) GetAsString() (credential string, err error) {
	credentialBytes, err := c.GetAsBytes()
	if err != nil {
		return "", err
	}

	credential = stringsutils.RemoveTailingNewline(string(credentialBytes))

	return credential, nil
}

func (c *GopassCredential) GetName() (name string, err error) {
	if len(c.name) <= 0 {
		return "", tracederrors.TracedError("name is not set")
	}

	return c.name, nil
}

func (c *GopassCredential) GetSslCertificate() (sslCert *x509utils.X509Certificate, err error) {
	contentBytes, err := c.GetAsBytes()
	if err != nil {
		return nil, err
	}

	sslCert = x509utils.NewX509Certificate()
	err = sslCert.LoadFromBytes(contentBytes)
	if err != nil {
		return nil, err
	}

	return sslCert, nil
}

func (c *GopassCredential) IncrementIntValue() (err error) {
	currentValue, err := c.GetAsInt()
	if err != nil {
		return err
	}

	err = c.SetByInt(currentValue + 1)
	if err != nil {
		return err
	}

	return err
}

func (c *GopassCredential) MustGetName() (name string) {
	name, err := c.GetName()
	if err != nil {
		logging.LogFatalf("gopassCredential.GetName failed: '%v'", err)
	}

	return name
}

func (c *GopassCredential) SetByInt(newValue int) (err error) {
	valueString := strconv.Itoa(newValue)
	err = c.SetByString(valueString)
	if err != nil {
		return err
	}

	return nil
}

func (c *GopassCredential) SetByString(newValue string) (err error) {
	if strings.Contains(newValue, "\n") {
		return tracederrors.TracedError("Unable to set copass value by string. newlines currenlty not supported.")
	}

	name, err := c.GetName()
	if err != nil {
		return err
	}

	insertCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf("echo '%s' | gopass insert -f '%s'", newValue, name),
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

func (c *GopassCredential) SetName(name string) (err error) {
	name = strings.TrimSpace(name)
	if len(name) <= 0 {
		return tracederrors.TracedError("name is empty string")
	}

	c.name = name

	return nil
}

func (c *GopassCredential) WriteIntoFile(ctx context.Context, outputFile files.File) (err error) {
	if outputFile == nil {
		return tracederrors.TracedError("outputFile is nil")
	}

	contentBytes, err := c.GetAsBytes()
	if err != nil {
		return err
	}

	err = outputFile.WriteBytes(contentBytes, contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return err
	}

	filePath, err := outputFile.GetLocalPath()
	if err != nil {
		return err
	}

	credentialName, err := c.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Wrote credential from gopass '%v' to file '%v'.", credentialName, filePath)

	return nil
}

func (c *GopassCredential) WriteIntoTemporaryFile(ctx context.Context) (temporaryFile files.File, err error) {
	temporaryFile, err = tempfiles.CreateEmptyTemporaryFile(contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return nil, err
	}

	err = c.WriteIntoFile(ctx, temporaryFile)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}
