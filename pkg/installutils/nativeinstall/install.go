package nativeinstall

import (
	"context"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/archiveutils/tarutils"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/httputils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpoptions"
	"github.com/asciich/asciichgolangpublic/pkg/installutils/installoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// isAlreadyInstalled returns true if the file at installPath already exists and
// its sha256 sum matches the expected one. Returns false when no checksum is
// given so the caller always performs the installation in that case.
func isAlreadyInstalled(ctx context.Context, installPath string, sha256Sum string) (bool, error) {
	if sha256Sum == "" {
		return false, nil
	}

	exists := nativefiles.Exists(ctx, installPath)
	if !exists {
		return false, nil
	}

	actualSha256Sum, err := checksumutils.GetSha256SumFromFile(ctx, installPath)
	if err != nil {
		return false, err
	}

	return actualSha256Sum == sha256Sum, nil
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

	if options.Sha256Sum != "" {
		actualSha256Sum, err := checksumutils.GetSha256SumFromFile(ctx, installPath)
		if err != nil {
			return err
		}

		if actualSha256Sum != options.Sha256Sum {
			return tracederrors.TracedErrorf(
				"Sha256Sum mismatch for installed file '%s'. Expected '%s' but got '%s'.",
				installPath,
				options.Sha256Sum,
				actualSha256Sum,
			)
		}
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' as '%s' finished.", sourcePath, installPath)

	return nil
}

func installFromSourceUrl(ctx context.Context, options *installoptions.InstallOptions) error {
	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	if options.SrcArchivePath != "" {
		return installFromSourceUrlArchive(ctx, options)
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

	archiveMember, err := options.GetSourceArchivePath()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' (archive member '%s') as '%s' started.", sourceUrl, archiveMember, installPath)

	alreadyInstalled, err := isAlreadyInstalled(ctx, installPath, options.Sha256Sum)
	if err != nil {
		return err
	}
	if alreadyInstalled {
		logging.LogInfoByCtxf(ctx, "Install '%s' (archive member '%s') as '%s' finished. Already in place with matching checksum.", sourceUrl, archiveMember, installPath)
		return nil
	}

	tempDir, err := tempfiles.CreateTempDir(ctx)
	if err != nil {
		return err
	}
	defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

	downloadedArchive := filepath.Join(tempDir, "downloaded_archive")
	_, err = httputils.DownloadAsFile(ctx, &httpoptions.DownloadAsFileOptions{
		RequestOptions:    &httpoptions.RequestOptions{Url: sourceUrl},
		OutputPath:        downloadedArchive,
		OverwriteExisting: true,
	})
	if err != nil {
		return err
	}

	// ExtractFileFromTarArchive auto-detects gzip by magic bytes, so both
	// .tar and .tar.gz archives are supported.
	extractedPath := filepath.Join(tempDir, "extracted")
	err = tarutils.ExtractFileFromTarArchive(ctx, downloadedArchive, archiveMember, extractedPath)
	if err != nil {
		return err
	}

	// Reuse the local-file install path so Mode + Sha256Sum validation
	// (of the extracted file) live in one place.
	err = installFromSourcePath(ctx, &installoptions.InstallOptions{
		SrcPath:         extractedPath,
		InstallPath:     options.InstallPath,
		Mode:            options.Mode,
		UseSudo:         options.UseSudo,
		ReplaceExisting: options.ReplaceExisting,
		Sha256Sum:       options.Sha256Sum,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install '%s' (archive member '%s') as '%s' finished.", sourceUrl, archiveMember, installPath)

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
