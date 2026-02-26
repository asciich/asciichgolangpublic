package installutils

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type InstallOptions struct {
	SrcPath     string
	InstallPath string
	Mode        string
	UseSudo     bool
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
