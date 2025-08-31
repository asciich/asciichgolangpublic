package asciichgolangpublic

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitlabCreateDeployKeyOptions struct {
	Name          string
	WriteAccess   bool
	PublicKeyFile filesinterfaces.File
}

func NewGitlabCreateDeployKeyOptions() (g *GitlabCreateDeployKeyOptions) {
	return new(GitlabCreateDeployKeyOptions)
}

func (g *GitlabCreateDeployKeyOptions) GetPublicKeyFile() (publicKeyFile filesinterfaces.File, err error) {
	if g.PublicKeyFile == nil {
		return nil, tracederrors.TracedErrorf("PublicKeyFile not set")
	}

	return g.PublicKeyFile, nil
}

func (g *GitlabCreateDeployKeyOptions) GetWriteAccess() (writeAccess bool, err error) {

	return g.WriteAccess, nil
}

func (g *GitlabCreateDeployKeyOptions) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorf("name is empty string")
	}

	g.Name = name

	return nil
}

func (g *GitlabCreateDeployKeyOptions) SetPublicKeyFile(publicKeyFile filesinterfaces.File) (err error) {
	if publicKeyFile == nil {
		return tracederrors.TracedErrorf("publicKeyFile is nil")
	}

	g.PublicKeyFile = publicKeyFile

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

func (o *GitlabCreateDeployKeyOptions) GetPublicKeyfile() (keyFile filesinterfaces.File, err error) {
	if o.PublicKeyFile == nil {
		return nil, tracederrors.TracedError("PublicKeyFile is nil")
	}

	return o.PublicKeyFile, nil
}
