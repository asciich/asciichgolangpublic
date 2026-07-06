package installutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

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

func installFromSourceUrlArchive(ctx context.Context, options *installoptions.InstallOptions) error {
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

	archivePath, err := options.GetSourceArchivePath()
	if err != nil {
		return err
	}

	targetFile, err := nativefilesoo.NewFileByPath(installPath)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' from archive '%s' as '%s' started.", archivePath, sourceUrl, installPath)

	exists, err := targetFile.Exists(ctx)
	if err != nil {
		return err
	}

	var isMatching bool
	if exists {
		if options.Sha256Sum != "" {
			isMatching, err = targetFile.IsMatchingSha256Sum(options.Sha256Sum)
			if err != nil {
				return err
			}
		}
	}

	if isMatching {
		logging.LogInfoByCtxf(ctx, "Target file '%s' is already matching the sha256sum '%s'.", installPath, options.Sha256Sum)
	} else {
		downloadedArchive, err := httputils.DownloadAsTemporaryFile(
			ctx,
			&httpoptions.DownloadAsTemporaryFileOptions{
				RequestOptions: &httpoptions.RequestOptions{
					Url: options.SrcUrl,
				},
			},
		)
		if err != nil {
			return err
		}
		defer downloadedArchive.Delete(ctx, &filesoptions.DeleteOptions{})

		downloadedArchivePath, err := downloadedArchive.GetPath()
		if err != nil {
			return err
		}

		var extractPath string
		if options.UseSudo {
			extractPath, err = tempfiles.CreateTemporaryFile(ctx)
			if err != nil {
				return err
			}
		} else {
			extractPath = installPath
		}

		err = tarutils.ExtractFileFromTarArchive(ctx, downloadedArchivePath, archivePath, extractPath)
		if err != nil {
			return err
		}

		if options.UseSudo {
			nativefiles.Move(ctx, extractPath, installPath, &filesoptions.MoveOptions{UseSudo: options.UseSudo})
		}
	}

	if options.Mode != "" {
		err = nativefiles.Chmod(ctx, installPath, &filesoptions.ChmodOptions{
			PermissionsString: options.Mode,
			UseSudo:           options.UseSudo,
		})
		if err != nil {
			return err
		}
	}

	if options.Sha256Sum != "" {
		expectedSha256 := options.Sha256Sum

		logging.LogInfoByCtxf(ctx, "Going to validate installed file '%s' using expected sha256sum %s", installPath, expectedSha256)

		sha256, err := targetFile.GetSha256Sum(ctx)
		if err != nil {
			return err
		}

		if expectedSha256 == sha256 {
			logging.LogInfoByCtxf(ctx, "Installed file '%s' matches expected sha256sum %s", installPath, expectedSha256)
		} else {
			return tracederrors.TracedErrorf(
				"%w: Installed file '%s' has checksum '%s' and is not matching expected '%s'.",
				httpgeneric.ErrChecksumMismatch,
				installPath,
				sha256,
				expectedSha256,
			)
		}
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' from archive '%s' as '%s' finished.", archivePath, sourceUrl, installPath)

	return nil
}

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

func Install(ctx context.Context, options *installoptions.InstallOptions) (err error) {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.IsSourcePathSet() {
		return installFromSourcePath(ctx, options)
	}

	if options.IsSourceUrlSet() {
		if options.SrcArchivePath == "" {
			return installFromSourceUrl(ctx, options)
		} else {
			return installFromSourceUrlArchive(ctx, options)
		}
	}

	return tracederrors.TracedError("No source to install set.")
}
