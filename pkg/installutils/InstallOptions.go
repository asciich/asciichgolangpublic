package installutils

import (
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type InstallOptions struct {
	SrcPath     string
	InstallPath string
	Mode        string
	Verbose     bool
}

func NewInstallOptions() (i *InstallOptions) {
	return new(InstallOptions)
}

func (i *InstallOptions) IsSourcePathSet() (isSet bool) {
	return i.SrcPath != ""
}

func (i *InstallOptions) IsModeSet() (isSet bool) {
	return i.Mode != ""
}

func (i *InstallOptions) GetInstallPath() (installPath string, err error) {
	if i.InstallPath == "" {
		return "", tracederrors.TracedErrorf("InstallPath not set")
	}

	return i.InstallPath, nil
}

func (i *InstallOptions) GetMode() (mode string, err error) {
	if i.Mode == "" {
		return "", tracederrors.TracedErrorf("Mode not set")
	}

	return i.Mode, nil
}

func (i *InstallOptions) GetSrcPath() (srcPath string, err error) {
	if i.SrcPath == "" {
		return "", tracederrors.TracedErrorf("SrcPath not set")
	}

	return i.SrcPath, nil
}

func (i *InstallOptions) GetVerbose() (verbose bool) {

	return i.Verbose
}

func (i *InstallOptions) MustGetInstallPath() (installPath string) {
	installPath, err := i.GetInstallPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return installPath
}

func (i *InstallOptions) MustGetMode() (mode string) {
	mode, err := i.GetMode()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return mode
}

func (i *InstallOptions) MustGetSrcPath() (srcPath string) {
	srcPath, err := i.GetSrcPath()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return srcPath
}

func (i *InstallOptions) MustSetInstallPath(installPath string) {
	err := i.SetInstallPath(installPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (i *InstallOptions) MustSetMode(mode string) {
	err := i.SetMode(mode)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (i *InstallOptions) MustSetSrcPath(srcPath string) {
	err := i.SetSrcPath(srcPath)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (i *InstallOptions) SetInstallPath(installPath string) (err error) {
	if installPath == "" {
		return tracederrors.TracedErrorf("installPath is empty string")
	}

	i.InstallPath = installPath

	return nil
}

func (i *InstallOptions) SetMode(mode string) (err error) {
	if mode == "" {
		return tracederrors.TracedErrorf("mode is empty string")
	}

	i.Mode = mode

	return nil
}

func (i *InstallOptions) SetSrcPath(srcPath string) (err error) {
	if srcPath == "" {
		return tracederrors.TracedErrorf("srcPath is empty string")
	}

	i.SrcPath = srcPath

	return nil
}

func (i *InstallOptions) SetVerbose(verbose bool) {
	i.Verbose = verbose
}
