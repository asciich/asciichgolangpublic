package install

import (
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func MustInstall(options *InstallOptions) {
	err := Install(options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func getSourceFileAndPathFromOptions(options *InstallOptions) (sourceFile files.File, sourcePath string, err error) {
	if options == nil {
		return nil, "", tracederrors.TracedErrorNil("options")
	}

	sourcePath, err = options.GetSrcPath()
	if err != nil {
		return nil, "", err
	}

	sourceFile, err = files.GetLocalFileByPath(sourcePath)
	if err != nil {
		return nil, "", err
	}

	return sourceFile, sourcePath, nil
}

func getInstallFileFromOptions(options *InstallOptions) (installFile files.File, installPath string, err error) {
	if options == nil {
		return nil, "", tracederrors.TracedErrorNil("options")
	}

	installPath, err = options.GetInstallPath()
	if err != nil {
		return nil, "", err
	}

	installFile, err = files.GetLocalFileByPath(installPath)
	if err != nil {
		return nil, "", err
	}

	return installFile, installPath, nil
}

func setInstallFileAccessByOptions(installFile files.File, options *InstallOptions) (err error) {
	if installFile == nil {
		tracederrors.TracedErrorNil("installFile")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if !options.IsModeSet() {
		return nil
	}

	mode, err := options.GetMode()
	if err != nil {
		return nil
	}

	err = installFile.Chmod(
		&parameteroptions.ChmodOptions{
			Verbose:           options.Verbose,
			PermissionsString: mode,
		},
	)
	if err != nil {
		return nil
	}

	return nil
}

func installFromSourcePath(options *InstallOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	sourceFile, sourcePath, err := getSourceFileAndPathFromOptions(options)
	if err != nil {
		return err
	}

	installFile, installPath, err := getInstallFileFromOptions(options)
	if err != nil {
		return err
	}

	if options.Verbose {
		logging.LogInfof("Install '%s' as '%s' started.", sourcePath, installPath)
	}

	err = sourceFile.CopyToFile(installFile, options.Verbose)
	if err != nil {
		return err
	}

	err = setInstallFileAccessByOptions(installFile, options)
	if err != nil {
		return err
	}

	if options.Verbose {
		logging.LogInfof("Install '%s' as '%s' finished.", sourcePath, installPath)
	}

	return nil
}

func Install(options *InstallOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.IsSourcePathSet() {
		return installFromSourcePath(options)
	}

	return tracederrors.TracedError("No source to install set.")
}
