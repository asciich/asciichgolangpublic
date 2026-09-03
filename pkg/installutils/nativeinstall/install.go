package nativeinstall

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func installFromSourcePath(ctx context.Context, options *installoptions.InstallOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	sourcePath, err := options.GetSrcPath()
	if err != nil {
		return err
	}

	installPath, err := options.GetInstallPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' started.", sourcePath, installPath)

	err = nativefiles.Copy(ctx, sourcePath, installPath, &filesoptions.CopyOptions{
		UseSudo:         options.UseSudo,
		ReplaceExisting: options.ReplaceExisting,
	})
	if err != nil {
		return err
	}

	if options.IsModeSet() {
		permissions, err := options.GetMode()
		if err != nil {
			return err
		}

		err = nativefiles.Chmod(ctx, installPath, &filesoptions.ChmodOptions{
			PermissionsString: permissions,
		})
		if err != nil {
			return err
		}
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' finished.", sourcePath, installPath)

	return nil
}

func installFromSourceUrl(ctx context.Context, options *installoptions.InstallOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	sourceUrl, err := options.GetSrcUrl()
	if err != nil {
		return err
	}

	installPath, err := options.GetInstallPath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' started.", sourceUrl, installPath)

	_, err = httputils.DownloadAsFile(ctx, &httpoptions.DownloadAsFileOptions{
		RequestOptions: &httpoptions.RequestOptions{
			Url: sourceUrl,
		},
		OutputPath:        installPath,
		OverwriteExisting: options.ReplaceExisting,
		Sha256Sum:         options.Sha256Sum,
		UseSudo:           options.UseSudo,
		PermissionsString: options.Mode,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' finished.", sourceUrl, installPath)

	return nil
}

func Install(ctx context.Context, options *installoptions.InstallOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.IsSourcePathSet() {
		return installFromSourcePath(ctx, options)
	}

	if options.IsSourceUrlSet() {
		return installFromSourceUrl(ctx, options)
	}

	return tracederrors.TracedError("No source to install set.")
}
