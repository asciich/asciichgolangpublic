package parameteroptions

import (
	"path/filepath"
	"regexp"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

const baseNameRegexString = "^[a-zA-Z0-9.]*$"

const secretPathRegexString = "^[a-zA-Z0-9.][a-zA-Z0-9.\\/]*$"

var baseNameRegex = regexp.MustCompile(baseNameRegexString)

var secretPathRegex = regexp.MustCompile(secretPathRegexString)

type GopassSecretOptions struct {
	SecretPath string

	Overwrite bool
	Verbose   bool
}

func NewGopassSecretOptions() (g *GopassSecretOptions) {
	return new(GopassSecretOptions)
}

func (g *GopassSecretOptions) GetBaseName() (baseName string, err error) {
	path, err := g.GetSecretPath()
	if err != nil {
		return "", err
	}

	baseName = filepath.Base(path)

	if baseName == "" {
		return "", tracederrors.TracedErrorf(
			"base name is empty string after evaluation of path='%s'", path,
		)
	}

	return baseName, nil
}

func (g *GopassSecretOptions) GetDeepCopy() (copy *GopassSecretOptions) {
	copy = NewGopassSecretOptions()

	*copy = *g

	return copy
}

func (g *GopassSecretOptions) GetDirName() (dirName string, err error) {
	path, err := g.GetSecretPath()
	if err != nil {
		return "", err
	}

	dirName = filepath.Dir(path)

	if dirName == "" {
		return "", tracederrors.TracedErrorf("dirName is empty string after evaluation.")
	}

	return dirName, nil
}

func (g *GopassSecretOptions) GetOverwrite() (overwrite bool) {

	return g.Overwrite
}

func (g *GopassSecretOptions) GetSecretPath() (secretPath string, err error) {
	secretPath = g.SecretPath

	if secretPath == "" {
		return "", tracederrors.TracedErrorf("SecretPath not set")
	}

	secretPath = stringsutils.TrimAllPrefix(secretPath, "/")
	if secretPath == "" {
		return "", tracederrors.TracedErrorf("secret path is empty string after evaluation of '%s'", g.SecretPath)
	}

	if !secretPathRegex.MatchString(secretPath) {
		return "", tracederrors.TracedErrorf("given secretPath '%s' does not match '%s'", secretPath, secretPathRegexString)
	}

	return secretPath, nil
}

func (g *GopassSecretOptions) GetVerbose() (verbose bool) {

	return g.Verbose
}

func (g *GopassSecretOptions) MustGetBaseName() (baseName string) {
	baseName, err := g.GetBaseName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return baseName
}

func (g *GopassSecretOptions) MustGetDirName() (dirName string) {
	dirName, err := g.GetDirName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return dirName
}

func (g *GopassSecretOptions) MustGetSecretPath() (secretPath string) {
	secretPath, err := g.GetSecretPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return secretPath
}

func (g *GopassSecretOptions) MustSetBaseName(newBaseName string) {
	err := g.SetBaseName(newBaseName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GopassSecretOptions) MustSetSecretPath(secretPath string) {
	err := g.SetSecretPath(secretPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GopassSecretOptions) SetBaseName(newBaseName string) (err error) {
	if newBaseName == "" {
		return tracederrors.TracedErrorEmptyString("newBaseName")
	}

	if !baseNameRegex.MatchString(newBaseName) {
		return tracederrors.TracedErrorf(
			"newBaseName '%s' does not match regex '%s'",
			newBaseName,
			baseNameRegexString,
		)
	}

	dirName, err := g.GetDirName()
	if err != nil {
		return err
	}

	newPath := filepath.Join(dirName, newBaseName)

	err = g.SetSecretPath(newPath)
	if err != nil {
		return err
	}

	return nil
}

func (g *GopassSecretOptions) SetOverwrite(overwrite bool) {
	g.Overwrite = overwrite
}

func (g *GopassSecretOptions) SetSecretPath(secretPath string) (err error) {
	if secretPath == "" {
		return tracederrors.TracedErrorf("secretPath is empty string")
	}

	secretPath = stringsutils.TrimAllPrefix(secretPath, "/")
	if secretPath == "" {
		return tracederrors.TracedErrorf("secret path is empty string after evaluation of '%s'", g.SecretPath)
	}

	if !secretPathRegex.MatchString(secretPath) {
		return tracederrors.TracedErrorf("given secretPath '%s' does not match '%s'", secretPath, secretPathRegexString)
	}

	g.SecretPath = secretPath

	return nil
}

func (g *GopassSecretOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}
