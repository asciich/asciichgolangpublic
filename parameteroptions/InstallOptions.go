package parameteroptions

import (
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type InstallOptions struct {
	SourcePath            string
	BinaryName            string
	InstallationPath      string
	InstallBashCompletion bool
	UseSudoToInstall      bool
	Verbose               bool
}

func NewInstallOptions() (i *InstallOptions) {
	return new(InstallOptions)
}

func (i *InstallOptions) GetBinaryName() (binaryName string, err error) {
	if i.BinaryName == "" {
		return "", tracederrors.TracedErrorf("BinaryName not set")
	}

	return i.BinaryName, nil
}

func (i *InstallOptions) GetInstallBashCompletion() (installBashCompletion bool) {

	return i.InstallBashCompletion
}

func (i *InstallOptions) GetInstallationPath() (installationPath string, err error) {
	if i.InstallationPath == "" {
		return "", tracederrors.TracedErrorf("InstallationPath not set")
	}

	return i.InstallationPath, nil
}

func (i *InstallOptions) GetInstallationPathOrDefaultIfUnset() (installationPath string, err error) {
	if i.InstallationPath != "" {
		return i.InstallationPath, nil
	}

	binaryName, err := i.GetBinaryName()
	if err != nil {
		return "", err
	}

	installationPath = filepath.Join("/bin", binaryName)

	return installationPath, nil
}

func (i *InstallOptions) GetSourcePath() (sourcePath string, err error) {
	if i.SourcePath == "" {
		return "", tracederrors.TracedErrorf("SourcePath not set")
	}

	return i.SourcePath, nil
}

func (i *InstallOptions) GetUseSudoToInstall() (useSudoToInstall bool) {

	return i.UseSudoToInstall
}

func (i *InstallOptions) GetVerbose() (verbose bool) {

	return i.Verbose
}

func (i *InstallOptions) SetBinaryName(binaryName string) (err error) {
	if binaryName == "" {
		return tracederrors.TracedErrorf("binaryName is empty string")
	}

	i.BinaryName = binaryName

	return nil
}

func (i *InstallOptions) SetInstallBashCompletion(installBashCompletion bool) {
	i.InstallBashCompletion = installBashCompletion
}

func (i *InstallOptions) SetInstallationPath(installationPath string) (err error) {
	if installationPath == "" {
		return tracederrors.TracedErrorf("installationPath is empty string")
	}

	i.InstallationPath = installationPath

	return nil
}

func (i *InstallOptions) SetSourcePath(sourcePath string) (err error) {
	if sourcePath == "" {
		return tracederrors.TracedErrorf("sourcePath is empty string")
	}

	i.SourcePath = sourcePath

	return nil
}

func (i *InstallOptions) SetUseSudoToInstall(useSudoToInstall bool) {
	i.UseSudoToInstall = useSudoToInstall
}

func (i *InstallOptions) SetVerbose(verbose bool) {
	i.Verbose = verbose
}
