package installoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type InstallOptions struct {
	SrcPath           string
	SrcUrl            string
	InstallPath       string
	Mode              string
	UseSudo           bool
	ReplaceExisting   bool
	Sha256Sum         string
	SkipTLSvalidation bool

	// Perform the download and validation in a local directory before the installation.
	//
	// This allows the installation on a machine which can not reach the SrcUrl because:
	//  1. The host running this command performs the download locally while the assumption is this host can reach the SrcUrl.
	//  2. The downloaded binary to install is validated by the checksums if specified.
	//  3. The downloaded binary is copied to the target host and installed.
	//  4. The installed file on the target host is validated by the checksums if specifed.
	ViaLocalTempDirectory bool
}

func NewInstallOptions() (i *InstallOptions) {
	return new(InstallOptions)
}

func (i *InstallOptions) IsSourcePathSet() bool {
	return i.SrcPath != ""
}

func (i *InstallOptions) IsModeSet() bool {
	return i.Mode != ""
}

func (i *InstallOptions) IsSourceUrlSet() bool {
	return i.SrcUrl != ""
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

func (i *InstallOptions) GetSrcUrl() (string, error) {
	if i.SrcUrl == "" {
		return "", tracederrors.TracedError("Src URL not set")
	}

	return i.SrcUrl, nil
}
