package asciichgolangpublic

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GitlabCreateDeployKeyOptions struct {
	Name          string
	WriteAccess   bool
	PublicKeyFile File
	Verbose       bool
}

func NewGitlabCreateDeployKeyOptions() (g *GitlabCreateDeployKeyOptions) {
	return new(GitlabCreateDeployKeyOptions)
}

func (g *GitlabCreateDeployKeyOptions) GetPublicKeyFile() (publicKeyFile File, err error) {
	if g.PublicKeyFile == nil {
		return nil, tracederrors.TracedErrorf("PublicKeyFile not set")
	}

	return g.PublicKeyFile, nil
}

func (g *GitlabCreateDeployKeyOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitlabCreateDeployKeyOptions) GetWriteAccess() (writeAccess bool, err error) {

	return g.WriteAccess, nil
}

func (g *GitlabCreateDeployKeyOptions) MustGetName() (name string) {
	name, err := g.GetName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return name
}

func (g *GitlabCreateDeployKeyOptions) MustGetPublicKeyFile() (publicKeyFile File) {
	publicKeyFile, err := g.GetPublicKeyFile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return publicKeyFile
}

func (g *GitlabCreateDeployKeyOptions) MustGetPublicKeyMaterialString() (keyMaterial string) {
	keyMaterial, err := g.GetPublicKeyMaterialString()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyMaterial
}

func (g *GitlabCreateDeployKeyOptions) MustGetPublicKeyfile() (keyFile File) {
	keyFile, err := g.GetPublicKeyfile()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyFile
}

func (g *GitlabCreateDeployKeyOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitlabCreateDeployKeyOptions) MustGetWriteAccess() (writeAccess bool) {
	writeAccess, err := g.GetWriteAccess()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return writeAccess
}

func (g *GitlabCreateDeployKeyOptions) MustSetName(name string) {
	err := g.SetName(name)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateDeployKeyOptions) MustSetPublicKeyFile(publicKeyFile File) {
	err := g.SetPublicKeyFile(publicKeyFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateDeployKeyOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateDeployKeyOptions) MustSetWriteAccess(writeAccess bool) {
	err := g.SetWriteAccess(writeAccess)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabCreateDeployKeyOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreateDeployKeyOptions) SetPublicKeyFile(publicKeyFile File) (err error) {
	if publicKeyFile == nil {
		return tracederrors.TracedErrorf("publicKeyFile is nil")
	}

	g.PublicKeyFile = publicKeyFile

	return nil
}

func (g *GitlabCreateDeployKeyOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (g *GitlabCreateDeployKeyOptions) SetWriteAccess(writeAccess bool) (err error) {
	g.WriteAccess = writeAccess

	return nil
}

func (o *GitlabCreateDeployKeyOptions) GetName() (name string, err error) {
	if len(o.Name) <= 0 {
		return "", tracederrors.TracedError("Name not set")
	}

	return o.Name, nil
}

func (o *GitlabCreateDeployKeyOptions) GetPublicKeyMaterialString() (keyMaterial string, err error) {
	keyFile, err := o.GetPublicKeyfile()
	if err != nil {
		return "", err
	}

	keyFilePath, err := keyFile.GetLocalPath()
	if err != nil {
		return "", err
	}

	keyMaterial, err = keyFile.ReadAsString()
	if err != nil {
		return "", err
	}

	keyMaterial = strings.TrimSpace(keyMaterial)
	if len(keyMaterial) <= 0 {
		return "", tracederrors.TracedErrorf(
			"Key material from '%s' failed. Got empty key material string",
			keyFilePath,
		)
	}

	return keyMaterial, nil
}

func (o *GitlabCreateDeployKeyOptions) GetPublicKeyfile() (keyFile File, err error) {
	if o.PublicKeyFile == nil {
		return nil, tracederrors.TracedError("PublicKeyFile is nil")
	}

	return o.PublicKeyFile, nil
}
