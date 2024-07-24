package asciichgolangpublic

import (
	"path/filepath"
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
		return "", TracedErrorf("BinaryName not set")
	}

	return i.BinaryName, nil
}

func (i *InstallOptions) GetInstallBashCompletion() (installBashCompletion bool) {

	return i.InstallBashCompletion
}

func (i *InstallOptions) GetInstallationPath() (installationPath string, err error) {
	if i.InstallationPath == "" {
		return "", TracedErrorf("InstallationPath not set")
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

func (i *InstallOptions) GetSourceFile() (sourceFile File, err error) {
	sourcePath, err := i.GetSourcePath()
	if err != nil {
		return nil, err
	}

	sourceFile, err = GetLocalFileByPath(sourcePath)
	if err != nil {
		return nil, err
	}

	return sourceFile, nil
}

func (i *InstallOptions) GetSourcePath() (sourcePath string, err error) {
	if i.SourcePath == "" {
		return "", TracedErrorf("SourcePath not set")
	}

	return i.SourcePath, nil
}

func (i *InstallOptions) GetUseSudoToInstall() (useSudoToInstall bool) {

	return i.UseSudoToInstall
}

func (i *InstallOptions) GetVerbose() (verbose bool) {

	return i.Verbose
}

func (i *InstallOptions) MustGetBinaryName() (binaryName string) {
	binaryName, err := i.GetBinaryName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return binaryName
}

func (i *InstallOptions) MustGetInstallationPath() (installationPath string) {
	installationPath, err := i.GetInstallationPath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return installationPath
}

func (i *InstallOptions) MustGetInstallationPathOrDefaultIfUnset() (installationPath string) {
	installationPath, err := i.GetInstallationPathOrDefaultIfUnset()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return installationPath
}

func (i *InstallOptions) MustGetSourceFile() (sourceFile File) {
	sourceFile, err := i.GetSourceFile()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sourceFile
}

func (i *InstallOptions) MustGetSourcePath() (sourcePath string) {
	sourcePath, err := i.GetSourcePath()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return sourcePath
}

func (i *InstallOptions) MustSetBinaryName(binaryName string) {
	err := i.SetBinaryName(binaryName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (i *InstallOptions) MustSetInstallationPath(installationPath string) {
	err := i.SetInstallationPath(installationPath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (i *InstallOptions) MustSetSourcePath(sourcePath string) {
	err := i.SetSourcePath(sourcePath)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (i *InstallOptions) SetBinaryName(binaryName string) (err error) {
	if binaryName == "" {
		return TracedErrorf("binaryName is empty string")
	}

	i.BinaryName = binaryName

	return nil
}

func (i *InstallOptions) SetInstallBashCompletion(installBashCompletion bool) {
	i.InstallBashCompletion = installBashCompletion
}

func (i *InstallOptions) SetInstallationPath(installationPath string) (err error) {
	if installationPath == "" {
		return TracedErrorf("installationPath is empty string")
	}

	i.InstallationPath = installationPath

	return nil
}

func (i *InstallOptions) SetSourcePath(sourcePath string) (err error) {
	if sourcePath == "" {
		return TracedErrorf("sourcePath is empty string")
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
