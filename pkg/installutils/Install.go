package installutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func getSourceFileAndPathFromOptions(options *InstallOptions) (sourceFile filesinterfaces.File, sourcePath string, err error) {
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

func getInstallFileFromOptions(options *InstallOptions) (installFile filesinterfaces.File, installPath string, err error) {
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

func setInstallFileAccessByOptions(ctx context.Context, installFile filesinterfaces.File, options *InstallOptions) (err error) {
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
		ctx,
		&parameteroptions.ChmodOptions{
			PermissionsString: mode,
		},
	)
	if err != nil {
		return nil
	}

	return nil
}

func installFromSourcePath(ctx context.Context, options *InstallOptions) (err error) {
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

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' started.", sourcePath, installPath)

	err = sourceFile.CopyToFile(installFile, contextutils.GetVerboseFromContext(ctx))
	if err != nil {
		return err
	}

	err = setInstallFileAccessByOptions(ctx, installFile, options)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' finished.", sourcePath, installPath)

	return nil
}

func Install(ctx context.Context, options *InstallOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.IsSourcePathSet() {
		return installFromSourcePath(ctx, options)
	}

	return tracederrors.TracedError("No source to install set.")
}
