package asciichgolangpublic

import (
	"fmt"
	"strconv"
	"strings"
)

type GopassCredential struct {
	name string
}

func GetGopassCredentialByName(name string) (credential *GopassCredential, err error) {
	name = strings.TrimSpace(name)
	if len(name) <= 0 {
		return nil, TracedError("name is empty string")
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
		LogFatalf("GetGopassCredentialByName failed: '%v'", err)
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

	output, err := Bash().RunCommand(
		&RunCommandOptions{
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
		return false, nil
	}

	stderr, err := output.GetStderrAsString()
	if err != nil {
		return false, nil
	}

	return stderr == "", nil
}

func (c *GopassCredential) GetAsBytes() (credential []byte, err error) {
	name, err := c.GetName()
	if err != nil {
		return nil, err
	}

	credential, err = Bash().RunCommandAndGetStdoutAsBytes(
		&RunCommandOptions{
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

	credential = Strings().RemoveTailingNewline(string(credentialBytes))

	return credential, nil
}

func (c *GopassCredential) GetName() (name string, err error) {
	if len(c.name) <= 0 {
		return "", TracedError("name is not set")
	}

	return c.name, nil
}

func (c *GopassCredential) GetSslCertificate() (sslCert *X509Certificate, err error) {
	contentBytes, err := c.GetAsBytes()
	if err != nil {
		return nil, err
	}

	sslCert = NewX509Certificate()
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
		LogFatalf("gopassCredential.GetName failed: '%v'", err)
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
		return TracedError("Unable to set copass value by string. newlines currenlty not supported.")
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

func (c *GopassCredential) SetName(name string) (err error) {
	name = strings.TrimSpace(name)
	if len(name) <= 0 {
		return TracedError("name is empty string")
	}

	c.name = name

	return nil
}

func (c *GopassCredential) WriteIntoFile(outputFile File, verbose bool) (err error) {
	if outputFile == nil {
		return TracedError("outputFile is nil")
	}

	contentBytes, err := c.GetAsBytes()
	if err != nil {
		return err
	}

	err = outputFile.WriteBytes(contentBytes, verbose)
	if err != nil {
		return err
	}

	if verbose {
		filePath, err := outputFile.GetLocalPath()
		if err != nil {
			return err
		}

		credentialName, err := c.GetName()
		if err != nil {
			return err
		}

		LogInfof("Wrote credential from gopass '%v' to file '%v'.", credentialName, filePath)
	}

	return nil
}

func (c *GopassCredential) WriteIntoTemporaryFile(verbose bool) (temporaryFile File, err error) {
	temporaryFile, err = TemporaryFiles().CreateEmptyTemporaryFile(verbose)
	if err != nil {
		return nil, err
	}

	err = c.WriteIntoFile(temporaryFile, verbose)
	if err != nil {
		return nil, err
	}

	return temporaryFile, nil
}

func (g *GopassCredential) MustExists() (exists bool) {
	exists, err := g.Exists()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return exists
}

func (g *GopassCredential) MustGetAsBytes() (credential []byte) {
	credential, err := g.GetAsBytes()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credential
}

func (g *GopassCredential) MustGetAsInt() (value int) {
	value, err := g.GetAsInt()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return value
}

func (g *GopassCredential) MustGetAsString() (credential string) {
	credential, err := g.GetAsString()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return credential
}

func (g *GopassCredential) MustGetSslCertificate() (sslCert *X509Certificate) {
	sslCert, err := g.GetSslCertificate()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sslCert
}

func (g *GopassCredential) MustIncrementIntValue() {
	err := g.IncrementIntValue()
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassCredential) MustSetByInt(newValue int) {
	err := g.SetByInt(newValue)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassCredential) MustSetByString(newValue string) {
	err := g.SetByString(newValue)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassCredential) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassCredential) MustWriteIntoFile(outputFile File, verbose bool) {
	err := g.WriteIntoFile(outputFile, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GopassCredential) MustWriteIntoTemporaryFile(verbose bool) (temporaryFile File) {
	temporaryFile, err := g.WriteIntoTemporaryFile(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return temporaryFile
}
