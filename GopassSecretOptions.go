package asciichgolangpublic

import (
	"path/filepath"
	"strings"

	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type GopassSecretOptions struct {
	SecretRootDirectoryPath string
	SecretBasename          string

	Overwrite bool
	Verbose   bool
}

func NewGopassSecretOptions() (gopassSecretOptions *GopassSecretOptions) {
	return new(GopassSecretOptions)
}

func (g *GopassSecretOptions) GetOverwrite() (overwrite bool, err error) {

	return g.Overwrite, nil
}

func (g *GopassSecretOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GopassSecretOptions) MustGetGopassPath() (gopassPath string) {
	gopassPath, err := g.GetGopassPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gopassPath
}

func (g *GopassSecretOptions) MustGetOverwrite() (overwrite bool) {
	overwrite, err := g.GetOverwrite()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return overwrite
}

func (g *GopassSecretOptions) MustGetSecretBasename() (basename string) {
	basename, err := g.GetSecretBasename()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return basename
}

func (g *GopassSecretOptions) MustGetSecretRootDirectoryPath() (rootDirectoryPath string) {
	rootDirectoryPath, err := g.GetSecretRootDirectoryPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return rootDirectoryPath
}

func (g *GopassSecretOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GopassSecretOptions) MustSetGopassPath(fullPath string) {
	err := g.SetGopassPath(fullPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GopassSecretOptions) MustSetOverwrite(overwrite bool) {
	err := g.SetOverwrite(overwrite)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GopassSecretOptions) MustSetSecretBasename(secretBasename string) {
	err := g.SetSecretBasename(secretBasename)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GopassSecretOptions) MustSetSecretRootDirectoryPath(secretRootDirectoryPath string) {
	err := g.SetSecretRootDirectoryPath(secretRootDirectoryPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GopassSecretOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GopassSecretOptions) SetOverwrite(overwrite bool) (err error) {
	g.Overwrite = overwrite

	return nil
}

func (g *GopassSecretOptions) SetSecretBasename(secretBasename string) (err error) {
	if secretBasename == "" {
		return tracederrors.TracedErrorf("secretBasename is empty string")
	}

	g.SecretBasename = secretBasename

	return nil
}

func (g *GopassSecretOptions) SetSecretRootDirectoryPath(secretRootDirectoryPath string) (err error) {
	if secretRootDirectoryPath == "" {
		return tracederrors.TracedErrorf("secretRootDirectoryPath is empty string")
	}

	g.SecretRootDirectoryPath = secretRootDirectoryPath

	return nil
}

func (g *GopassSecretOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}

func (o *GopassSecretOptions) GetDeepCopy() (copy *GopassSecretOptions) {
	copy = new(GopassSecretOptions)

	*copy = *o

	return copy
}

func (o *GopassSecretOptions) GetGopassPath() (gopassPath string, err error) {
	rootDir, err := o.GetSecretRootDirectoryPath()
	if err != nil {
		return "", err
	}

	basename, err := o.GetSecretBasename()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(basename, "/") {
		return "", tracederrors.TracedErrorf("absolute secret gopass paths not allowed, but got: '%v'", basename)
	}

	gopassPath = filepath.Join(rootDir, basename)
	return gopassPath, nil
}

func (o *GopassSecretOptions) GetSecretBasename() (basename string, err error) {
	basename = o.SecretBasename
	basename = strings.TrimSpace(basename)
	if len(basename) <= 0 {
		return "", tracederrors.TracedError("basename is empty string")
	}

	if strings.HasPrefix(basename, "/") {
		return "", tracederrors.TracedErrorf("absolute secret basenames not allowed for gopass, but got: '%v'", basename)
	}

	return basename, nil
}

func (o *GopassSecretOptions) GetSecretRootDirectoryPath() (rootDirectoryPath string, err error) {
	rootDirectoryPath = o.SecretRootDirectoryPath
	rootDirectoryPath = strings.TrimSpace(rootDirectoryPath)
	if len(rootDirectoryPath) <= 0 {
		return "", tracederrors.TracedError("rootDirectoryPath is empty string")
	}

	if strings.HasPrefix(rootDirectoryPath, "/") {
		return "", tracederrors.TracedErrorf("absolute secret rootDirectoryPaths not allowed for gopass, but got: '%v'", rootDirectoryPath)
	}

	return rootDirectoryPath, nil
}

func (o *GopassSecretOptions) SetGopassPath(fullPath string) (err error) {
	if len(fullPath) <= 0 {
		return tracederrors.TracedError("fullPath is empty string")
	}

	o.SecretBasename = filepath.Base(fullPath)
	o.SecretRootDirectoryPath = filepath.Dir(fullPath)

	return nil
}
